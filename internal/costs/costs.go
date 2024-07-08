package costs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/imdario/mergo"
	"github.com/lib/pq"
	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/users"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var (
	mu                             sync.Mutex
	ErrInvalidCostItemName         = fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	ErrInvalidCostItemInvoiceLink  = fmt.Errorf("Unable to use the provided invoice link, as it appears to not be a valid link")
	ErrInvalidCostItemInvoiceDate  = fmt.Errorf("Unable to use the provided invoice date, because it is invalid")
	ErrInvalidCostItemReoccurUntil = fmt.Errorf("Unable to use the provided reoccur until, because it is invalid")
	ErrInvalidCostItemFrequency    = fmt.Errorf("Unable to use the provided frequency, because it is invalid")
)

type Manager struct {
	db    *sql.DB
	users *users.Manager
}

func NewManager(db *sql.DB, users *users.Manager) *Manager {
	return &Manager{
		db:    db,
		users: users,
	}
}
func (m *Manager) Validate(item types.CostItem) error {
	if len(item.Title) == 0 || len(item.Title) >= 30 || item.Title == "" {
		return ErrInvalidCostItemName
	}
	if u, err := url.Parse(item.InvoiceLink); u == nil || err != nil {
		return ErrInvalidCostItemInvoiceLink
	}
	if item.InvoiceDate == 0 {
		return ErrInvalidCostItemInvoiceDate
	}
	if item.ReoccurUntil != 0 && item.ReoccurUntil < item.InvoiceDate {
		log.Println(item.ReoccurUntil, item.InvoiceDate)
		return ErrInvalidCostItemReoccurUntil
	}
	switch item.Frequency {
	case
		types.FrequencyNever,
		types.FrequencyDaily,
		types.FrequencyWeekly,
		types.FrequencyFornightly,
		types.FrequencyMonthly:
	default:
		return ErrInvalidCostItemFrequency
	}
	// TODO validate fromID
	return nil
}

func (m *Manager) GetView() (costs types.Costs, err error) {
	sqlStatement := `select costs_view from costs_view()`
	rows, err := m.db.Query(sqlStatement)
	if err != nil {
		return types.Costs{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		costs, err = getCostsViewFromRows(rows)
		if err != nil {
			return types.Costs{}, err
		}
	}
	return costs, nil
}

func (m *Manager) Get(id string) (types.CostItem, error) {
	sqlStatement := `
        select * from costs
        where id = $1 and deletionTimestamp = 0`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return types.CostItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	item, err := getCostItemObjectFromRows(rows)
	if err != nil {
		return types.CostItem{}, err
	}
	return item, nil
}
func (m *Manager) GetExistsBatch(ids []string) (count int, err error) {
	sqlStatement := `
        select count(*) from costs
        where id = any($1)`
	rows, err := m.db.Query(sqlStatement, pq.Array(ids))
	if err != nil {
		return -1, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	if err := rows.Scan(&count); err != nil {
		return -1, err
	}
	if err := rows.Err(); err != nil {
		return -1, err
	}
	return count, nil
}

func (m *Manager) List(options types.CostListOptions) (costItems []types.CostItem, err error) {
	sqlStatement := `
        select * from costs
        where deletionTimestamp = 0 `
	fields := []interface{}{}

	if options.UnfulfilledReoccuring != nil && *options.UnfulfilledReoccuring {
		sqlStatement += ` and toId = '' `
	}

	if len(options.IDs) > 0 {
		sqlStatement += fmt.Sprintf(` and id = any($%v) `, len(fields)+1)
		fields = append(fields, pq.Array(options.IDs))
	}

	sqlStatement += ` order by creationTimestamp desc `

	if options.Limit > 0 {
		sqlStatement += fmt.Sprintf(` limit $%v `, len(fields)+1)
		fields = append(fields, options.Limit)
	}

	rows, err := m.db.Query(sqlStatement, fields...)
	if err != nil {
		return []types.CostItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		costItem, err := getCostItemObjectFromRows(rows)
		if err != nil {
			return []types.CostItem{}, err
		}
		costItems = append(costItems, costItem)
	}
	return costItems, nil
}

func (m *Manager) Create(item types.CostItem) (types.CostItem, error) {
	if err := m.Validate(item); err != nil {
		return types.CostItem{}, err
	}

	item.AuthorLast = item.Author
	if item.InvoicedBy == "" {
		item.InvoicedBy = item.Author
	}

	sqlStatement := `
        insert into costs (
          title,
          frequency,
          reoccurUntil,
          notes,
          amount,
          invoiceLink,
          invoiceDate,
          invoicedBy,
          fromId,
          toId,
          author,
          authorLast
        )
        values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        returning *`
	rows, err := m.db.Query(sqlStatement, item.Title, item.Frequency, item.ReoccurUntil, item.Notes, item.Amount, item.InvoiceLink, item.InvoiceDate, item.InvoicedBy, item.FromID, item.ToID, item.Author, item.AuthorLast)
	if err != nil {
		return types.CostItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	itemCreated, err := getCostItemObjectFromRows(rows)
	if err != nil {
		return types.CostItem{}, err
	}
	return itemCreated, nil
}

func (m *Manager) Update(id string, item types.CostItem) (itemUpdated types.CostItem, err error) {
	if err := m.Validate(item); err != nil {
		return types.CostItem{}, err
	}
	if _, err := m.Get(id); err != nil {
		return types.CostItem{}, err
	}
	sqlStatement := `
        update costs
        set
          title = $2,
          frequency = $3,
          reoccurUntil = $4,
          notes = $5,
          amount = $6,
          invoiceDate = $7,
          invoiceLink = $8,
          modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int
        where id = $1
        returning *`
	rows, err := m.db.Query(sqlStatement, id, item.Title, item.Frequency, item.ReoccurUntil, item.Notes, item.Amount, item.InvoiceDate, item.InvoiceLink)
	if err != nil {
		return types.CostItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		itemUpdated, err = getCostItemObjectFromRows(rows)
		if err != nil {
			return types.CostItem{}, err
		}
	}
	return itemUpdated, nil
}

func (m *Manager) Patch(id string, item types.CostItem) (itemPatched types.CostItem, err error) {
	existingItem, err := m.Get(id)
	if err != nil {
		return types.CostItem{}, err
	}
	if err := mergo.Merge(&item, existingItem); err != nil {
		return types.CostItem{}, err
	}
	if err := m.Validate(item); err != nil {
		return types.CostItem{}, err
	}
	sqlStatement := `
        update costs
        set
          title = $2,
          frequency = $3,
          reoccurUntil = $4,
          notes = $5,
          amount = $6,
          invoiceDate = $7,
          invoiceLink = $8,
          modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int
        where id = $1
        returning *`
	rows, err := m.db.Query(sqlStatement, id, item.Title, item.Frequency, item.ReoccurUntil, item.Notes, item.Amount, item.InvoiceDate, item.InvoiceLink)
	if err != nil {
		return types.CostItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		itemPatched, err = getCostItemObjectFromRows(rows)
		if err != nil {
			return types.CostItem{}, err
		}
	}
	return itemPatched, nil
}
func (m *Manager) PatchToIDLink(id string, toID string) (itemPatched types.CostItem, err error) {
	if _, err := m.Get(id); err != nil {
		return types.CostItem{}, err
	}
	sqlStatement := `
        update costs
        set
          toID = $2,
          modificationTimestamp = date_part('epoch',CURRENT_TIMESTAMP)::int
        where id = $1
        returning *`
	rows, err := m.db.Query(sqlStatement, id, toID)
	if err != nil {
		return types.CostItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		itemPatched, err = getCostItemObjectFromRows(rows)
		if err != nil {
			return types.CostItem{}, err
		}
	}
	return itemPatched, nil
}

func (m *Manager) unlinkItems(ids []string) error {
	sqlStatement := `update costs set toid = '' where toid = any($1)`
	rows, err := m.db.Query(sqlStatement, pq.Array(ids))
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	return nil
}
func (m *Manager) Delete(id string) error {
	item, err := m.Get(id)
	if err != nil {
		return err
	}
	sqlStatement := `delete from costs where id = $1`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	if err := m.unlinkItems([]string{item.ID}); err != nil {
		return err
	}
	return nil
}
func (m *Manager) DeleteBatch(ids []string) error {
	items, err := m.List(types.CostListOptions{
		IDs: ids,
	})
	if err != nil {
		return err
	}
	foundIDs := []string{}
	for _, item := range items {
		foundIDs = append(foundIDs, item.ID)
	}
	sqlStatement := `delete from costs where id = any($1)`
	rows, err := m.db.Query(sqlStatement, pq.Array(foundIDs))
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	if err := m.unlinkItems(foundIDs); err != nil {
		return err
	}
	return nil
}

// getCostItemObjectFromRows
// returns a shopping list object from rows
func getCostItemObjectFromRows(rows *sql.Rows) (item types.CostItem, err error) {
	if err := rows.Scan(
		&item.ID,
		&item.Title,
		&item.Frequency,
		&item.ReoccurUntil,
		&item.Notes,
		&item.Amount,
		&item.InvoiceLink,
		&item.InvoiceDate,
		&item.InvoicedBy,
		&item.FromID,
		&item.ToID,
		&item.Author,
		&item.AuthorLast,
		&item.CreationTimestamp,
		&item.ModificationTimestamp,
		&item.DeletionTimestamp,
	); err != nil {
		return types.CostItem{}, err
	}
	if err := rows.Err(); err != nil {
		return types.CostItem{}, err
	}
	return item, nil
}

// getCostItemViewFromRows
// returns a shopping list view from rows
func getCostsViewFromRows(rows *sql.Rows) (item types.Costs, err error) {
	var data []byte
	if err := rows.Scan(&data); err != nil {
		return types.Costs{}, err
	}
	if err := rows.Err(); err != nil {
		return types.Costs{}, err
	}
	if err := json.Unmarshal(data, &item); err != nil {
		return types.Costs{}, err
	}
	return item, nil
}

func (m *Manager) ScheduleNextCosts() error {
	mu.Lock()
	defer mu.Unlock()

	list, err := m.List(types.CostListOptions{
		UnfulfilledReoccuring: common.ToPointer(true),
	})
	if err != nil {
		return err
	}
	for _, item := range list {
		var nextOccurance time.Time
		lastOccurance := time.Unix(item.InvoiceDate, 0)
		scheduleNow := false
		switch item.Frequency {
		case types.FrequencyNever:
			continue
		case types.FrequencyDaily:
			nextOccurance = lastOccurance.AddDate(0, 0, 1)
			scheduleNow = time.Now().After(nextOccurance)
		case types.FrequencyWeekly:
			nextOccurance = lastOccurance.AddDate(0, 0, 7)
			scheduleNow = time.Now().After(nextOccurance)
		case types.FrequencyFornightly:
			nextOccurance = lastOccurance.AddDate(0, 0, 14)
			scheduleNow = time.Now().After(nextOccurance)
		case types.FrequencyMonthly:
			nextOccurance = lastOccurance.AddDate(0, 1, 0)
			scheduleNow = time.Now().After(nextOccurance)
		}
		if item.ReoccurUntil != 0 && nextOccurance.Unix() > item.ReoccurUntil {
			scheduleNow = false
		}
		if !scheduleNow {
			continue
		}
		log.Println("[scheduler] generating new cost item from", item.ID)
		item.InvoiceDate = nextOccurance.Unix()
		item.FromID = item.ID
		newItem, err := m.Create(item)
		if err != nil {
			return err
		}
		log.Println("new item:", newItem)
		if _, err := m.PatchToIDLink(item.ID, newItem.ID); err != nil {
			return err
		}
	}
	return nil
}

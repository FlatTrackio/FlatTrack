package costs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"gitlab.com/flattrack/flattrack/internal/users"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var (
	ErrInvalidCostItemName = fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
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
		costs, err = getCostsObjectFromRows(rows)
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

func (m *Manager) List(options types.CostListOptions) (costItems []types.CostItem, err error) {
	sqlStatement := `
        select * from cost
        where deletionTimestamp = 0 `
	fields := []interface{}{}

	sqlStatement += `order by creationTimestamp desc `

	if options.Limit > 0 {
		sqlStatement += fmt.Sprintf(`limit $%v `, len(fields)+1)
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

	sqlStatement := `
        insert into costs (title, paymentType, notes, amount, invoiceDate, invoicedBy, author, authorLast)
                         values ($1, $2, $3, $4, $5, $6, $7, $8)
                         returning *`
	rows, err := m.db.Query(sqlStatement, item.Title, item.PaymentType, item.Notes, item.Amount, item.InvoiceDate, item.InvoicedBy, item.Author, item.Author)
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

func (m *Manager) Update(id string, item types.CostItem) (types.CostItem, error) {
	return types.CostItem{}, nil
}

func (m *Manager) Patch(id string, item types.CostItem) (types.CostItem, error) {
	return types.CostItem{}, nil
}

func (m *Manager) Delete(id string) error {
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
	return nil
}

// getCostItemObjectFromRows
// returns a shopping list object from rows
func getCostItemObjectFromRows(rows *sql.Rows) (item types.CostItem, err error) {
	if err := rows.Scan(
		&item.ID,
		&item.Title,
		&item.PaymentType,
		&item.Notes,
		&item.Amount,
		&item.InvoiceDate,
		&item.InvoiceLink,
		&item.InvoicedBy,
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

// getCostItemObjectFromRows
// returns a shopping list object from rows
func getCostsObjectFromRows(rows *sql.Rows) (item types.Costs, err error) {
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

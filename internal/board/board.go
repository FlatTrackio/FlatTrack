package board

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.com/flattrack/flattrack/internal/users"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var (
	ErrInvalidBoardItemName = fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
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
func (m *Manager) Validate(item types.BoardItem) error {
	if len(item.Title) == 0 || len(item.Title) >= 30 || item.Title == "" {
		return ErrInvalidBoardItemName
	}
	return nil
}

func (m *Manager) Get(id string) (types.BoardItem, error) {
	sqlStatement := `select * from BoardItems where id = $1 and deletionTimestamp = 0`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return types.BoardItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	item, err := getBoardItemObjectFromRows(rows)
	if err != nil {
		return types.BoardItem{}, err
	}
	return item, nil
}

func (m *Manager) List(options types.BoardListOptions) (boardItems []types.BoardItem, err error) {
	sqlStatement := `select * from board_items where deletionTimestamp = 0 `
	fields := []interface{}{}

	sqlStatement += `order by creationTimestamp desc `

	if options.Limit > 0 {
		sqlStatement += fmt.Sprintf(`limit $%v `, len(fields)+1)
		fields = append(fields, options.Limit)
	}

	rows, err := m.db.Query(sqlStatement, fields...)
	if err != nil {
		return []types.BoardItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		boardItem, err := getBoardItemObjectFromRows(rows)
		if err != nil {
			return []types.BoardItem{}, err
		}
		boardItems = append(boardItems, boardItem)
	}
	return boardItems, nil
}

func (m *Manager) Create(item types.BoardItem) (types.BoardItem, error) {
	if err := m.Validate(item); err != nil {
		return types.BoardItem{}, err
	}

	sqlStatement := `insert into board_items (title, body, author)
                         values ($1, $2, $3, $4, $5, $6, $7)
                         returning *`
	rows, err := m.db.Query(sqlStatement, item.Title, item.Body, item.Author)
	if err != nil {
		return types.BoardItem{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	itemCreated, err := getBoardItemObjectFromRows(rows)
	if err != nil {
		return types.BoardItem{}, err
	}
	return itemCreated, nil
}

func (m *Manager) Update(id string, item types.BoardItem) (types.BoardItem, error) {
	return types.BoardItem{}, nil
}

func (m *Manager) Patch(id string, item types.BoardItem) (types.BoardItem, error) {
	return types.BoardItem{}, nil
}

func (m *Manager) Delete(id string) error {
	sqlStatement := `delete from board_items where id = $1`
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

// getBoardItemObjectFromRows
// returns a shopping list object from rows
func getBoardItemObjectFromRows(rows *sql.Rows) (item types.BoardItem, err error) {
	if err := rows.Scan(
		&item.ID,
		&item.Title,
		&item.Body,
		&item.Author,
		&item.CreationTimestamp,
		&item.ModificationTimestamp,
		&item.DeletionTimestamp,
	); err != nil {
		return types.BoardItem{}, err
	}
	if err := rows.Err(); err != nil {
		return types.BoardItem{}, err
	}
	return item, nil
}

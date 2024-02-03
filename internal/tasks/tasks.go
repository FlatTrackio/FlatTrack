/*
  tasks
    manage flat tasks
*/

// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package tasks

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"sync"
	"time"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/users"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var (
	ErrInvalidTaskName               = fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	ErrTargetStartTimestampBeforeNow = fmt.Errorf("Unable to use target timestamp because it is before the current time")
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

func (m *Manager) Validate(task types.Task) error {
	if len(task.Name) == 0 || len(task.Name) >= 30 || task.Name == "" {
		return ErrInvalidTaskName
	}
	if t := time.Unix(int64(task.TargetStartTimestamp), 0); t.Before(time.Now()) {
		return ErrTargetStartTimestampBeforeNow
	}
	return nil
}

func (m *Manager) Get(id string) (types.Task, error) {
	sqlStatement := `select * from tasks where id = $1 and deletionTimestamp = 0`
	rows, err := m.db.Query(sqlStatement, id)
	if err != nil {
		return types.Task{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	task, err := getTaskObjectFromRows(rows)
	if err != nil {
		return types.Task{}, err
	}
	return task, nil
}

func (m *Manager) List(options types.TaskListOptions) (tasks []types.Task, err error) {
	sqlStatement := `select * from tasks where deletionTimestamp = 0 `
	fields := []interface{}{}

	if options.Selector.ModificationTimestampAfter != 0 {
		sqlStatement += fmt.Sprintf(`and modificationTimestamp > $%v `, len(fields)+1)
		fields = append(fields, options.Selector.ModificationTimestampAfter)
	}
	if options.Selector.CreationTimestampAfter != 0 {
		sqlStatement += fmt.Sprintf(`and creationTimestamp > $%v `, len(fields)+1)
		fields = append(fields, options.Selector.CreationTimestampAfter)
	}
	if options.Selector.Name != "" {
		// NOTE needs to be split up like this so it can be handled properly by Sprintf
		sqlStatement += `and name ilike '%' ||` + fmt.Sprintf("$%v", len(fields)+1) + `|| '%' `
		fields = append(fields, options.Selector.Name)
	}
	if options.Selector.Paused != nil {
		sqlStatement += fmt.Sprintf(` and paused = $%v `, len(fields)+1)
		fields = append(fields, *options.Selector.Paused)
	}
	if options.Selector.Completed != nil {
		sqlStatement += fmt.Sprintf(` and completed = $%v `, len(fields)+1)
		fields = append(fields, *options.Selector.Completed)
	}
	if options.Selector.TemplateID != nil {
		sqlStatement += fmt.Sprintf(` and templateID = $%v `, len(fields)+1)
		fields = append(fields, *options.Selector.Completed)
	}

	switch options.SortBy {
	case types.TaskSortByRecentlyUpdated:
		sqlStatement += `order by modificationTimestamp desc `
	case types.TaskSortByLastUpdated:
		sqlStatement += `order by modificationTimestamp asc `
	case types.TaskSortByRecentlyAdded:
		sqlStatement += `order by creationTimestamp asc `
	case types.TaskSortByLastAdded:
		sqlStatement += `order by creationTimestamp asc `
	case types.TaskSortByAlphabeticalDescending:
		sqlStatement += `order by name asc `
	case types.TaskSortByAlphabeticalAscending:
		sqlStatement += `order by name desc `
	default:
		sqlStatement += `order by creationTimestamp desc `
	}

	if options.Limit > 0 {
		sqlStatement += fmt.Sprintf(`limit $%v `, len(fields)+1)
		fields = append(fields, options.Limit)
	}

	rows, err := m.db.Query(sqlStatement, fields...)
	if err != nil {
		return []types.Task{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		task, err := getTaskObjectFromRows(rows)
		if err != nil {
			return []types.Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (m *Manager) Create(task types.Task) (types.Task, error) {
	if err := m.Validate(task); err != nil {
		return types.Task{}, err
	}
	task.AuthorLast = task.Author

	sqlStatement := `insert into tasks (name, notes, assignee, assigneeType, targetStartTimestamp, frequency, templateId, author, authorLast)
                         values ($1, $2, $3, $4, $5, $6, $7)
                         returning *`
	rows, err := m.db.Query(sqlStatement, task.Name, task.Assignee, task.AssigneeType, task.TargetStartTimestamp, task.Frequency, task.TemplateID, task.Author, task.AuthorLast)
	if err != nil {
		return types.Task{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	rows.Next()
	taskCreated, err := getTaskObjectFromRows(rows)
	if err != nil {
		return types.Task{}, err
	}
	return taskCreated, nil
}

func (m *Manager) Update(id string, task types.Task) (types.Task, error) {
	return types.Task{}, nil
}

func (m *Manager) Patch(id string, task types.Task) (types.Task, error) {
	return types.Task{}, nil
}

func (m *Manager) Delete(id string) error {
	sqlStatement := `delete from tasks where id = $1`
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

// NOTE assign next logic

//   users                 : a, b, c, d, e
//   taskAssignedWithUsers :
//   nextAssignedUser      : a

//   users                 : a, b, c, d, e
//   taskAssignedWithUsers : 1:a
//   nextAssignedUser      : b

//   users                 : a, b, c, d, e
//   taskAssignedWithUsers : 1:a, 2:b
//   nextAssignedUser      : c

//   users                 : a, b, c, d, e
//   taskAssignedWithUsers : 1:a, 2:b, 3:c
//   nextAssignedUser      : d

//   users                 : a, b, c, d, e
//   taskAssignedWithUsers : 1:a, 2:b, 3:c, 4:d
//   nextAssignedUser      : e

//   users                 : a, b, c, d, e
//   taskAssignedWithUsers : 1:a, 2:b, 3:c, 4:d, 5:e
//   nextAssignedUser      : a

//   users                 : a, b, c, d
//   taskAssignedWithUsers : 1:a, 2:b, 3:c, 4:d, 5:e, 6:a
//   nextAssignedUser      : b

//   users                 : a, b, c
//   taskAssignedWithUsers : 1:a, 2:b, 3:c, 4:d, 5:e, 6:a, 7:b
//   nextAssignedUser      : c

//   users                 : a, b, c
//   taskAssignedWithUsers : 1:a, 2:b, 3:c, 4:d, 5:e, 6:a, 7:b, 8:c
//   nextAssignedUser      : a

func (m *Manager) AssignNext(task types.Task) (string, error) {
	tasks, err := m.List(types.TaskListOptions{
		SortBy: types.TaskSortByRecentlyAdded,
		Selector: types.TaskListSelector{
			Paused:     common.Pointer(false),
			TemplateID: &task.ID,
		},
	})
	if err != nil {
		return "", err
	}
	if len(tasks) == 0 && task.Assignee != "" {
		return task.Assignee, nil
	}
	users, err := m.users.List(false, types.UserSelector{
		Disabled: common.Pointer(false),
	})
	if err != nil {
		return "", err
	}
	userIDs := []string{}
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	lastAssignee := ""
	for _, t := range tasks {
		if t.Assignee == "" {
			continue
		}
		lastAssignee = t.Assignee
	}
	if lastAssignee == "" || len(tasks) == 0 {
		//nolint:gosec
		return users[rand.Intn(len(users))].ID, nil
	}
	i := slices.Index(userIDs, lastAssignee) + 1
	if i >= len(users) {
		i = 0
	}
	return userIDs[i%len(userIDs)], nil
}
func (m *Manager) AssignRandom(task types.Task) (string, error) {
	users, err := m.users.List(false, types.UserSelector{})
	if err != nil {
		return "", err
	}
	//nolint:gosec
	user := users[rand.Intn(len(users))]
	return user.ID, nil
}

func (m *Manager) AssignSelf(task types.Task) (string, error) {
	return task.Author, nil
}

func (m *Manager) Assign(task types.Task) (string, error) {
	switch task.AssigneeType {
	case types.TaskAssigneeNext:
		return m.AssignNext(task)
	case types.TaskAssigneeRandom:
		return m.AssignRandom(task)
	case types.TaskAssigneeSelf:
		return m.AssignSelf(task)
	}
	return "", nil
}

func (m *Manager) scheduleNextTask(task types.Task) error {
	// TODO list all tasks templated from this one
	// TODO create new task or use existing uncompleted and templated off this one

	instancesOfThisTask, err := m.List(types.TaskListOptions{
		SortBy: types.TaskSortByRecentlyAdded,
		Selector: types.TaskListSelector{
			Paused:     common.Pointer(false),
			Completed:  common.Pointer(false),
			TemplateID: &task.ID,
		},
	})
	if err != nil {
		return err
	}
	if len(instancesOfThisTask) == 0 {
		assignee, err := m.Assign(task)
		if err != nil {
			return err
		}
		latestInstance, err := m.Create(types.Task{
			Name:                 task.Name,
			Notes:                task.Notes,
			Assignee:             assignee,
			AssigneeType:         task.AssigneeType,
			TargetStartTimestamp: task.TargetStartTimestamp,
			Frequency:            task.Frequency,
			TemplateID:           task.TemplateID,
			Author:               task.Author,
			AuthorLast:           task.AuthorLast,
		})
		if err != nil {
			return err
		}
		task.LatestInstanceId = &latestInstance.ID
		if _, err := m.Update(task.ID, task); err != nil {
			return err
		}
		return nil
	}
	lastTask := instancesOfThisTask[0]
	if _, err := m.Update(lastTask.ID, types.Task{}); err != nil {
		return err
	}
	return nil
}

func (m *Manager) getNextDate(task types.Task) time.Time {
	switch task.Frequency {
	case types.TaskFrequencyOnce:
		return time.Unix(task.TargetStartTimestamp, 0)
	case types.TaskFrequencyDaily:
		return time.Unix(task.TargetStartTimestamp, 0).AddDate(0, 0, 1)
	case types.TaskFrequencyWeekly:
		return time.Unix(task.TargetStartTimestamp, 0).AddDate(0, 0, 7)
	case types.TaskFrequencyFortnightly:
		return time.Unix(task.TargetStartTimestamp, 0).AddDate(0, 0, 14)
	case types.TaskFrequencyMonthly:
		return time.Unix(task.TargetStartTimestamp, 0).AddDate(0, 1, 0)
	}
	return time.Time{}
}

func (m *Manager) Reconcile() error {
	list, err := m.List(types.TaskListOptions{
		Selector: types.TaskListSelector{
			Completed: common.Pointer(false),
		},
	})
	if err != nil {
		return err
	}
	// NOTE unsure current whether this needs to be synchronous or not
	var errs []error
	var wg sync.WaitGroup
	for _, task := range list {
		wg.Add(1)
		go func(task types.Task) {
			log.Printf("%+v\n", task)
			nextDate := m.getNextDate(task)
			// TODO check for not already assigned
			if task.LatestInstanceId == nil && (time.Now().Equal(nextDate) || time.Now().After(nextDate)) {
				if err := m.scheduleNextTask(task); err != nil {
					errs = append(errs, err)
				}
			}
			wg.Done()
		}(task)
	}
	wg.Wait()
	if len(errs) > 0 {
		return fmt.Errorf("errors reconciling tasks: %v", errs)
	}
	return nil
}

// getListObjectFromRows ...
// returns a shopping list object from rows
func getTaskObjectFromRows(rows *sql.Rows) (task types.Task, err error) {
	var assignee sql.NullString
	if err := rows.Scan(
		&task.ID,
		&task.Name,
		&task.Notes,
		&assignee,
		&task.AssigneeType,
		&task.TargetStartTimestamp,
		&task.Frequency,
		&task.TemplateID,
		&task.LatestInstanceId,
		&task.Paused,
		&task.Author,
		&task.AuthorLast,
		&task.Completed,
		&task.CreationTimestamp,
		&task.ModificationTimestamp,
		&task.DeletionTimestamp,
	); err != nil {
		return types.Task{}, err
	}
	if err := rows.Err(); err != nil {
		return types.Task{}, err
	}
	task.Assignee = assignee.String
	return task, nil
}

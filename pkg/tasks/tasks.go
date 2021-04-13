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
	"strings"

	"gitlab.com/flattrack/flattrack/pkg/users"
)

// TaskSpec ...
// task entry
type TaskSpec struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	Notes                 string        `json:"notes,omitempty"`
	Frequency             TaskFrequency `json:"frequency,omitempty"`
	Completed             bool          `json:"completed"`
	Assignee              string        `json:"assignee"`
	Rotation              TaskRotation  `json:"rotation"`
	RotatesBetween        []string      `json:"rotatesBetween"`
	StartDate             int           `json:"startDate"`
	TemplateID            string        `json:"templateId"`
	Author                string        `json:"author"`
	AuthorLast            string        `json:"authorLast"`
	CreationTimestamp     int           `json:"creationTimestamp"`
	ModificationTimestamp int           `json:"modificationTimestamp"`
	DeletionTimestamp     int           `json:"deletionTimestamp"`
}

// TaskListSelector ...
// filters for listing tasks
type TaskListSelector struct {
	Frequency string `json:"frequency"`
	Completed string `json:"completed"`
	Assignee  string `json:"assignee"`
	Rotation  string `json:"rotation"`
	Author    string `json:"author"`
}

// TaskListOptions ...
// options for listing tasks
type TaskListOptions struct {
	SortBy   string           `json:"sortBy"`
	Selector TaskListSelector `json:"selector"`
}

// TaskFrequency ...
// regularity of the task
type TaskFrequency string

// TaskFrequencies ...
// valid frequencies for tasks
const (
	TaskFrequencyDaily    TaskFrequency = "daily"
	TaskFrequencyWeekly   TaskFrequency = "weekly"
	TaskFrequencyBiWeekly TaskFrequency = "bi-weekly"
	TaskFrequencyMonthly  TaskFrequency = "monthly"
)

// TaskRotation ...
// how the task is rotated
type TaskRotation string

// TaskRotations ...
// valid rotations for a task
const (
	TaskRotationNever            TaskRotation = "never"
	TaskRotationOnNewAssignation TaskRotation = "onNewAssignation"
)

// TaskListSortType ...
// ways of sorting task lists
type TaskListSortType string

// TaskListSortTypes ...
// ways of sorting task lists
const (
	TaskListSortByRecentlyAdded          = "recentlyAdded"
	TaskListSortByRecentlyUpdated        = "recentlyUpdated"
	TaskListSortByLastAdded              = "lastAdded"
	TaskListSortByLastUpdated            = "lastUpdated"
	TaskListSortByAlphabeticalDescending = "alphabeticalDescending"
	TaskListSortByAlphabeticalAscending  = "alphabeticalAscending"
)

// ValidateTask ...
// validates a taskSpec
func ValidateTask(db *sql.DB, task TaskSpec) error {
	if len(task.Name) == 0 || len(task.Name) >= 30 || task.Name == "" {
		return fmt.Errorf("Unable to use the provided name, as it is either empty or too long or too short")
	}
	if len(task.Notes) > 100 {
		return fmt.Errorf("Unable to save task notes, as they are too long")
	}
	if task.TemplateID != "" {
		task, err := GetTask(db, task.TemplateID)
		if err != nil || task.ID == "" {
			return fmt.Errorf("Unable to find list to use as template from provided id")
		}
	}
	if user, err := users.GetUserByID(db, task.Assignee, false); err != nil || user.ID == "" {
		return fmt.Errorf("Unable to use user id with field assigned, as it is not found or invalid")
	}
	for _, includedUser := range task.RotatesBetween {
		if user, err := users.GetUserByID(db, includedUser, false); err != nil || user.ID == "" {
			return fmt.Errorf("Unable to use user id with field assigned, as it is not found or invalid")
		}
	}
	switch taskFrequency := task.Frequency; taskFrequency {
	case TaskFrequencyDaily, TaskFrequencyWeekly, TaskFrequencyBiWeekly, TaskFrequencyMonthly:
		break
	default:
		return fmt.Errorf("Unable to use task frequency '%v', as it is invalid", taskFrequency)
	}
	switch taskRotation := task.Rotation; taskRotation {
	case TaskRotationNever, TaskRotationOnNewAssignation:
		break
	default:
		return fmt.Errorf("Unable to use task rotation '%v', as it is invalid", taskRotation)
	}
	return nil
}

// GetTasks ...
// returns a list of tasks
func GetTasks(db *sql.DB, options TaskListOptions) (tasks []TaskSpec, err error) {
	// recentlyAdded
	sqlStatement := `select * from tasks
                         where deletionTimestamp = 0
	                 order by creationTimestamp desc`
	if options.SortBy == TaskListSortByRecentlyUpdated {
		sqlStatement = `select * from tasks
                         where deletionTimestamp = 0
	                 order by modificationTimestamp desc`
	} else if options.SortBy == TaskListSortByLastAdded {
		sqlStatement = `select * from tasks
                         where deletionTimestamp = 0
	                 order by creationTimestamp asc`
	} else if options.SortBy == TaskListSortByLastUpdated {
		sqlStatement = `select * from tasks
                         where deletionTimestamp = 0
	                 order by modificationTimestamp asc`
	} else if options.SortBy == TaskListSortByAlphabeticalDescending {
		sqlStatement = `select * from tasks
                         where deletionTimestamp = 0
	                 order by name asc`
	} else if options.SortBy == TaskListSortByAlphabeticalAscending {
		sqlStatement = `select * from tasks
                         where deletionTimestamp = 0
	                 order by name desc`
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		task, err := GetTaskObjectFromRows(rows)
		if err != nil {
			return tasks, err
		}
		if options.Selector.Completed == "true" && task.Completed != true {
			continue
		} else if options.Selector.Completed == "false" && task.Completed != false {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

// GetTask ...
// returns a given task
func GetTask(db *sql.DB, taskID string) (task TaskSpec, err error) {
	sqlStatement := `select * from tasks where id = $1 and deletionTimestamp = 0`
	rows, err := db.Query(sqlStatement, taskID)
	if err != nil {
		return task, err
	}
	defer rows.Close()
	rows.Next()
	task, err = GetTaskObjectFromRows(rows)
	if err != nil {
		return task, err
	}
	return task, err
}

// CreateTask ...
// creates a new task
func CreateTask(db *sql.DB, task TaskSpec) (taskInserted TaskSpec, err error) {
	err = ValidateTask(db, task)
	if err != nil {
		return taskInserted, err
	}

	task.AuthorLast = task.Author
	task.Completed = false

	rotatesBetween := strings.Join(task.RotatesBetween, " ")
	sqlStatement := `insert into tasks (name, notes, frequency, assignee, rotation, rotatesBetween, startDate, templateId, author, authorLast)
                         values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                         returning *`
	rows, err := db.Query(sqlStatement, task.Name, task.Notes, task.Frequency, task.Assignee, task.Rotation, rotatesBetween, task.StartDate, task.TemplateID, task.Author, task.AuthorLast)
	if err != nil {
		return taskInserted, err
	}
	defer rows.Close()
	rows.Next()
	taskInserted, err = GetTaskObjectFromRows(rows)
	if err != nil || taskInserted.ID == "" {
		return taskInserted, fmt.Errorf("Failed to create task")
	}
	if task.TemplateID == "" {
		return taskInserted, err
	}
	return taskInserted, err
}

// PatchTask ...
// patches a task
func PatchTask(db *sql.DB, task TaskSpec) (taskUpdated TaskSpec, err error) {
	return taskUpdated, err
}

// UpdateTask ...
// updates a task
func UpdateTask(db *sql.DB, task TaskSpec) (taskUpdated TaskSpec, err error) {
	return taskUpdated, err
}

// DeleteTask ...
// deletes a task
func DeleteTask(db *sql.DB, taskID string) (err error) {
	sqlStatement := `delete from tasks where id = $1`
	rows, err := db.Query(sqlStatement, taskID)
	defer rows.Close()
	return err
}

// GetTaskObjectFromRows ...
// returns a task object from rows
func GetTaskObjectFromRows(rows *sql.Rows) (task TaskSpec, err error) {
	var rotatesBetween string
	rows.Scan(&task.ID, &task.Name, &task.Notes, &task.Frequency, &task.Completed, &task.Assignee, &task.Rotation, &rotatesBetween, &task.StartDate, &task.TemplateID, &task.Author, &task.AuthorLast, &task.CreationTimestamp, &task.ModificationTimestamp, &task.DeletionTimestamp)
	err = rows.Err()
	task.RotatesBetween = strings.Split(rotatesBetween, " ")
	return task, err
}

/*
  tasks
    manage tasks
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

import Request from '@/requests/requests'

// GetTasks
// returns a list of all shopping lists
function GetTasks (
  completed,
  sortBy,
  creationTimestampAfter,
  modificationTimestampAfter,
  limit
) {
  return Request({
    url: '/api/apps/tasks/lists',
    method: 'GET',
    params: {
      completed,
      sortBy,
      creationTimestampAfter,
      modificationTimestampAfter,
      limit
    }
  })
}

// GetTask
// returns a shopping list
function GetTask (id) {
  return Request({
    url: `/api/apps/tasks/lists/${id}`,
    method: 'GET'
  })
}

// PostTask
// given a name and optional notes, create a shopping list
function PostTask (name, notes, templateId, templateListItemSelector) {
  return Request({
    url: '/api/apps/tasks/lists',
    method: 'POST',
    data: {
      name,
      notes,
      templateId
    }
  })
}

// PatchTask
// given a name and optional notes, patch a shopping list
function PatchTask (id, name, notes) {
  return Request({
    url: `/api/apps/tasks/lists/${id}`,
    method: 'PATCH',
    data: {
      name,
      notes
    }
  })
}

// UpdateTask
// given a name and optional notes, patch a shopping list
function UpdateTask (id, name, notes, completed, totalTagExclude) {
  return Request({
    url: `/api/apps/tasks/lists/${id}`,
    method: 'PUT',
    data: {
      name,
      notes,
      completed,
      totalTagExclude
    }
  })
}

// DeleteTask
// deletes a shopping list
function DeleteTask (id) {
  return Request({
    url: `/api/apps/tasks/lists/${id}`,
    method: 'DELETE'
  })
}

export default {
  GetTasks,
  GetTask,
  PostTask,
  PatchTask,
  UpdateTask,
  DeleteTask
}

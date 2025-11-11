/*
  settings
    manage admin settings
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

import Request from "@/requests/requests";

// PostFlatName
// changes the FlatName
function PostFlatName(flatName) {
  return Request({
    url: `/api/admin/settings/flatName`,
    method: "POST",
    data: {
      flatName,
    },
  });
}

// PutFlatNotes
// changes the FlatNotes
function PutFlatNotes(notes) {
  return Request({
    url: `/api/admin/settings/flatNotes`,
    method: "PUT",
    data: {
      notes,
    },
  });
}

// GetShoppingListKeepPolicy
// gets the policy for time to keep shopping lists
function GetShoppingListKeepPolicy() {
  return Request({
    url: `/api/admin/settings/shoppingListKeepPolicy`,
    method: "GET",
  });
}

// PutShoppingListKeepPolicy
// changes the time to keep shopping lists for
function PutShoppingListKeepPolicy(keepPolicy) {
  return Request({
    url: `/api/admin/settings/shoppingListKeepPolicy`,
    method: "PUT",
    data: {
      keepPolicy,
    },
  });
}

export default {
  PostFlatName,
  PutFlatNotes,
  GetShoppingListKeepPolicy,
  PutShoppingListKeepPolicy,
};

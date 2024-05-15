/*
  board
    manage board items
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

// GetBoardItems
// returns a list of all board items
function GetBoardItems(
  creationTimestampAfter,
  modificationTimestampAfter,
  limit
) {
  return Request({
    url: "/api/apps/board/items",
    method: "GET",
    params: {
      creationTimestampAfter,
      modificationTimestampAfter,
      limit,
    },
  });
}

// GetBoardItem
// returns a board items
function GetBoardItem(id) {
  return Request({
    url: `/api/apps/board/items/${id}`,
    method: "GET",
  });
}

// PostBoardItem
// given a name and optional notes, create a board items
function PostBoardItem(title, body) {
  return Request({
    url: "/api/apps/board/items",
    method: "POST",
    data: {
      title,
      body,
    },
  });
}

// PatchBoardItem
// given a name and optional notes, patch a board items
function PatchBoardItem(id, title, body) {
  return Request({
    url: `/api/apps/board/items/${id}`,
    method: "PATCH",
    data: {
      title,
      body,
    },
  });
}

// UpdateBoardItem
// given a name and optional notes, patch a board items
function UpdateBoardItem(id, title, body) {
  return Request({
    url: `/api/apps/board/items/${id}`,
    method: "PUT",
    data: {
      title,
      body,
    },
  });
}

// DeleteBoardItem
// deletes a board items
function DeleteBoardItem(id) {
  return Request({
    url: `/api/apps/board/items/${id}`,
    method: "DELETE",
  });
}

export default {
  GetBoardItems,
  GetBoardItem,
  PostBoardItem,
  PatchBoardItem,
  UpdateBoardItem,
  DeleteBoardItem,
};

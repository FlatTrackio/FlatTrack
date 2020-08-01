/*
  flatmates
    get your Flatmates
*/

import Request from '@/frontend/requests/requests'

// GetAllFlatmates
// get a list of all Flatmates
function GetAllFlatmates (id, notSelf, group) {
  return Request({
    url: '/api/users',
    method: 'GET',
    params: {
      id,
      notSelf,
      group
    }
  })
}

// GetFlatmate
// get a Flatmate
function GetFlatmate (id) {
  return Request({
    url: `/api/users/${id}`,
    method: 'GET'
  })
}

export default {
  GetAllFlatmates,
  GetFlatmate
}

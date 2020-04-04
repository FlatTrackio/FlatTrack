/*
  flatmates
    get your Flatmates
*/

import Request from '@/frontend/requests/requests'

// GetAllFlatmates
// get a list of all Flatmates
function GetAllFlatmates (params) {
  return Request({
    url: '/api/users',
    method: 'GET',
    params: params
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

export default { GetAllFlatmates, GetFlatmate }

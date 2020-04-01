/*
  flatmates
    get your Flatmates
*/

import Request from '@/frontend/requests/requests'

// GetAllFlatmates
// get a list of all Flatmates
function GetAllFlatmates () {
  return Request({
    url: '/api/flatmates',
    method: 'GET'
  })
}

// GetFlatmate
// get a Flatmate
function GetFlatmate (id) {
  return Request({
    url: `/api/flatmates/${id}`,
    method: 'GET'
  })
}

export default { GetAllFlatmates, GetFlatmate }

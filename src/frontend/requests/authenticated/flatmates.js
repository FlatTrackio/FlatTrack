/*
  flatmates
    get your Flatmates
*/

import Request from '@/frontend/requests/requests'

// GetAllFlatmates
// get a list of all Flatmates
function GetAllFlatmates (names, email, phoneNumber, birthday, groups, password) {
  return Request({
    url: '/api/users',
    method: 'GET',
    params: {
      names,
      email,
      phoneNumber,
      birthday,
      groups,
      password
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

export default { GetAllFlatmates, GetFlatmate }

/*
  groups
    fetch groups
*/

import Request from '@/frontend/requests/requests'

// GetGroups
// returns a list of groups
function GetGroups () {
  return Request({
    url: '/api/groups',
    method: 'GET'
  })
}

// GetGroup
// returns a group by id
function GetGroup (id) {
  return Request({
    url: `/api/groups/${id}`,
    method: 'GET'
  })
}

export default { GetGroups, GetGroup }

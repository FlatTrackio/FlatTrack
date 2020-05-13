/*
  can-i
    determine user account privilege
*/

import Request from '@/frontend/requests/requests'

// GetCanIgroup
// returns whether a user is in a group
function GetCanIgroup (name) {
  return Request({
    url: `/api/user/can-i/group/${name}`,
    method: 'GET'
  })
}

export default {
  GetCanIgroup
}

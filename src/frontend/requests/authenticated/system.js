/*
  system
    general instance information
*/

import Request from '@/frontend/requests/requests'

// GetVersion
// returns version information about the instance
function GetVersion (redirect) {
  return Request({
    url: '/api/system/version',
    method: 'GET'
  })
}

export default {
  GetVersion
}

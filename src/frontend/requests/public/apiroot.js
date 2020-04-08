/*
  apiroot
    get public metadata of API
*/

import Request from '@/frontend/requests/requests'

// GetAPIroot
// returns public metadata of API
function GetAPIroot () {
  return Request({
    url: '/api',
    method: 'GET'
  }, false, true)
}

export default { GetAPIroot }

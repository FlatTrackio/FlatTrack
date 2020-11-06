/*
  healthz
    check the health of an instance
*/

import Request from '@/requests/requests'

// GetHealthz
// returns if the instance is heathly
function GetHealthz () {
  return Request({
    url: '/_healthz',
    method: 'GET'
  })
}

export default {
  GetHealthz
}

/*
  flatInfo
    fetch info
*/

import Request from '@/frontend/requests/requests'

// GetFlatName
// gets the name of the flat
function GetFlatName () {
  return Request({
    url: '/api/system/flatName',
    method: 'GET'
  })
}

export default { GetFlatName }

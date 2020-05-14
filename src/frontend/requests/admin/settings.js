/*
  settings
    manage admin settings
*/

import Request from '@/frontend/requests/requests'

// PostFlatName
// changes the FlatName
function PostFlatName (flatName) {
  return Request({
    url: `/api/admin/settings/flatName`,
    method: 'POST',
    data: {
      flatName
    }
  })
}

export default {
  PostFlatName
}

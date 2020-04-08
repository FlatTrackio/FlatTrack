/*
  registration
    register the FlatTrack instance
*/

import Request from '@/frontend/requests/requests'

// GetInstanceRegistered
// determine if the instance of FlatTrack is registered
function GetInstanceRegistered () {
  return Request({
    url: '/api/system/initialized',
    method: 'GET'
  }, undefined, true)
}

// PostAdminRegister
// login and return a JWT
function PostAdminRegister (form) {
  return Request({
    url: '/api/admin/register',
    method: 'POST',
    data: form
  }, false, true)
}

export default { GetInstanceRegistered, PostAdminRegister }

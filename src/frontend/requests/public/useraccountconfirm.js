/*
  useraccountconfirm
    complete account registration
*/

import Request from '@/frontend/requests/requests'

// GetTokenValid
// returns if a confirm token is available
function GetTokenValid (id) {
  return Request({
    url: `/api/user/confirm/${id}`,
    method: 'GET'
  }, false, true)
}

// PostUserConfirm
// posts to confirm an account
function PostUserConfirm (id, secret, phoneNumber, birthday, password) {
  return Request({
    url: `/api/user/confirm/${id}`,
    method: 'POST',
    params: {
      secret
    },
    data: {
      phoneNumber,
      birthday,
      password
    }
  }, false, true)
}

export default { GetTokenValid, PostUserConfirm }

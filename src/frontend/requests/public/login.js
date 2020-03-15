import Request from '@/frontend/requests/requests'

function PostAPIauth (credentials) {
  return Request({
    url: '/api/login',
    method: 'POST',
    data: credentials
  })
}

function VerifyAuthToken () {
  return new Promise((resolve, reject) => {
    Request({
      url: '/api/meta',
      method: 'GET'
    }).then(res => resolve())
      .catch(err => {
        if (err.response.status === 401) {
          reject(err)
        }
      })
  })
}

export { PostAPIauth, VerifyAuthToken }

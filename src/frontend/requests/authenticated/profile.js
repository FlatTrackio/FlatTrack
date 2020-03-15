import Request from '@/frontend/requests/requests'

function GetAPIprofile () {
  return Request({
    url: '/api/profile',
    method: 'GET'
  })
}

export { GetAPIprofile }

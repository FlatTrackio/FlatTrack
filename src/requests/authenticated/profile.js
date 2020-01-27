import Request from '@/requests/requests'

function GetAPIprofile () {
  return Request({
    url: '/api/profile',
    method: 'GET'
  })
}

export { GetAPIprofile }

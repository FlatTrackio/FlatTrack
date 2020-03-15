import Request from '@/frontend/requests/requests'

function GetAPIflatInfo () {
  return Request({
    url: '/api/flatinfo',
    method: 'GET'
  })
}

export { GetAPIflatInfo }

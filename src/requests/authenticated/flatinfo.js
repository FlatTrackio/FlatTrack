import Request from '@/requests/requests'

function GetAPIflatInfo () {
  return Request({
    url: '/api/flatinfo',
    method: 'GET'
  })
}

export { GetAPIflatInfo }

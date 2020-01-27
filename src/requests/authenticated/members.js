import Request from '@/requests/requests'

function GetAPImemberById (id) {
  return Request({
    url: `/api/admin/members/${id}`,
    method: 'GET'
  })
}

function GetAPImembers () {
  return Request({
    url: `/api/admin/members`,
    method: 'GET'
  })
}

function PostAPImember (member) {
  return Request({
    url: `/api/admin/members`,
    method: 'POST',
    data: member
  })
}

function PutAPImember (id, member) {
  return Request({
    url: `/api/admin/members/${id}`,
    method: 'PUT',
    data: member
  })
}

function DeleteAPImember (id) {
  return Request({
    url: `/api/admin/members/${id}`,
    method: 'DELETE'
  })
}

export { GetAPImemberById, GetAPImembers, PostAPImember, PutAPImember, DeleteAPImember }

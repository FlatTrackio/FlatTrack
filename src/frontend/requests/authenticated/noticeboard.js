import Request from '@/frontend/requests/requests'

function GetAPInoticeboardEntryById (id) {
  return Request({
    url: `/api/noticeboard/${id}`,
    method: 'GET'
  })
}

function PostAPInoticeboardEntry (post) {
  return Request({
    url: `/api/noticeboard`,
    method: 'POST',
    data: post
  })
}

function PutAPInoticeboardEntry (id, post) {
  return Request({
    url: `/api/noticeboard/${id}`,
    method: 'PUT',
    data: post
  })
}

function DeleteAPInoticeboardEntry (id) {
  return Request({
    url: `/api/noticeboard/${id}`,
    method: 'DELETE'
  })
}
export { GetAPInoticeboardEntryById, PostAPInoticeboardEntry, PutAPInoticeboardEntry, DeleteAPInoticeboardEntry }

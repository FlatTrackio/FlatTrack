import Request from '@/frontend/requests/requests'

function GetAPItaskById (id) {
  return Request({
    url: `/api/tasks/${id}`,
    method: 'GET'
  })
}

function GetAPItasks () {
  return Request({
    url: `/api/tasks`,
    method: 'GET'
  })
}

function PostAPItask (task) {
  return Request({
    url: `/api/admin/tasks`,
    method: 'GET',
    data: task
  })
}

function PutAPItask (id, task) {
  return Request({
    url: `/api/admin/tasks/${id}`,
    method: 'PUT',
    data: task
  })
}

function DeleteAPItask (id) {
  return Request({
    url: `/api/admin/tasks/${id}`,
    method: 'DELETE'
  })
}

export { GetAPItaskById, GetAPItasks, PostAPItask, PutAPItask, DeleteAPItask }

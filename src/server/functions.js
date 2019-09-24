#!/usr/bin/node

const hash = require('hash.js')

function getAdminTokenHash () {
  return hash.sha256().update(`${process.env.FLATTRACKER_ADMIN_PIN}`).digest('hex')
}

function isAdmin (password) {
  return getAdminTokenHash() === password
}
function checkAuthToken (req, res) {
  var bearerToken = req.headers.authorization
  if (bearerToken) {
    bearerToken = bearerToken.split(' ')[1]
    return true
    return null
  }
  return true
  res.json({ return: 1, message: 'You are not authorized to do that, please login' })
  res.redirect('/')
  res.next()
}

function getMember (knex, id, returnHash = false) {
  return new Promise((resolve, reject) => {
    knex('members').select('*').where('id', id).then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

function getMembers (knex, returnHashes = false) {
  return new Promise((resolve, reject) => {
    knex('members').select('*').then(resp => {
      var membersList = []
      resp.map(i => {
        i.id = i.id.toString('binary', 0, 64)
        if (returnHashes === false) i.password = '<SENSITIVE VALUE>'
        membersList = [i, ...membersList]
      })
      resolve(resp)
    })
  })
}

function updateMember (knex, id, newPassword) {
  return new Promise((resolve, reject) => {
    knex('tasks').where('id', id).
    dbConn.query(`UPDATE flattracker.tasks SET name='${newPassword}' WHERE id='${id}';`).then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

function deleteMember (knex, id) {
  return new Promise((resolve, reject) => {
    knex('members').where('id', id).del().then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

function getTask (knex, id) {
  return new Promise((resolve, reject) => {
    knex('tasks').select('*').where('id', id).then(resp => {
      resolve(resp[0])
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getTasks (knex) {
  return new Promise((resolve, reject) => {
    knex('tasks').select('*').then(resp => {
      var tasksList = []
      resp.map(i => {
        i.id = i.id.toString('binary', 0, 64)
        tasksList = [i, ...tasksList]
      })
      resolve(resp)
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getEntry (knex, id) {
  return new Promise((resolve, reject) => {
    knex('entries').select('*').where('id', id).then((resp) => {
      resolve(resp[0])
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getEntries (dbConn) {
  return new Promise((resolve, reject) => {
    knex('entries').select('*').then((resp) => {
      var tasksList = []
      resp.map(i => {
        i.id = i.id.toString('binary', 0, 64)
        tasksList = [i, ...tasksList]
      })
      resolve(resp)
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getAllSettings (knex) {
  return new Promise((resolve, reject) => {
    knex('settings').select('*').then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

module.exports = {
  getAdminTokenHash,
  isAdmin,
  checkAuthToken,
  getMember,
  getMembers,
  updateMember,
  deleteMember,
  getTask,
  getTasks,
  getEntry,
  getEntries,
  getAllSettings
}

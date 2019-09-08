#!/usr/bin/node

const hash = require('hash.js')

function getAdminTokenHash () {
  return hash.sha256().update(`${process.env.FLATTRACKER_ADMIN_PIN}`).digest('hex')
}

function isAdmin (password) {
  return getAdminTokenHash() === password
}
function verifyAdminHeaderBearer (req) {

  var bearerToken = req.headers.authorization
  if (bearerToken) {
    bearerToken = bearerToken.split(' ')[1]
    return null
  }
  return isAdmin(bearerToken)
}

function getMember (dbConn, names, returnHashes = false) {
  return new Promise((resolve, reject) => {
    dbConn.query(`SELECT * FROM members WHERE id = '${names}'`).then((resp) => {
      // console.log(resp)
      if (returnHashes === false) resp[0].password = '<SENSITIVE VALUE>'
      resolve(resp[0])
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getMembers (dbConn, returnHashes = false) {
  return new Promise((resolve, reject) => {
    dbConn.query('SELECT * FROM members').then((resp) => {
      var membersList = []
      resp.map(i => {
        i.id = i.id.toString('binary', 0, 64)
        if (returnHashes === false) i.password = '<SENSITIVE VALUE>'
        membersList = [i, ...membersList]
      })
      resolve(resp)
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function updateMember (dbConn, id, newPassword) {
  return new Promise((resolve, reject) => {
    dbConn.query(`UPDATE flattracker.tasks SET name='${newPassword}' WHERE id='${id}';`).then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

function deleteMember (dbConn, id) {
  return new Promise((resolve, reject) => {
    dbConn.query(`DELETE FROM flattracker.tasks WHERE id='${id}';`).then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

function getTask (dbConn, id) {
  return new Promise((resolve, reject) => {
    dbConn.query(`SELECT * FROM tasks WHERE id = '${id}'`).then((resp) => {
      resolve(resp[0])
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getTasks (dbConn) {
  return new Promise((resolve, reject) => {
    dbConn.query('SELECT * FROM tasks').then((resp) => {
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

function getEntry (dbConn, id) {
  return new Promise((resolve, reject) => {
    dbConn.query(`SELECT * FROM entries WHERE id = '${id}'`).then((resp) => {
      resolve(resp[0])
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getEntries (dbConn) {
  return new Promise((resolve, reject) => {
    dbConn.query('SELECT * FROM entries').then((resp) => {
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

function getAllSettings (dbConn) {
  return new Promise((resolve, reject) => {
    dbConn.query('SELECT * FROM flattracker.settings;').then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

module.exports = {
  getAdminTokenHash,
  isAdmin,
  verifyAdminHeaderBearer,
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

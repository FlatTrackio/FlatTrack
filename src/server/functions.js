#!/usr/bin/node

const hash = require('hash.js')
const jwt = require('jsonwebtoken')

function verifyAuthToken (req, res, next) {
  const config = require('../../deployment/config.json')
  var bearerToken = req.headers.authorization
  if (bearerToken) {
    bearerToken = bearerToken.split(' ')[1]
    jwt.verify(bearerToken, config.system.ACCESS_TOKEN_SECRET, (err, flatmember) => {
      if (err) {
        console.log(err)
        res.json(err)
        return res.status(403).send()
      }
      req.flatmember = flatmember.flatmember
      next()
    })
  } else {
    return res.status(401).send()
  }
}

function generateToken (email, knex) {
  // create tokens for user authentication
  return new Promise((resolve, reject) => {
    getMemberProfileByEmail(knex, email).then(resp => {
      const config = require('../../deployment/config.json')
      const flatmember = resp
      const accessToken = jwt.sign({flatmember}, config.system.REFRESH_TOKEN_SECRET, { expiresIn: '2h' })
      const refreshToken = jwt.sign({flatmember}, config.system.ACCESS_TOKEN_SECRET)
      resolve({accessToken, refreshToken})
    }).catch(err => {
      reject(err)
    })
  })
}

function checkGroupForAdmin (req, res, next) {
  // middleware to only allow admin users to access certain areas of the API

  // TODO untie checking from token, instead fetch using ID 
  if (req.flatmember.group === 'admin') {
    next()
    return
  } else {
    res.json({message: 'An admin account is required for this action'}).status(401).send().end()
    return
  }
}

function generateSecret () {
  return require('crypto').randomBytes(64).toString('hex')
}

function getMember (knex, id, returnHash = false) {
  return new Promise((resolve, reject) => {
    knex('members').select('*').where('id', id).first().then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

function getMemberProfileByEmail (knex, email) {
  return new Promise((resolve, reject) => {
    knex('members').select('id', 'email', 'names', 'joinTimestamp', 'phoneNumber', 'allergies').where('email', email).first().then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

function getMembers (knex, returnHashes = false, userID) {
  return new Promise((resolve, reject) => {
    knex('members').select('*').whereNot('id', userID).then(resp => {
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
    knex('tasks').select('*').where('id', id).first().then(resp => {
      resolve(resp)
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getTasks (req, knex) {
  return new Promise((resolve, reject) => {
    knex('tasks').select('*').where('assignee', req.flatmember.id).then(resp => {
      var tasksList = []
      resolve(resp)
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getEntry (knex, id) {
  return new Promise((resolve, reject) => {
    knex('entries').select('*').where('id', id).first().then((resp) => {
      resolve(resp)
    }).catch(err => {
      // handle error
      reject(err)
    })
  })
}

function getEntries (knex) {
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

function getAllPoints (knex) {
  return new Promise((resolve, reject) => {
    knex('flatInfo').select('*').then(resp => {
      resolve(resp)
    }).catch(err => {
      reject(err)
    })
  })
}

module.exports = {
  verifyAuthToken,
  generateToken,
  generateSecret,
  getMember,
  getMemberProfileByEmail,
  getMembers,
  updateMember,
  deleteMember,
  checkGroupForAdmin,
  getTask,
  getTasks,
  getEntry,
  getEntries,
  getAllSettings,
  getAllPoints
}

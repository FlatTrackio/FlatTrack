#!/usr/bin/node

const hash = require('hash.js')
const jwt = require('jsonwebtoken')

function verifyAuthToken (req, res) {
  const config = require('../../deployment/config.json')
  var bearerToken = req.headers.authorization
  if (bearerToken) {
    bearerToken = bearerToken.split(' ')[1]
    jwt.verify(token, config.system.ACCESS_TOKEN_SECRET, (err, email) => {
      console.log(err)
      if (err) return res.sendStatus(403)
      req.email = email
      next()
    })
  }
  return true
  res.json({ return: 1, message: 'You are not authorized to do that, please login' })
  res.redirect('/api/login')
  res.next()
}

function generateAccessToken(email) {
  const config = require('../../deployment/config.json')
  return jwt.sign(email, config.system.REFRESH_TOKEN_SECRET, { expiresIn: '2h' })
}

function generateToken (req, res) {
    // Authenticate flatmember
    const config = require('../../deployment/config.json')
    const email = req.body.email
    const flatmember = {
      email
    }
    const accessToken = generateAccessToken(flatmember)
    const refreshToken = jwt.sign(flatmember, config.system.ACCESS_TOKEN_SECRET)
    res.json({ accessToken: accessToken, refreshToken: refreshToken })
}

function generateSecret () {
  return require('crypto').randomBytes(64).toString('hex')
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
  verifyAuthToken,
  generateToken,
  generateSecret,
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

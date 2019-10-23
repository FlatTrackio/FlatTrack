#!/usr/bin/node

const hash = require('hash.js')
const jwt = require('jsonwebtoken')
const fs = require('fs')
const path = require('path')
const packageJSON = require('../../package.json')
const configPath = path.resolve(path.join('.', 'deployment', 'config.json'))

function verifyAuthToken (req, res, next) {
  const config = require('../../deployment/config.json')
  var bearerToken = req.headers.authorization
  if (bearerToken) {
    bearerToken = bearerToken.split(' ')[1]
    jwt.verify(bearerToken, config.system.ACCESS_TOKEN_SECRET, (err, flatmember) => {
      if (err) {
        console.log(err)
        res.status(403)
        res.json(err).send().end()
        return
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
    }).catch(err => reject(err))
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

function getMember (knex, id) {
  return new Promise((resolve, reject) => {
    knex('members').select('*').where('id', id).first()
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getMemberProfileByEmail (knex, email) {
  return new Promise((resolve, reject) => {
    knex('members').select('id', 'email', 'names', 'joinTimestamp', 'phoneNumber', 'allergies', 'group').where('email', email).first()
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getMembers (knex, notID) {
  return new Promise((resolve, reject) => {
    knex('members').select('id', 'email', 'names', 'joinTimestamp', 'phoneNumber', 'allergies', 'group').whereNot('id', notID)
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getAllMembers (knex) {
  return new Promise((resolve, reject) => {
    knex('members').select('id', 'email', 'names', 'joinTimestamp', 'phoneNumber', 'allergies', 'group')
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function updateMember (knex, id, props) {
  // TODO hash and set password (make a standard function?)
  props = {
    email: props.email,
    phoneNumber: props.phoneNumber || null,
    password: props.password || undefined,
    allergies: props.allergies || null,
    group: props.group || 'flatmember'
  }
  return new Promise((resolve, reject) => {
    knex('members').where('id', id).update(props)
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function deleteMember (knex, id) {
  return new Promise((resolve, reject) => {
    knex('members').where('id', id).del()
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getTaskOfMember (knex, req, id) {
  return new Promise((resolve, reject) => {
    knex('tasks').select('*').where('id', id).where('assignee', req.flatmember.id).first()
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getTasksOfMember (knex, req) {
  return new Promise((resolve, reject) => {
    knex('tasks').select('*').where('assignee', req.flatmember.id)
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getTasks (knex) {
  return new Promise((resolve, reject) => {
    knex('tasks').select('*')
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getTask (knex, id) {
  return new Promise((resolve, reject) => {
    knex('tasks').select('*').where('id', id).first()
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function createTask (knex, form) {
  form = {
    id: form.id,
    name: form.name,
    description: form.description,
    location: form.location,
    rotation: form.rotation,
    frequency: form.frequency
  }
  return new Promise((resolve, reject) => {
    knex('tasks').insert(form)
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function updateTask (knex, id, props) {
  props = {
    name: props.name,
    description: props.description,
    location: props.location,
    rotation: props.rotation,
    assignee: props.assignee,
    frequency: props.frequency
  }

  if (props.rotation !== 'never') {
    props.assignee = null
  }

  return new Promise((resolve, reject) => {
    knex('tasks').where('id', id).update(props)
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function deleteTask (knex, id) {
  return new Promise((resolve, reject) => {
    knex('tasks').where('id', id).first().del()
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getEntry (knex, id) {
  return new Promise((resolve, reject) => {
    knex('entries').select('*').where('id', id).first()
      .then(resp => resolve(resp))
      .catch(err => reject(err))
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
    }).catch(err => reject(err))
  })
}

function updateEntry (knex, id, props) {
  props = {
    timestamp: props.timestamp,
    status: props.status,
    approvedBy: props.approvedBy,
    amendStatus: props.amendStatus
  }
  return new Promise((resolve, reject) => {
    knex('entries').where('id', id).update(props)
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getAllSettings (knex) {
  return new Promise((resolve, reject) => {
    knex('settings').select('*')
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function getAllPoints (knex) {
  return new Promise((resolve, reject) => {
    knex('flatInfo').select('*')
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

function updateTaskNotificationFrequency (knex, id, frequency) {
  return new Promise((resolve, reject) => {
    knex('members').where('id', id).update({ taskNotificationFrequency: frequency })
      .then(resp => resolve(resp))
      .catch(err => reject(err))
  })
}

const configJSONTemplate = {
  "system": {
    "installedVersion": packageJSON.version,
    "maintenence": false,
    "hasInitialised": false,
    "dbInstalled": false,
    "DB_ROOT_PASSWORD": process.env.DB_ROOT_PASSWORD || "",
    "DB_PASSWORD": process.env.DB_PASSWORD || "",
    "DB_DATABASE": process.env.DB_DATABASE || "",
    "DB_USER": process.env.DB_USER || "",
    "DB_HOST": process.env.DB_HOST || "",
    "DB_FLAVOR": process.env.DB_FLAVOR || "",
    "ACCESS_TOKEN_SECRET": generateSecret() || "",
    "REFRESH_TOKEN_SECRET": generateSecret() || "",
    "MAIL_SMTP_USER": process.env.MAIL_SMTP_USER || "",
    "MAIL_SMTP_PASSWORD": process.env.MAIL_SMTP_PASSWORD || "",
    "MAIL_SMTP_MODE": process.env.MAIL_SMTP_MODE || "",
    "MAIL_FROM_ADDRESS": process.env.MAIL_FROM_ADDRESS || "",
    "MAIL_DOMAIN": process.env.MAIL_DOMAIN || "",
    "MAIL_SMTP_AUTH": process.env.MAIL_SMTP_AUTH || "",
    "MAIL_SMTP_SERVER": process.env.MAIL_SMTP_SERVER || "",
    "MAIL_SMTP_PORT": process.env.MAIL_SMTP_PORT || "",
    "MAIL_SMTP_NAME": process.env.MAIL_SMTP_NAME || ""
  },
  "apps": {}
}

function doesExistConfigJSON () {
  return fs.existsSync(configPath)
}

function initConfigJSON () {
  if (! doesExistConfigJSON()) {
    if (fs.mkdirSync(path.resolve(path.join('.', 'deployment')), { recursive: true })) {
      return writeConfigJSON(configJSONTemplate)
    } else if (! doesExistConfigJSON()) {
      return writeConfigJSON(configJSONTemplate)
    } else  return false
  } else return true
}

function deinitConfigJSON () {
  if (fs.unlinkSync(configPath)) {
    return true
  } else {
    return false
  }
}

function readConfigJSON () {
  return require(configPath)
}

function writeConfigJSON (content) {
  return fs.writeFileSync(configPath, JSON.stringify(content, null, 2))
}

module.exports = {
  general: {
    entry: {
      get: getEntry,
      update: updateEntry,
      all: {
        get: getEntries
      }
    },
    member: {
      all: {
        get: getMembers
      }
    },
    verifyAuthToken,
    generateToken,
    generateSecret,
    checkGroupForAdmin,
    getTaskOfMember,
    getTasksOfMember,
    getEntry,
    getEntries,
    getAllSettings,
    getAllPoints,
    updateTaskNotificationFrequency
  },
  admin: {
    task: {
      get: getTask,
      create: createTask,
      update: updateTask,
      delete: deleteTask,
      all: {
        get: getTasks
      }
    },
    member: {
      get: getMember,
      getByEmail: getMemberProfileByEmail,
      update: updateMember,
      delete: deleteMember,
      all: {
        get: getAllMembers
      }
    },
    config: {
      exists: doesExistConfigJSON,
      init: initConfigJSON,
      deinit: deinitConfigJSON,
      read: readConfigJSON,
      write: writeConfigJSON
    }
  }
}

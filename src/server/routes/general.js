const express = require('express')
const router = express.Router()
const functions = require('../functions')
const moment = require('moment')
const uuid = require('uuid/v4')
const jwt = require('jsonwebtoken')
const passport = require('passport')
const hash = require('hash.js')
const packageJSON = require('../../../package.json')

module.exports = (knex) => {
  router.get('/', (req, res) => {
      // return list of api endpoints
      res.json({
        message: 'Hello from FlatTracker API v1',
        version: packageJSON.version,
        return: 0
      })
      res.end()
      return
    })
    
  router.route('/login')
    .post((req, res) => {
      if (req.body.email === '' || req.body.password === '') {
        res.status(401).send()
        res.end()
        return
      }
      knex.select('id', 'email', 'password').from('members').where('email', req.body.email).first().then(resp => {
        // hash sent password, check it against saved password from database
        var hashedSentPassword = hash.sha256().update(req.body.password).digest('hex')
        if (hashedSentPassword === resp.password) {
          functions.general.generateToken(req.body.email, knex).then(resp => {
            res.json(resp).status(200).end()
            return
          }).catch(err => {
            console.log(err)
            res.json({message: 'Failed generating a token'})
            res.status(403).send().end()
            return
          })
        } else {
          res.status(401).send().end()
          return
        }
      }).catch(err => {
        res.status(401).send().end()
        console.log(err)
        return
      })
    })
    
  router.route('/entry')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get all entries
      functions.general.getEntries(knex).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
  router.route('/entry/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a particular entry
      var id = req.params.id
      functions.general.entry.get(knex, id).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
    .put(functions.general.verifyAuthToken, (req, res) => {
      var id = req.params.id
      functions.general.entry.get(knex, id).then(resp => {
        if (resp.member == req.flatmember.id) {
          res.json(resp)
          res.end()
          return
        } else {
          res.status(403)
          res.end()
          return
        }
      }).catch(err => {
        res.status(400)
        res.json(err)
        res.end()
        return
      })
      var props = {
        status: req.body.status || null,
        timestamp: moment().unix() || null,
        approvedBy: req.body.approvedBy || null,
        amendStatus: req.body.amendStatus || null
      }
      functions.general.entry.update(knex, id, props).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.status(400)
        res.json(err)
        res.end()
        return
      })
    })
  
  router.route('/members')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a list of all flat members
      var memberSearch = '*'
      if (!req.query.allMembers) {
        memberSearch = req.flatmember.id
      }

      functions.general.getMembers(knex, returnHashes = false, memberSearch).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
    .post([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // add a new flat member (requires admin)
      var form = req.body
      console.log(req)
      form = {
        id: uuid(),
        names: form.names,
        password: form.password,
        email: form.email,
        group: form.group,
        phoneNumber: form.phoneNumber || null,
        allergies: form.allergies || null,
        contractAgreement: form.contractAgreement || 1,
        joinTimestamp: moment().unix(),
        memberSetPassword: form.memberSetPassword
      }

      if (form.memberSetPassword === "true") {
        form.password = "__SETME__"
      }

      // validate fields
      var regexNames = /([A-Za-z])\w+/
      switch (form) {
        case typeof form.names !== 'string' || (!regexNames.test(form.names) && form.names.length >= 100):
          res.json({ status: 1, message: 'Please enter a valid name, containing only letters' })
          res.end()
          return
  
        case (form.memberSetPassword !== true && !(typeof form.password === 'string' && form.password.length > 30)):
          res.json({ status: 1, message: 'Please enter a valid password, 30 characters max' })
          res.end()
          return
      }
      functions.general.getMember(knex, form.member).then(resp => {
        res.json({ return: 1, message: 'User already exists.' })
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })

      // hash the password
      form.password = hash.sha256().update(form.password).digest('hex')
  
      console.log("Request to create new member:")
      console.log(JSON.stringify(form))

      knex('members').insert(form).then(resp => {
        // handle user creating sucessfully
        res.json({ message: 'Added new user successfully' })
        res.end()
        return
      }).catch(err => {
        // handle error
        console.log(err)
        res.status(201)
        res.json({ return: 1, message: 'Failed to create user' })
        return
      })
    })
  router.route('/members/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a given flat member
      var id = req.params.id
      functions.general.getMember(knex, id).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
    .put([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // update a password for a given flat member (requires admin or previous password)
      var id = req.params.id
      var form = req.body
      form = {
        email: form.email,
        phoneNumber: form.phoneNumber || null,
        password: form.password,
        allergies: form.allergies || null,
        group: form.group,
      }
  
      functions.general.updateMember(knex, id, form).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        console.error(err)
        res.status(400)
        res.json({return: 1, message: 'An error occurred'})
        res.end()
        return
      })
    })
    .delete([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // delete a flat member (requires admin)    
      var id = req.params.id
      functions.general.deleteMember(knex, id).then(resp => {
        res.json(resp)
        res.end()
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
  
  router.route('/tasks')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a list of all tasks
      functions.general.getTasksOfMember(req, knex).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json({status: 1, message: err})
        res.end()
        return
      })
    })
    .post([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // add a new task (requires admin)
      var id = req.params.id
      var form = req.body
      form = {
        id: uuid(),
        name: form.name,
        description: form.description,
        location: form.location,
        importance: form.importance
      }
      // validate fields
      var regexNames = /([A-Za-z])\w+/
      switch (form) {
        case typeof form.name !== 'string' || (!regexNames.test(form.name) && form.name.length >= 100):
          res.json({ status: 1, message: 'Please enter a valid name, containing only letters' })
          res.end()
          return
  
        case typeof form.description !== 'string' || (!regexNames.test(form.description) && form.description.length >= 100):
          res.json({ status: 1, message: 'Please enter a valid description, containing only letters' })
          res.end()
          return
  
        case typeof form.location !== 'string' || (!regexNames.test(form.location) && form.location.length >= 100):
          res.json({ status: 1, message: 'Please enter a valid location, containing only letters' })
          res.end()
          return
      }
      functions.general.getTaskOfMember(knex, form.id).then(resp => {
        res.json({ return: 1, message: 'Task already exists.' })
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
  
      knex('tasks').insert(form).then((resp) => {
        console.log(resp)
        res.json({message: "Task created"})
        res.end()
        return
      }).catch(err => {
        // handle error
        console.log(err)
        res.json({message: "Task failed to create"})
        res.status(201)
        res.end()
        return
      })
    })
  router.route('/task/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a given task
      var id = req.params.id
      functions.general.getTaskOfMember(knex, id).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
    .put([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // update a given task
      var id = req.params.id
    })
    .delete([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // delete a task (requires admin)
      var id = req.params.id
    })

  router.route('/settings')
    .get((req, res) => {
      functions.general.getAllSettings(knex).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json({ message: err, return: 1 })
        res.end()
        return
      })
    })

  router.route('/settings/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      var id = req.params.id
      functions.general.getAllSettings(knex).then(resp => {  
        res.json(resp[0])
        res.end()
        return
      })
    })
    .put([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      
    })
  
  router.route('/profile')
    .get(functions.general.verifyAuthToken, (req, res, next) => {
      if (!typeof req.flatmember === 'object' || req.flatmember === '') {
        res.status(403).send()
        res.end()
        return
      } else {
        // TODO fetch by id, instead of email
        // select('id').select('names').select('email').select('phoneNumber').select('allergies')
        knex('members').select('*').where('email', req.flatmember.email).first().then(resp => {
          res.json(resp)
          res.end()
          return
        }).catch(err => {
          console.log(err)
          res.status(403).send()
          res.end()
          return
        })
      }
    })
    .post(functions.general.verifyAuthToken, (req, res, next) => {
      if (req.body.frequency) {
        switch (req.body.frequency) {
          case '0': case '1': case '2': case '3':
            break

          default:
            res.status(403).send().end()
            return
            break
        }
        functions.general.updateTaskNotificationFrequency(knex, req.flatmember.id, req.body.frequency).then(resp => {
          res.status(200).send().end()
          return
        }).catch(err => {
          console.error(err)
          res.status(403).send().end()
          return
        })
      } else {
        res.status(403).send().end()
        return
      }
    })

  router.route('/flatinfo')
    .get(functions.general.verifyAuthToken, (req, res, next) => {
      functions.general.getAllPoints(knex).then(resp => {
        res.json(resp)
        res.status(400).send().end()
        return
      }).catch(err => {
        res.status(403).send().end()
        return
      })
    })

  router.get('/meta', functions.general.verifyAuthToken, (req, res) => {
    res.json({ version: packageJSON.version })
    res.end()
    return
  })

  router.get('/health', (req, res) => {
    // get health state
  
    var health = {
      return: 0,
      healthy: undefined
    }
    knex.raw('SELECT 0;').then(conn => {
      if (conn) {
        health.healthy = true
      }
      res.json(health)
      res.end()
    }).catch(err => {
      health.healthy = false
      health.description = err
      res.json(health)
      res.end()
    })
  })
  
  if (process.env.NODE_ENV !== "production")
    router.route('/httptest')
      .all((req, res) => {
        console.log(req)
        res.json({message: "Check the output in the console"})
        res.end()
        return
      })
  return router
}

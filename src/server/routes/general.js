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
      var notID = req.flatmember.id
      functions.admin.member.all.get(knex, notID).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
  router.route('/members/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a given flat member
      var id = req.params.id
      functions.admin.member.get(knex, id).then(resp => {
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
  
  router.route('/tasks')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a list of all tasks
      functions.general.getTasksOfMember(knex, req).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.status(400)
        res.json({message: err})
        res.end()
        return
      })
    })
  router.route('/task/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a given task
      var id = req.params.id
      functions.general.getTaskOfMember(knex, req, id).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
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

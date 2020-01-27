/*
  routes/general.js

  API routes which requires Flatmember group access
*/

const express = require('express')
const router = express.Router()
const functions = require('../functions')
const moment = require('moment')
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
    return res.end()
  })

  router.route('/login')
    .post((req, res) => {
      if (req.body.email === '' || req.body.password === '') {
        res.status(401).send()
        return res.end()
      }
      knex.select('id', 'email', 'password').from('members').where('email', req.body.email).first().then(resp => {
        // hash sent password, check it against saved password from database
        var hashedSentPassword = hash.sha256().update(req.body.password).digest('hex')
        if (hashedSentPassword === resp.password) {
          functions.general.generateToken(req.body.email, knex).then(resp => {
            res.status(200)
            res.json(resp)
            return res.end()
          }).catch(err => {
            console.log(err)
            res.status(500)
            res.json({message: 'Failed generating a token'})
            res.send()
            return res.end()
          })
        } else {
          res.status(401)
          res.json({ message: 'The password provided is incorrect, unable to login' })
          return res.end
        }
      }).catch(err => {
        res.status(401)
        res.json({ message: 'Unable to find flatmember' })
        console.log(err)
        return res.end()
      })
    })

  router.route('/entry')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get all entries
      functions.general.entry.all.get(knex, req).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.status(400)
        res.json(err)
        return res.end()
      })
    })
  router.route('/entry/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a particular entry
      var id = req.params.id
      functions.general.entry.get(knex, id).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.status(400)
        res.json(err)
        return res.end()
      })
    })
    .put(functions.general.verifyAuthToken, (req, res) => {
      var id = req.params.id
      functions.general.entry.get(knex, id).then(resp => {
        if (resp.member === req.flatmember.id) {
          res.json(resp)
          return res.end()
        } else {
          res.status(500)
          res.json({ message: 'Failed updating entry' })
          return res.end()
        }
      }).catch(err => {
        res.status(400)
        console.log(err)
        res.json({ message: 'Failed updating entry' })
        return res.end()
      })
      var props = {
        status: req.body.status || null,
        timestamp: moment().unix() || null,
        approvedBy: req.body.approvedBy || null,
        amendStatus: req.body.amendStatus || null
      }
      functions.general.entry.update(knex, id, props).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.status(400)
        res.json({ message: 'Failed updating entry' })
        console.log(err)
        return res.end()
      })
    })

  router.route('/members')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a list of all flat members
      var notID = req.query.all || req.flatmember.id
      functions.general.member.all.get(knex, notID).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.status(400)
        res.json(err)
        return res.end()
      })
    })
  router.route('/members/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a given flat member
      var id = req.params.id
      functions.admin.member.get(knex, id).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(400)
        res.json({ message: 'Failed to fetch member' })
        return res.end()
      })
    })

  router.route('/tasks')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a list of all tasks
      functions.general.getTasksOfMember(knex, req).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(400)
        res.json({ message: 'Failed to fetch tasks' })
        return res.end()
      })
    })
  router.route('/task/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a given task
      var id = req.params.id
      functions.general.getTaskOfMember(knex, req, id).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(400)
        res.json({ message: 'Failed to fetch task' })
        return res.end()
      })
    })

  router.route('/noticeboard')
    .get(functions.general.verifyAuthToken, (req, res) => {
      functions.general.noticeboard.all.get(knex).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(400)
        res.json({ message: 'Failed to fetch posts' })
        return res.end()
      })
    })
    .post(functions.general.verifyAuthToken, (req, res) => {
      var form = req.body
      functions.general.noticeboard.create(knex, req, form).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(400)
        res.json({ message: 'Failed to create post' })
        return res.end()
      })
    })

  router.route('/settings')
    .get((req, res) => {
      functions.general.getAllSettings(knex).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(400)
        res.json({ message: 'Failed to fetch settings' })
        return res.end()
      })
    })

  router.route('/settings/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      var id = req.params.id

      // TODO whitelist settings to fetch from non-admin
      functions.admin.setting.get(knex, id).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(400)
        res.json({ message: 'Failed to fetch setting' })
        return res.end()
      })
    })

  router.route('/profile')
    .get(functions.general.verifyAuthToken, (req, res, next) => {
      if (!typeof req.flatmember === 'object' || req.flatmember === '') {
        res.status(403)
        res.json({ message: 'Failed to fetch profile information' })
        return res.end()
      } else {
        // TODO fetch by id, instead of email
        // select('id').select('names').select('email').select('phoneNumber').select('allergies')
        knex('members').select('*').where('email', req.flatmember.email).first().then(resp => {
          res.json(resp)
          return res.end()
        }).catch(err => {
          console.log(err)
          res.status(403)
          res.json({ message: 'Failed to fetch profile information' })
          return res.end()
        })
      }
    })
    .post(functions.general.verifyAuthToken, (req, res, next) => {
      if (req.body.frequency) {
        switch (req.body.frequency) {
          case '0': case '1': case '2': case '3':
            break

          default:
            res.status(403)
            res.json({ message: 'Failed to update profile information' })
            return res.end()
        }
        functions.general.updateTaskNotificationFrequency(knex, req.flatmember.id, req.body.frequency).then(resp => {
          res.status(200)
          res.json({ message: 'Updated profile successfully' })
          return res.end()
        }).catch(err => {
          console.error(err)
          res.status(403)
          res.json({ message: 'Failed to updated profile information' })
          return res.end()
        })
      } else {
        res.status(403)
        res.json({ message: 'Failed to updated profile information' })
        return res.end()
      }
    })

  router.route('/flatinfo')
    .get(functions.general.verifyAuthToken, (req, res, next) => {
      functions.general.getAllPoints(knex).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(403)
        res.json({ message: 'Failed fetching flat info' })
        return res.end()
      })
    })

  router.route('/features')
    .get(functions.general.verifyAuthToken, (req, res) => {
      functions.general.features.get(knex).then(resp => {
        res.json(resp)
      }).catch(err => {
        console.log(err)
        res.status(403)
        res.json({ message: 'Failed fetching features' })
        return res.end()
      })
    })

  router.get('/meta', functions.general.verifyAuthToken, (req, res) => {
    res.json({ version: packageJSON.version })
    return res.end()
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

  if (process.env.NODE_ENV !== 'production') {
    router.route('/httptest')
      .all((req, res) => {
        console.log(req)
        res.json({message: 'Check the output in the console'})
        return res.end()
      })
  }
  return router
}

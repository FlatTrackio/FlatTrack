const express = require('express')
const router = express.Router()
const mariadb = require('mariadb')
const functions = require('./functions')
const moment = require('moment')
const uuid = require('uuid/v4')
const packageJSON = require('../../package.json')
const pool = mariadb.createPool({
    host: process.env.DB_HOST,
    user: process.env.DB_USER,
    password: process.env.DB_PASSWORD,
    database: process.env.DB_DATABASE,
    connectionLimit: 5
  })
var dbConn
pool.getConnection().then(conn => {
  dbConn = conn
}).catch(err => {
  console.log(err)
})


/*const knex = require('knex')({
  client: 'mysql',
  connection: '127.0.0.1',
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_DATABASE
})*/

router.get('/', (req, res) => {
    // return list of api endpoints
    res.json({
      message: 'Hello from FlatTracker API v1',
      version: packageJSON.version,
      return: 0
    })
    res.end()
  })
  
  router.route('/entry')
    .get((req, res) => {
      // get all entries
      functions.getEntries(dbConn).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
    .post((req, res) => {
      // add a new entry
      // TODO add non-admin user auth via password
      if (!functions.verifyAdminHeaderBearer(req)) {
        res.json({ return: 1, message: 'Whoops! you need to be admin to do that.' })
        res.end()
        return
      }
      var form = req.body
      form = {
        member: form.member,
        taskName: form.taskName,
        timeSpent: form.timeSpent
      }
      // validate fields
      var regexNames = /([A-Za-z])\w+/
      functions.getMember(dbConn, form.member).then(resp => {
        console.log('Found user', resp)
      }).catch(resp => {
        res.json({ return: 1, message: 'Unable to find that user.' })
        res.end()
        return
      })
      switch (form) {
        case typeof form.taskName !== 'string' || (!regexNames.test(form.taskName) && form.taskName.length >= 100):
          res.json({ status: 1, message: 'Please enter a valid task name, containing only letters' })
          res.end()
          return
      }
      dbConn.query(`INSERT INTO flattracker.entries (id,timestamp,member,taskName) VALUES ('${uuid()}','${moment().unix()}','${form.member}','${form.taskName}');`).then((resp) => {
        console.log(resp)
        res.json({status: 0, message: "Entry added"})
        res.end()
        return
      }).catch(err => {
        // handle error
        console.log(err)
        res.json({status: 1, message: "Failed to add entry :("})
        res.end()
        return
      })
    })
  router.get('/entry/:id', (req, res) => {
    // get a particular entry
    var id = req.params.id
    functions.getEntry(dbConn, id).then(resp => {
      res.json(resp)
      res.end()
    }).catch(err => {
      res.json(err)
      res.end()
      return
    })
  })
  
  router.route('/members')
    .get((req, res) => {
      // get a list of all flat members
      functions.getMembers(dbConn).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
    .post((req, res) => {
      // add a new flat member (requires admin)
      if (!functions.verifyAdminHeaderBearer(req)) {
        res.status(201)
        res.json({ return: 1, message: 'Whoops! you need to be admin to do that.' })
        res.end()
        return
      }
      var form = req.body
      form = {
        names: form.names,
        password: form.password,
        email: form.email,
        phoneNumber: form.phoneNumber,
        allergies: form.allergies,
        contractAgreement: form.contractAgreement,
        group: form.group
      }

      // TODO implement contract agreements
      form.contractAgreement = "FALSE"

      // validate fields
      var regexNames = /([A-Za-z])\w+/
      switch (form) {
        case typeof form.names !== 'string' || (!regexNames.test(form.names) && form.names.length >= 100):
          res.json({ status: 1, message: 'Please enter a valid name, containing only letters' })
          res.end()
          return
  
        case !(typeof form.password === 'string' || typeof form.body.password === 'number'):
          res.json({ status: 1, message: 'Please enter a valid password, containing only numbers' })
          res.end()
          return
      }
      functions.getMember(dbConn, form.member).then(resp => {
        res.json({ return: 1, message: 'User already exists.' })
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
  
      console.log("Request to create new member:")
      console.log(JSON.stringify(form))

      dbConn.query(`INSERT INTO flattracker.members (id,names,email,password,joinTimestamp,phoneNumber,allergies,contractAgreement,group) VALUES ('${uuid()}','${form.names}','${form.email}','${form.password}','${moment().unix()}','${form.phoneNumber}','${form.allergies}',${form.contractAgreement},'${form.group}');`).then((resp) => {
        console.log(resp)
      }).catch(err => {
        // handle error
        console.log(err)
        res.status(201)
        res.json({ return: 1, message: 'Failed to create user' })
      })
    })
  router.route('/members/:id')
    .get((req, res) => {
      // get a given flat member
      var id = req.params.id
      functions.getMember(dbConn, id).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
    .put((req, res) => {
      // update a password for a given flat member (requires admin or previous password)
      var id = req.params.id
      var form = req.body
      form = {
        password: form.password,
        newPassword: form.newPassword
      }
  
      // get the user's row
      functions.getMember(dbConn, id, returnHashes = true).then(resp => {
        // verify the user
        if (!(functions.isAdmin(form.password) || form.password === resp.password)) {
          res.json({ return: 1, message: `${resp.names}'s old password or the Administrator password must be provided to do this.` })
          res.end()
        }
  
        // change the password
        return functions.updateMember(dbConn, id, form.newPassword)
      }).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json({return: 1, message: err})
        res.end()
        return
      })
    })
    .delete((req, res) => {
      // delete a flat member (requires admin)
      if (!functions.verifyAdminHeaderBearer(req)) {
        res.json({ return: 1, message: 'Whoops! you need to be admin to do that.' })
        res.end()
        return
      }
  
      var id = req.params.id
      functions.deleteMember(dbConn, id).then(resp => {
        res.json(resp)
        res.end()
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
  
  router.route('/tasks')
    .get((req, res) => {
      // get a list of all tasks
      functions.getTasks(dbConn).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json({status: 1, message: err})
        res.end()
        return
      })
    })
    .post((req, res) => {
      // add a new task (requires admin)
      var id = req.params.id
      if (!functions.verifyAdminHeaderBearer(req)) {
        res.json({ return: 1, message: 'Whoops! you need to be admin to do that.' })
        res.end()
        return
      }
      var form = req.body
      form = {
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
      functions.getTask(dbConn, form.id).then(resp => {
        res.json({ return: 1, message: 'Task already exists.' })
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
  
      dbConn.query(`INSERT INTO flattracker.tasks (id,name,description,location,importance,shortid) VALUES ('${uuid()}','${form.name}','${form.description}','${form.location}','${form.importance}','${shortid.generate()}');`).then((resp) => {
        console.log(resp)
      }).catch(err => {
        // handle error
        console.log(err)
      })
    })
  router.route('/task/:id')
    .get((req, res) => {
      // get a given task
      var id = req.params.id
      functions.getTask(dbConn, id).then(resp => {
        res.json(resp)
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
    .put((req, res) => {
      // update a given task
      var id = req.params.id
    })
    .delete((req, res) => {
      // delete a task (requires admin)
      var id = req.params.id
    })

  router.route('/settings')
    .get((req, res) => {
      functions.getAllSettings(dbConn).then(resp => {
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
    .get((req, res) => {
      var id = req.params.id
      functions.getAllSettings(dbConn).then(resp => {  
        res.json(resp[0])
        res.end()
        return
      })
    })
    .put((req, res) => {
      
    })
  
  router.get('/health', (req, res) => {
    // get health state
  
    var health = {
      return: 0,
      healthy: undefined
    }
    pool.getConnection().then(conn => {
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

if (process.env.NODE_ENV !== "production") router.route('/httptest')
  .all((req, res) => {
    console.log(req)
    res.json({message: "Check the output in the console"})
    res.end()
    return
  })

module.exports = router

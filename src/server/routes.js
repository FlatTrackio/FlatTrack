const express = require('express')
const router = express.Router()
const mariadb = require('mariadb')
const functions = require('./functions')
const pool = mariadb.createPool({
    host: process.env.MYSQL_HOST,
    user: process.env.MYSQL_USER,
    password: process.env.MYSQL_PASSWORD,
    database: process.env.MYSQL_DATABASE,
    connectionLimit: 5
  })
  var dbConn
  pool.getConnection().then(conn => {
    dbConn = conn
  }).catch(err => {
    console.log(err)
  })

router.get('/', (req, res) => {
    // return list of api endpoints
    res.json({
      message: 'Hello from FlatTracker API v1.0.0',
      version: '1.0.0',
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
      // TODO add non-admin user auth via pin
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
      var regexPin = /([0-9])/
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
  
        case !(typeof form.timeSpent === 'string' || typeof form.body.pin === 'number') || (!regexPin.test(form.timeSpent) && form.timeSpent.length >= 10):
          res.json({ status: 1, message: 'Please enter a valid amount of time, containing only numbers' })
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
        res.json({ return: 1, message: 'Whoops! you need to be admin to do that.' })
        res.end()
        return
      }
      var form = req.body
      form = {
        names: form.names,
        pin: form.pin
      }
      // validate fields
      var regexNames = /([A-Za-z])\w+/
      var regexPin = /([0-9])/
      switch (form) {
        case typeof form.names !== 'string' || (!regexNames.test(form.names) && form.names.length >= 100):
          res.json({ status: 1, message: 'Please enter a valid name, containing only letters' })
          res.end()
          return
  
        case !(typeof form.pin === 'string' || typeof form.body.pin === 'number') || (!regexPin.test(form.pin) && form.pin.length >= 10):
          res.json({ status: 1, message: 'Please enter a valid pin, containing only numbers' })
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
  
      dbConn.query(`INSERT INTO flattracker.members (id,names,pin,joinTimestamp) VALUES ('${uuid()}','${form.names}','${form.pin}','${moment().unix()}');`).then((resp) => {
        console.log(resp)
      }).catch(err => {
        // handle error
        console.log(err)
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
      // update a pin for a given flat member (requires admin or previous pin)
      var id = req.params.id
      var form = req.body
      form = {
        pin: form.pin,
        newPin: form.newPin
      }
  
      // get the user's row
      functions.getMember(dbConn, id, returnHashes = true).then(resp => {
        // verify the user
        if (!(functions.isAdmin(form.pin) || form.pin === resp.pin)) {
          res.json({ return: 1, message: `${resp.names}'s old pin or the Administrator pin must be provided to do this.` })
          res.end()
        }
  
        // change the pin
        return functions.updateMember(dbConn, id, form.newPin)
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
      var regexPin = /([0-9])/
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
  
        case typeof form.importance !== 'string' || (!regexPin.test(form.importance) && form.importance.length !== 1):
          res.json({ status: 1, message: 'Please enter a valid importance rating, containing only letters' })
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

    })

  router.route('/settings/:id')
    .get((req, res) => {

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

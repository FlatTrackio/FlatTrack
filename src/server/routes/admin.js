/*
  routes/admin.js

  API routes which requires Admin group access
*/

const express = require('express')
const router = express.Router()
const functions = require('../functions')
const moment = require('moment')
const uuid = require('uuid/v4')
const hash = require('hash.js')

module.exports = (knex) => {
  router.route('/tasks')
    .get([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res, next) => {
      functions.admin.task.all.get(knex).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.log(err)
        res.status(400)
        res.json({ message: 'Failed to fetch tasks' })
        return res.end()
      })
    })
    .post([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // add a new task (requires admin)
      var form = req.body
      form = {
        id: uuid(),
        name: form.name,
        description: form.description,
        location: form.location,
        assignee: form.assignee || null,
        rotation: form.rotation,
        frequency: form.frequency
      }

      if (form.rotation !== 'never') {
        form.frequency = form.rotation
      }

      // validate fields
      var regexNames = /([A-Za-z])\w+/
      switch (form) {
        case typeof form.name !== 'string' || (!regexNames.test(form.name) && form.name.length >= 100):
          res.status(400)
          res.json({ message: 'Please enter a valid name, containing only letters' })
          return res.end()

        case typeof form.description !== 'string' || (!regexNames.test(form.description) && form.description.length >= 100):
          res.status(400)
          res.json({ message: 'Please enter a valid description, containing only letters' })
          return res.end()

        case typeof form.location !== 'string' || (!regexNames.test(form.location) && form.location.length >= 100):
          res.status(400)
          res.json({ message: 'Please enter a valid location, containing only letters' })
          return res.end()
      }

      functions.admin.task.create(knex, form).then((resp) => {
        console.log(resp)
        res.json({ message: 'Task created' })
        return res.end()
      }).catch(err => {
        // handle error
        console.log(err)
        res.json({ message: 'Task failed to create' })
        res.status(201)
        return res.end()
      })
    })
  router.route('/task/:id')
    .get([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // get a given task
      var id = req.params.id

      functions.admin.task.get(knex, id).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.json(err)
        return res.end()
      })
    })
    .put([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // update a given task
      var id = req.params.id

      var form = req.body
      form = {
        description: form.description,
        location: form.location,
        disabled: form.disabled || null,
        rotation: form.rotation,
        assignee: form.assignee || null
      }

      if (form.rotation !== 'never') {
        form.frequency = form.rotation
      }

      if (form.rotation === 'never' && (typeof form.assignee === 'undefined' || form.assignee === null)) {
        res.status(400)
        res.json({ message: 'There must be an assignee for a task, if it never rotates' })
        res.end()
        return
      }

      functions.admin.task.update(knex, id, form).then((resp) => {
        console.log(resp)
        res.json({ message: 'Task updated' })
        return res.end()
      }).catch(err => {
        // handle error
        console.log(err)
        res.status(201)
        res.json({ message: 'Failed to update task' })
        return res.end()
      })
    })
    .delete([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // delete a task (requires admin)
      var id = req.params.id

      functions.admin.task.delete(knex, id).then((resp) => {
        console.log(resp)
        res.json({ message: 'Task deleted' })
        return res.end()
      }).catch(err => {
        // handle error
        console.log(err)
        res.json({ message: 'Failed to delete task' })
        res.status(201)
        return res.end()
      })
    })

  router.route('/members')
    .get([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // get a list of all flat members
      functions.admin.member.all.get(knex).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.json(err)
        return res.end()
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

      if (form.memberSetPassword === 'true') {
        form.password = '__SETME__'
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
      functions.admin.member.get(knex, form.member).then(resp => {
        res.json({ return: 1, message: 'User already exists.' })
        return res.end()
      }).catch(err => {
        res.json(err)
        return res.end()
      })

      // hash the password
      form.password = hash.sha256().update(form.password).digest('hex')

      console.log('Request to create new member:')
      console.log(JSON.stringify(form))

      knex('members').insert(form).then(resp => {
        // handle user creating sucessfully
        res.json({ message: 'Added new user successfully' })
        return res.end()
      }).catch(err => {
        // handle error
        console.log(err)
        res.status(201)
        res.json({ return: 1, message: 'Failed to create user' })
        return res.end()
      })
    })
  router.route('/members/:id')
    .get([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // get a given flat member
      var id = req.params.id
      functions.admin.member.get(knex, id).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.json(err)
        return res.end()
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
        group: form.group
      }

      functions.admin.member.update(knex, id, form).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        console.error(err)
        res.status(400)
        res.json({return: 1, message: 'An error occurred'})
        return res.end()
      })
    })
    .delete([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // delete a flat member (requires admin)
      var id = req.params.id
      functions.admin.member.delete(knex, id).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.json(err)
        return res.end()
      })
    })

  router.route('/settings/:id')
    .get(functions.general.verifyAuthToken, (req, res) => {
      var id = req.params.id
      functions.admin.setting.all.get(knex, id).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.status(400)
        res.json(err)
        return res.end()
      })
    })
    .put([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // update a setting
    })

  router.route('/entry')
    .get([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // get all entries
      functions.admin.entry.all.get(knex, req).then(resp => {
        res.json(resp)
        return res.end()
      }).catch(err => {
        res.status(400)
        res.json(err)
        return res.end()
      })
    })

  router.get(/(.*)/, (req, res) => {
    res.status(404)
    res.json({ message: 'Unknown route' })
    return res.end()
  })

  return router
}

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
  router.route('/tasks')
    .get([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res, next) => {
      functions.admin.task.all.get(knex).then(resp => {
        res.json(resp).end()
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
          res.end()
          return

        case typeof form.description !== 'string' || (!regexNames.test(form.description) && form.description.length >= 100):
          res.status(400)
          res.json({ message: 'Please enter a valid description, containing only letters' })
          res.end()
          return

        case typeof form.location !== 'string' || (!regexNames.test(form.location) && form.location.length >= 100):
          res.status(400)
          res.json({ message: 'Please enter a valid location, containing only letters' })
          res.end()
          return
      }

      functions.admin.task.create(knex, form).then((resp) => {
        console.log(resp)
        res.json({ message: 'Task created' })
        res.end()
        return
      }).catch(err => {
        // handle error
        console.log(err)
        res.json({ message: 'Task failed to create' })
        res.status(201)
        res.end()
        return
      })
    })
  router.route('/task/:id')
    .get([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // get a given task
      var id = req.params.id

      functions.admin.task.get(knex, id).then(resp => {
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
        res.end()
        return
      }).catch(err => {
        // handle error
        console.log(err)
        res.status(201)
        res.json({ message: 'Failed to update task' })
        res.end()
        return
      })
    })
    .delete([functions.general.verifyAuthToken, functions.general.checkGroupForAdmin], (req, res) => {
      // delete a task (requires admin)
      var id = req.params.id

      functions.admin.task.delete(knex, id).then((resp) => {
        console.log(resp)
        res.json({ message: 'Task deleted' })
        res.end()
        return
      }).catch(err => {
        // handle error
        console.log(err)
        res.json({ message: 'Failed to delete task' })
        res.status(201)
        res.end()
        return
      })
    })

    router.route('/members')
    .get(functions.general.verifyAuthToken, (req, res) => {
      // get a list of all flat members
      functions.admin.member.all.get(knex).then(resp => {
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
        res.end()
        return
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })

      // hash the password
      form.password = hash.sha256().update(form.password).digest('hex')
  
      console.log('Request to create new member:')
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
      functions.admin.member.get(knex, id).then(resp => {
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
        group: form.group
      }
  
      functions.admin.member.update(knex, id, form).then(resp => {
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
      functions.admin.member.delete(knex, id).then(resp => {
        res.json(resp)
        res.end()
      }).catch(err => {
        res.json(err)
        res.end()
        return
      })
    })
  
  return router
}
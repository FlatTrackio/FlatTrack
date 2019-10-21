/* eslint-disable no-useless-return */
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
        rotation: form.rotation
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
  return router
}
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
            functions.admin.getTasks().then(resp => {
                res.json(resp).end()
            })
        })
    return router
}
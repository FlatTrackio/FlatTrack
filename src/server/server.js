#!/usr/bin/node

/*
  server.js

  backend server for FlatTrack.
  connects to database,
  starts webserver
*/

require('dotenv').config({ path: `./.dev.env` })
const os = require('os')
if (os.userInfo().uid < 1000) {
  console.log(`[error] flattrack was run as '${os.userInfo().username}' must not be run as a system user`)
  if (process.env.NODE_ENV === 'production') process.exit(1)
}

const express = require('express')
const app = express()
const bodyParser = require('body-parser')
const path = require('path')
const morgan = require('morgan')
const knex = require('knex')({
  client: 'mysql',
  connection: {
    host: process.env.DB_HOST,
    user: process.env.DB_USER,
    password: process.env.DB_PASSWORD,
    database: process.env.DB_DATABASE
  },
  pool: { min: 0, max: 7 }
})
knex.raw('select 0;').catch(err => {
  console.log('No database connection found.')
  console.log(err)
  process.exit(1)
})
require('./background-tasks')(knex)
var serverObject

require('./init.js')(knex)

// development port is 8080
var port = process.env.APP_PORT || 8080

if (process.env.NODE_ENV !== 'production') {
  app.use((req, res, next) => {
    res.append('Access-Control-Allow-Origin', ['*'])
    res.append('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE')
    res.append('Access-Control-Allow-Headers', 'Content-Type')
    next()
  })
}

app.use(bodyParser.urlencoded({
  extended: true
}))
app.use(express.json())
app.use(morgan('combined'))

var routesFlatmember = require('./routes/general')(knex)
app.use('/api', routesFlatmember)

var routesAdmin = require('./routes/admin')(knex)
app.use('/api/admin', routesAdmin)

// Sends static files from the public path directory
app.use(express.static(path.join(__dirname, '..', '..', 'dist')))

app.get('/#', (req, res) => {
  return res.redirect('/')
})

app.get(/(.*)/, (req, res) => {
  res.status(404)
  return res.redirect('/#/unknown-page')
})

function start () {
  return app.listen(port, () => {
    console.log(`Running on port ${port}`)
  })
}

// make requirable, so it can be started and stopped programatically
if (require.main !== module) {
  module.exports = {
    start: () => {
      serverObject = start()
    },
    stop: () => {
      serverObject.close()
    }
  }
} else start()

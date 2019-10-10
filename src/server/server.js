#!/usr/bin/node

require('dotenv').config({ path: `./.dev.env` })

const express = require('express')
const app = express()
const bodyParser = require('body-parser')
const path = require('path')
const os = require('os')
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

if (os.userInfo().uid < 1000) {
  console.log('[error] ftctl must not be run as a system user')
  process.exit(1)
}

knex.raw('select 0;').catch(err => {
  console.log("No database connection found.")
  process.exit(1)
})

require('./init.js')(knex)

// development port is 8080
var port = process.env.APP_PORT || 8080

if (process.env.NODE_ENV !== "production") app.use((req, res, next) => {
  res.append('Access-Control-Allow-Origin', ['*'])
  res.append('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE')
  res.append('Access-Control-Allow-Headers', 'Content-Type')
  next()
})

app.use(bodyParser.urlencoded({
  extended: true
}))
app.use(express.json())
app.use(morgan('combined'))

var routes = require('./routes')(knex)
app.use('/api', routes)

// Sends static files from the public path directory
app.use(express.static(path.join(__dirname, '..', '..', 'dist')))

app.get('/#', (req, res) => {
    res.redirect('/')
})

app.get(/(.*)/, (req, res) => {
    res.redirect('/#/unknown-page')
})

function start () {
  return app.listen(port, () => {
      console.log(`Running on port ${port}`)
  })
}

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

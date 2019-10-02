#!/usr/bin/node

require('dotenv').config({ path: `./.dev.env` })

const express = require('express')
const app = express()
const bodyParser = require('body-parser')
const path = require('path')
const fs = require('fs')
const morgan = require('morgan')
const packageJSON = require('../../package.json')
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
  console.log("No database connection found.")
  process.exit(1)
})

require('./init.js')(knex)
const config = require('../../deployment/config.json')

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

app.get(/(.*)/, (req, res) => {
  res.redirect('/#/unknown-page')
})

// start service
app.listen(port, () => {
  console.log(`Running on ${port}`)
})

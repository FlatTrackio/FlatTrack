#!/usr/bin/node

const express = require('express')
const app = express()
const bodyParser = require('body-parser')
const path = require('path')

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

var routes = require('./routes')
app.use('/api', routes)

// Sends static files  from the public path directory
app.use(express.static(path.join(__dirname, '..', '..', 'dist')))

// start service
app.listen(port, () => {
  console.log(`Running on ${port}`)
})

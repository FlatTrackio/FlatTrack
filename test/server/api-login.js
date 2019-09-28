#!/usr/bin/node

// this is a test to create an entry via the api.
// assumes that the database is already provisioned with correct tables and that the API is active

var request = require('request-promise')
var hash = require('hash.js')
// development port is 8080
var port = process.env.APP_PORT || 8080
var protocol = process.env.NODE_ENV === 'production' ? 'https' : 'http'
var host = 'localhost'
var uri = `${protocol}://${host}:${port}/api/login`

console.log(`Making request to '${uri}'`)

var form = {
  email: 'alduas@example.com',
  password: '1234567890'
}

request({
  method: 'POST',
  headers: {
    'content-type': 'application/json',
  },
  uri,
  form
}).then((resp) => {
  console.log(resp)
  process.exit()
}).catch((err) => {
  console.log({ err })
  process.exit()
})

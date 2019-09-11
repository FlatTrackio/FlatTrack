#!/usr/bin/node

// this is a test to create a task via the api.
// assumes that the database is already provisioned with correct tables and that the API is active

var request = require('request-promise')
var hash = require('hash.js')
// development port is 8080
var port = process.env.APP_PORT || 8080
var protocol = process.env.NODE_ENV === 'production' ? 'https' : 'http'
var host = 'localhost'
var uri = `${protocol}://${host}:${port}/api/task`

console.log(`Making request to '${uri}'`)

var pinUnhashed = Math.floor(Math.random() * (999999 - 100000 + 1) + 100000) + Math.floor(Math.random() * (999999 - 100000 + 1) + 100000)
var testInsertData = {
  name: `${['Dishes', 'Garbage', 'Washing', 'Shopping', 'Dusting'][Math.floor(Math.random() * 5)]}`,
  description: `Test description #${[Math.floor(Math.random() * 100)]}`,
  location: `Test location #${[Math.floor(Math.random() * 100)]}`,
  importance: 3
}

var adminTokenPin = hash.sha256().update(`${process.env.FLATTRACKER_ADMIN_PIN}`).digest('hex')
request({
  method: 'POST',
  headers: {
    'content-type': 'application/json',
    Authorization: `bearer ${adminTokenPin}`
  },
  uri,
  form: testInsertData
}).then((resp) => {
  console.log(resp)
  console.log(JSON.stringify({ testInsertData }, null, 4), { pinUnhashed })
  console.log(resp.status, ':: Request posted.')
  process.exit()
}).catch((err) => {
  console.log({ err })
  process.exit()
})

#!/usr/bin/node

// this is a test to create a flat member via the api.
// assumes that the database is already provisioned with correct tables and that the API is active

var request = require('request-promise')
var hash = require('hash.js')
// development port is 3000
var port = process.env.APP_PORT || 3000
var protocol = process.env.NODE_ENV === 'production' ? 'https' : 'http'
var host = 'localhost'
var uri = `${protocol}://${host}:${port}/members`

console.log(`Making request to '${uri}'`)

var pinUnhashed = Math.floor(Math.random() * (999999 - 100000 + 1) + 100000) + Math.floor(Math.random() * (999999 - 100000 + 1) + 100000)
var testInsertData = {
  names: `${['John', 'Mary', 'Phillip', 'Jess', 'Kyle', 'Jack'][Math.floor(Math.random() * 5)]} ${['Doe', 'Taylor', 'Smith', 'Richards', 'Miles', 'Jackson'][Math.floor(Math.random() * 5)]}`,
  password: hash.sha256().update(`${pinUnhashed}`).digest('hex')
}

var adminTokenPin = hash.sha256().update(`${process.env.FLATTRACKER_ADMIN_PIN}`).digest('hex')
request({
  method: 'DELETE',
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

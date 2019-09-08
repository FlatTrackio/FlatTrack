#!/usr/bin/node

// this is a test to create a flat member via the api.
// assumes that the database is already provisioned with correct tables and that the API is active

var request = require('request-promise')
var hash = require('hash.js')
// development port is 3000
var port = process.env.APP_PORT || 3000
var protocol = process.env.NODE_ENV === 'production' ? 'https' : 'http'
var host = 'localhost'

var pinUnhashed = Math.floor(Math.random() * (999999 - 100000 + 1) + 100000) + Math.floor(Math.random() * (999999 - 100000 + 1) + 100000)
var testInsertData = {
    id: "a798f50a-dfef-475c-958f-86b87c915a96",
    password: "b922f606c1545ac977887d338c18a7e42e7838abc9cc5b969ff4e97602acefbb",
    newPassword: hash.sha256().update(`1234567890`).digest('hex')
}

var uri = `${protocol}://${host}:${port}/members/${testInsertData.id}`

console.log(`Making request to '${uri}'`)

var adminTokenPin = hash.sha256().update(`${process.env.FLATTRACKER_ADMIN_PIN}`).digest('hex')
request({
  method: 'PUT',
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

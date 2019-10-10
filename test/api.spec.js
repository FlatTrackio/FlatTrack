const axios = require('axios')
const assert = require('assert')
const flattrack = require('../src/server/server.js')
const functions = require('../src/server/functions.js')
const packageJSON = require('../package.json')

beforeEach('initialise site\'s config file', function() {
  if (functions.admin.config.exists()) {
    functions.admin.config.deinit()
    functions.admin.config.init()
  }
  flattrack.start()
})

afterEach('deinitialise site\'s config file', function () {
  if (functions.admin.config.exists()) {
    functions.admin.config.deinit() 
  }
  flattrack.stop()
})

describe('general', function () {
  /*
    Testname: Server spawn
    Description: Should start server
	*/
  it('should start server', function () {
    return new Promise((resolve, reject) => {
      axios.get('http://localhost:8080/api').then(resp => {
        assert.equal(resp.data.message, 'Hello from FlatTracker API v1', 'Should return greeting')
        assert.equal(resp.data.version, packageJSON.version, 'Should contain matching version')
        resolve()
      }).catch(err => {
        assert.equal(err, null, 'Should not return error')
        reject(err)
      })
    }).catch(err => {
      assert.equal(err, null, 'Should not return error')
      reject(err)
    })
  })
})

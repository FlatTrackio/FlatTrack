const fs = require('fs')
const path = require('path')
const packageJSON = require('../../package.json')

function initialisedConfigJSON () {
  const functions = require('./functions')
  console.log('- configurations')
  functions.admin.config.init()
}

function triggerDBInitialisation (knex) {
  knex.raw('show databases;').then(resp => {
    var databases = []
    resp[0].map(i => {
      databases = [...databases, i.Database]
    })
    if (databases.includes(process.env.DB_DATABASE)) {
      console.log('- database')
      // if the version is different or tables don't exist
      require(`./migrations/db-${packageJSON.version}`).up(knex)
    } else {
      console.error(`Database ${process.env.DB_DATABASE} doesn't exist`)
      process.exit(1)
    }
  })
}

module.exports = (knex) => {    
  if (!fs.existsSync(path.resolve(path.join('.', 'deployment', 'config.json')))) {
    console.log('Initializing:')
    initialisedConfigJSON()
    triggerDBInitialisation(knex)
    return
  }
    
  const configJSON = require(path.resolve(path.join('.', 'deployment', 'config.json')))
  if (!configJSON.system.dbInstalled || packageJSON.version !== configJSON.system.installedVersion) {
    console.log('Initializing:')
    triggerDBInitialisation(knex)
  }
}
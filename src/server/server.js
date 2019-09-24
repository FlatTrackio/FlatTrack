#!/usr/bin/node

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

// initialise settings
if (! fs.existsSync(path.resolve(path.join('.', 'deployment')))) {
  fs.mkdirSync(path.resolve(path.join('.', 'deployment')))
}

if (! fs.existsSync(path.resolve(path.join('.', 'deployment', 'config.json')))) {
  fs.writeFileSync(
    path.resolve(path.join('.', 'deployment', 'config.json')),
    JSON.stringify({
      "system": {
        "installedVersion": packageJSON.version,
        "DB_ROOT_PASSWORD": process.env.DB_ROOT_PASSWORD || "",
        "DB_PASSWORD": process.env.DB_PASSWORD || "",
        "DB_DATABASE": process.env.DB_DATABASE || "",
        "DB_USER": process.env.DB_USER || "",
        "DB_HOST": process.env.DB_HOST || "",
        "DB_FLAVOR": process.env.DB_FLAVOR || ""
      },
      "apps": {}
    }, null, 4)
  )
}

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

app.use(morgan('combined'))

// if the version is different or tables don't exist
if (config.installedVersion != packageJSON.version) {
  // require(`./migrations/db-${packageJSON.version}`)(knex)
  // knex.schema.hasTable(['settings', 'members', 'groups'])
  console.log("DB MIGRATION")
}

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

/*
  run background tasks,
  such as entry creation,
  task reminders
*/

const schedule = require('node-schedule')
const nodemailer = require('nodemailer')
const moment = require('moment')
const functions = require('./functions')

module.exports = (knex) => {
  function findAndCreateEntries (knex) {
    /*
      fetch tasks,
      find the latest entires for each (if available),
      if the entry has met it's completedBy date,
      create a new entry
    */
    return new Promise((resolve, reject) => {
      Promise.all([
        functions.admin.task.all.get(knex),
        functions.admin.entry.all.get(knex),
        functions.admin.member.all.get(knex)
      ]).then(resp => {
        var tasks = resp[0]
        var entries = resp[1]
        var members = resp[2]
        var currentTime = moment().unix()

        tasks.map(task => {
          // create list of entries to create
          var entriesOfTask = entries.filter((entry) => {
            // if the entry's taskID matches the task's ID and it has gone past the completeBy date
            return entry.taskID === task.id && currentTime >= entry.completeBy
          })
          if (entriesOfTask.length === 0) return
          var assignee
          members.map((member) => {
            (member.id !== task.assignee || member.id !== task.assigneeLast)
            if (members.length <= 2) {
            } else {
            }
          })
          // functions.admin.entry.create(knex, task.id, member.id)
          functions.admin.task.update(knex, task.id, {assignee, assigneeLast: task.assignee})
        })

        // console.log({tasks, entries, members})
      }).catch(err => {
        console.log(err)
        return reject(err)
      })
    })
  }

  /*
    fetch entries,
    find tasks which haven't been completed yet,
    get the user which the entry is assigned to,
    send a reminder email according to the assigned user's frequency
  */
  function sendTaskReminders (knex) {
    return new Promise((resolve, reject) => {
      functions.admin.entry.all.get(knex).then(resp => {
        var items = resp.filter(item => {
          var currentTime = moment().unix()
          if (currentTime <= item.completeBy) {
            return item
          }
        })
        return items
      }).then(resp => {})
    })
  }

  // execute jobs
  findAndCreateEntries(knex)
  return
  return schedule.scheduleJob('*/15 * * * *', () => {
    findAndCreateEntries(knex)
    sendTaskReminders(knex)
  })
}

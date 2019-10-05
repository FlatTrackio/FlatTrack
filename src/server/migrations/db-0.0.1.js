// migrations for version 0.0.1

const uuid = require('uuid/v4')

module.exports = (knex) => {
    // create tables
    knex.schema.createTable('entries', (table) => {
        table.string('id', 37)
        table.string('timestamp', 100)
        table.string('member', 100)
        table.string('taskName', 100)
        table.foreign('approvedBy').references('id').inTable('members')
        table.string('amendStatus', 100)
    })

    knex.schema.createTable('members', (table) => {
        table.string('id', 37)
        table.string('names', 100)
        table.string('email', 100)
        table.string('password', 100)
        table.string('joinTimestamp', 100)
        table.string('phoneNumber', 100)
        table.string('allergies', 100)
        table.boolean('contractAgreement')
        table.boolean('disabled')
        table.string('group', 100)
        table.string('taskNotificationFrequency', 10)
    })

    knex.schema.createTable('task', (table) => {
        table.string('id', 37)
        table.string('name', 100)
        table.string('description', 100)
        table.string('location', 100)
        table.boolean('disabled')
        table.foreign('assignee').references('id').inTable('members')
        table.foreign('assigneeLast').references('id').inTable('members')
        table.string('rotation', 100)
    })

    knex.schema.createTable('noticeboard', (table) => {
        table.string('id', 37)
        table.string('title', 100)
        table.string('message', 100)
        table.foreign('author').references('id').inTable('members')
        table.string('timestamp', 100)
    })

    knex.schema.createTable('recipes', (table) => {
        table.string('id', 37)
        table.string('name', 100)
        table.string('comment', 100)
        table.foreign('addedBy').references('id').inTable('members')
        table.string('preparationTime', 100)
        table.string('timestamp', 100)
        table.string('steps', 500)
        table.string('countryOfOrigin', 100)
    })

    knex.schema.createTable('shoppinglist', (table) => {
        table.string('id', 37)
        table.string('name', 100)
        table.string('price', 100)
        table.string('comment', 500)
        table.integer('week')
        table.string('timestamp', 100)
        table.foreign('addedBy').references('id').inTable('members')
        table.string('standard', 10)
        table.boolean('obtained')
    })

    knex.schema.createTable('features', (table) => {
        table.string('id', 37)
        table.string('name', 100)
        table.boolean('enabled')
    })

    knex.schema.createTable('flatInfo', (table) => {
        table.string('id', 37)
        table.string('line', 100)
        table.string('subPointOf', 37)
    })

    knex.schema.createTable('settings', (table) => {
        table.string('id', 37)
        table.string('name', 100)
        table.string('value', 500)
    })

    knex.schema.createTable('highfives', (table) => {
        table.string('id', 37)
        table.string('timestamp', 100)
        table.string('message', 200)
        table.foreign('addedBy').references('id').inTable('members')
    })

    knex.schema.createTable('groups', (table) => {
        table.string('id', 37)
        table.string('name', 100)
    })

    knex.schema.createTable('recipesTags', (table) => {
        table.string('id', 37)
        table.foreign('recipeID').references('id').inTable('recipes')
        table.string('name', 100)
    })

    knex.schema.createTable('recipeIngredients', (table) => {
        table.string('id', 37)
        table.foreign('recipeID').references('id').inTable('recipes')
        table.string('ingredient', 100)
    })

    ['flatmember', 'admin', 'approver'].map(name => {
        if (!knex('groups').where('name', name)) {
            knex('groups').insert({id: uuid(), name: name})
        }
    })
}
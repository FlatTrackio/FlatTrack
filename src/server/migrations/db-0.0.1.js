// migrations for version 0.0.1

const uuid = require('uuid/v4')

module.exports.up = (knex) => {
    // create tables
    return knex.schema
        .createTable('entries', function(table) {
            console.log("ENTRIES")
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('timestamp', 100).notNullable()
            table.string('timestampAssign', 100)
            table.string('member', 100).notNullable()
            table.foreign('taskID').references('id').inTable('tasks').notNullable()
            table.string('status', 100)
            table.foreign('approvedBy').references('id').inTable('members')
            table.string('amendStatus', 100)
            table.collate('latin1_general_cs')
        })

        .createTable('members', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('names', 100).notNullable()
            table.string('email', 100).notNullable()
            table.string('password', 100).notNullable()
            table.string('joinTimestamp', 100).notNullable()
            table.string('phoneNumber', 100)
            table.string('allergies', 100)
            table.boolean('contractAgreement')
            table.boolean('disabled')
            table.string('group', 100)
            table.string('taskNotificationFrequency', 10)
            table.collate('latin1_general_cs')
        })

        .createTable('tasks', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('name', 100).notNullable()
            table.string('description', 100).notNullable()
            table.string('location', 100).notNullable()
            table.string('rotation', 100).notNullable()
            table.boolean('disabled')
            table.foreign('assignee').references('id').inTable('members')
            table.foreign('assigneeLast').references('id').inTable('members')
            table.collate('latin1_general_cs')
        })

        .createTable('noticeboard', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('title', 100).notNullable()
            table.string('message', 100).notNullable()
            table.foreign('author').references('id').inTable('members').notNullable()
            table.string('timestamp', 100).notNullable()
            table.collate('latin1_general_cs')
        })

        .createTable('recipes', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('name', 100).notNullable()
            table.string('comment', 100).notNullable()
            table.foreign('addedBy').references('id').inTable('members').notNullable()
            table.string('preparationTime', 100).notNullable()
            table.string('timestamp', 100).notNullable()
            table.string('steps', 500).notNullable()
            table.string('countryOfOrigin', 100).notNullable()
            table.collate('latin1_general_cs')
        })

        .createTable('shoppinglist', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('name', 100).notNullable()
            table.string('price', 100)
            table.string('comment', 500)
            table.integer('week')
            table.string('timestamp', 100)
            table.foreign('addedBy').references('id').inTable('members').notNullable()
            table.string('standard', 10)
            table.boolean('obtained')
            table.collate('latin1_general_cs')
        })

        .createTable('features', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('name', 100).notNullable()
            table.boolean('enabled').notNullable()
            table.collate('latin1_general_cs')
        })

        .createTable('flatInfo', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('line', 100).notNullable()
            table.string('subPointOf', 37)
            table.collate('latin1_general_cs')
        })

        .createTable('settings', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('name', 100).notNullable()
            table.string('value', 500).notNullable()
            table.collate('latin1_general_cs')
        })

        .createTable('highfives', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('timestamp', 100).notNullable()
            table.string('message', 200).notNullable()
            table.foreign('addedBy').references('id').inTable('members').notNullable()
            table.collate('latin1_general_cs')
        })

        .createTable('groups', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.string('name', 100).notNullable()
            table.collate('latin1_general_cs')
        })

        .createTable('recipesTags', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.foreign('recipeID').references('id').inTable('recipes').notNullable()
            table.string('name', 100).notNullable()
            table.collate('latin1_general_cs')
        })

        .createTable('recipeIngredients', (table) => {
            table.charset('utf8mb4')
            table.string('id', 37).notNullable()
            table.foreign('recipeID').references('id').inTable('recipes').notNullable()
            table.string('ingredient', 100).notNullable()
            table.collate('latin1_general_cs')
        })
}

module.exports.down = (knex) => {
    return knex.schema
        .dropTableIfExists('entries')
        .dropTableIfExists('members')
        .dropTableIfExists('tasks')
        .dropTableIfExists('noticeboard')
        .dropTableIfExists('shoppinglist')
        .dropTableIfExists('features')
        .dropTableIfExists('flatInfo')
        .dropTableIfExists('settings')
        .dropTableIfExists('highfives')
        .dropTableIfExists('groups')
        .dropTableIfExists('recipesTags')
        .dropTableIfExists('recipeIngredients')
}

module.exports.load = (knex) => {
    const groups = ['flatmember', 'admin', 'approver']
    groups.map(name => {
        if (!knex('groups').where('name', name)) {
            knex('groups').insert({id: uuid(), name: name})
        }
    })
}
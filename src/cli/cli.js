#!/usr/bin/node

const os = require('os')
const commander = require('commander')
const prog = new commander.Command()
const packageJSON = require('../../package.json')
const functions = require('../server/functions')

prog.name('ftctl')

prog.version(packageJSON.version)
prog.parse(process.argv)

if (os.userInfo().uid < 1000) {
    console.log('[error] ftctl must not be run as a system user')
    process.exit(1)
}

prog.command('config [area]', 'set system and app settings')
    .action((area) => {
        console.log("CONFIG", area)
        switch(area) {
            case 'system':
                break
            
            case 'app':
                break
            
            default:
                console.log('please state app or system')
                break
        }
    })

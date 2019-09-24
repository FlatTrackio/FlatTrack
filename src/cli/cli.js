#!/usr/bin/node

const commander = require('commander')
const prog = new commander.Command()
const packageJSON = require('../../package.json')

prog.name('ftctl')

prog.version(packageJSON.version)
prog.parse(process.argv)

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

const fs = require('fs')
const path = require('path')
const packageJSON = require('../../package.json')

module.exports = (knex) => {
    console.log("Initializing:")

    if (! fs.existsSync(path.resolve(path.join('.', 'deployment', 'config.json')))) {
        const functions = require('./functions')
        console.log("- configurations")
        fs.writeFileSync(
            path.resolve(path.join('.', 'deployment', 'config.json')),
            JSON.stringify({
            "system": {
                "installedVersion": packageJSON.version,
                "maintenence": false,
                "hasInitialised": true,
                "DB_ROOT_PASSWORD": process.env.DB_ROOT_PASSWORD || "",
                "DB_PASSWORD": process.env.DB_PASSWORD || "",
                "DB_DATABASE": process.env.DB_DATABASE || "",
                "DB_USER": process.env.DB_USER || "",
                "DB_HOST": process.env.DB_HOST || "",
                "DB_FLAVOR": process.env.DB_FLAVOR || "",
                "ACCESS_TOKEN_SECRET": functions.generateSecret() || "",
                "REFRESH_TOKEN_SECRET": functions.generateSecret() || "",
                "MAIL_SMTP_USER": process.env.MAIL_SMTP_USER || "",
                "MAIL_SMTP_PASSWORD": process.env.MAIL_SMTP_PASSWORD || "",
                "MAIL_SMTP_MODE": process.env.MAIL_SMTP_MODE || "",
                "MAIL_FROM_ADDRESS": process.env.MAIL_FROM_ADDRESS || "",
                "MAIL_DOMAIN": process.env.MAIL_DOMAIN || "",
                "MAIL_SMTP_AUTH": process.env.MAIL_SMTP_AUTH || "",
                "MAIL_SMTP_SERVER": process.env.MAIL_SMTP_SERVER || "",
                "MAIL_SMTP_PORT": process.env.MAIL_SMTP_PORT || "",
                "MAIL_SMTP_NAME": process.env.MAIL_SMTP_NAME || ""
            },
            "apps": {}
            }, null, 4)
        )
    }

    if (true) {
        console.log("- database")
        // if the version is different or tables don't exist
        // require(`./migrations/db-${packageJSON.version}`)(knex)
        // knex.schema.hasTable(['settings', 'members', 'groups'])
    }

    console.log("")
}
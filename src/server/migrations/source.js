const path = require('path')
const fs = require('fs')
const semver = require('semver')
const baseMigrationPath = path.join(process.cwd(), 'src', 'server', 'migrations', 'sql')

module.exports = {
  async getMigrations () {
    const migrationFiles = fs.readdirSync(baseMigrationPath)
    return migrationFiles.sort(semver.compare).map(file => ({
      file,
      directory: baseMigrationPath
    }))
  },

  getMigrationName (migration) {
    return migration.file
  },

  getMigration (migration) {
    return require(path.join(baseMigrationPath, migration.file))
  }
}

/*
  constants
*/

export default {
  appBuildVersion: process.env.VUE_APP_AppBuildVersion || '0.0.0',
  appBuildHash: process.env.VUE_APP_AppBuildHash || '???',
  appBuildDate: process.env.VUE_APP_AppBuildDate || '???',
  appBuildMode: process.env.VUE_APP_AppBuildMode || 'development'
}

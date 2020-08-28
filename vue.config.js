module.exports = {
  productionSourceMap: false,
  pwa: {
    name: "FlatTrack",
    themeColor: "#209cee",
    msTileColor: "#209cee",
    appleMobileWebAppCache: "yes",
    manifestOptions: {
      background_color: "#000000"
    },
    workboxOptions: {
      swSrc: 'service-worker.js',
      skipWaiting: true
    },
    workboxPluginMode: "GenerateSW"
  }
}

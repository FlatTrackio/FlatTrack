/*
  common
  commonly used JS functions
*/

const themes = [
  {
    name: 'light',
    icon: 'weather-sunny'
  }
]
const themeDefault = 'light'

// GetTheme
// returns the current theme
function GetTheme () {
  return JSON.parse(localStorage.getItem('appTheme')) || {}
}

// SetTheme
// sets the current theme
function SetTheme (name) {
  var theme = themes.filter((theme) => theme.name === name)
  if (typeof theme === 'undefined') {
    theme = themes.filter((theme) => theme.name === themeDefault)
  }
  return localStorage.setItem('appTheme', JSON.stringify(theme[0] || {}))
}

// SetThemeDefault
// sets the default theme
function SetThemeDefault () {
  SetTheme('light')
}

// SetThemeDefaultIfNotSet
// sets the default theme if there is no theme set
function SetThemeDefaultIfNotSet () {
  var currentTheme = GetTheme()
  if (typeof currentTheme.name === 'undefined') {
    SetTheme(themeDefault)
  }
}

export default {
  GetTheme,
  SetTheme,
  SetThemeDefault
}

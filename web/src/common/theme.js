/*
  common
  commonly used JS functions
*/

// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

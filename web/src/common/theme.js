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
    name: "system",
    icon: "cogs",
  },
  {
    name: "light",
    icon: "weather-sunny",
  },
  {
    name: "dark",
    icon: "weather-night",
  },
];
const themeDefault = themes[0];

// GetTheme
// returns the current theme
function GetTheme() {
  return JSON.parse(localStorage.getItem("appTheme")) || themes[0];
}

// SetTheme
// sets the current theme
function SetTheme(name) {
  let theme = themes.filter((theme) => theme.name === name);
  if (typeof theme === "undefined") {
    theme = themes.filter((theme) => theme.name === themeDefault);
  }
  document.documentElement.setAttribute("data-theme", theme[0].name);
  return localStorage.setItem("appTheme", JSON.stringify(theme[0] || {}));
}

function ListThemes() {
  return themes;
}

// SetThemeDefault
// sets the default theme
function SetThemeDefault() {
  SetTheme(themeDefault.name);
}

export default {
  GetTheme,
  SetTheme,
  SetThemeDefault,
  ListThemes,
};

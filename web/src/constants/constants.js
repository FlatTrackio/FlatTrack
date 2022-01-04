/*
  constants
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

export default {
  appBuildVersion: import.meta.env.VITE_APP_AppBuildVersion || '0.0.0',
  appBuildHash: import.meta.env.VITE_APP_AppBuildHash || '???',
  appBuildDate: import.meta.env.VITE_APP_AppBuildDate || '???',
  appBuildMode: import.meta.env.VITE_APP_AppBuildMode || 'development',
  appWebpackHotUpdate: typeof webpackHotUpdate === 'function' && import.meta.env.MODE === 'development'
}

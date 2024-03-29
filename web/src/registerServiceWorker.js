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

import { register } from 'register-service-worker'
import { ToastProgrammatic as Toast } from 'buefy'

if (
  'serviceWorker' in navigator &&
  localStorage.getItem('ft-no-sw') !== 'true'
) {
  register(`${import.meta.env.BASE_URL}sw.js`, {
    ready () {},
    registered () {},
    cached () {},
    updatefound () {
      window.location.reload(true)
    },
    updated () {
      Toast.open({
        message: 'FlatTrack was updated',
        type: 'is-info',
        position: 'is-bottom'
      })
    },
    offline () {},
    error (error) {
      console.error('Error during service worker registration:', error)
    }
  })
}

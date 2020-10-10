import { register } from 'register-service-worker'
import { SnackbarProgrammatic as Snackbar, ToastProgrammatic as Toast } from 'buefy'
import path from 'path'
import common from '@/frontend/common/common'
const subpath = common.GetSiteSubPath()

if (process.env.NODE_ENV === 'production') {
  register(path.join(subpath, '/service-worker.js'), {
    ready () {
    },
    registered () {
    },
    cached () {
    },
    updatefound () {
      window.location.reload(true)
    },
    updated () {
      console.log('Updated')
      Toast.open({
        message: 'FlatTrack was updated',
        type: 'is-success'
      })
    },
    offline () {
    },
    error (error) {
      console.error('Error during service worker registration:', error)
    }
  })
}

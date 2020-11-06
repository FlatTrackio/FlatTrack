import { register } from 'register-service-worker'
import { SnackbarProgrammatic as Snackbar, ToastProgrammatic as Toast } from 'buefy'

if (process.env.NODE_ENV === 'production') {
  register(`${process.env.BASE_URL}service-worker.js`, {
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

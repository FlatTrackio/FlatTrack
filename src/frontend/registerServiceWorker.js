import { register } from 'register-service-worker'

if (process.env.NODE_ENV === 'production') {
  register(`${process.env.BASE_URL}service-worker.js`, {
    ready () {
    },
    registered () {
    },
    cached () {
    },
    updatefound () {
      this.$buefy.snackbar.open({
        message: 'New FlatTrack version available, click to update',
        actionText: 'Update',
        type: 'is-white',
        position: 'is-right',
        indefinite: true,
        onAction: () => {
          const loadingComponent = this.$buefy.snackbar.open({
            container: null
          })
          setTimeout(() => {
            window.location.reload(true)
          }, 500)
        }
      })
    },
    updated () {
      console.log('Updated')
      this.$buefy.toast.open({
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

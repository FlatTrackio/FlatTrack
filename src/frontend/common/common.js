/*
  common
    commonly used JS functions
*/

import { ToastProgrammatic as Toast, DialogProgrammatic as Dialog, LoadingProgrammatic as Loading } from 'buefy'
import moment from 'moment'

// GetAuthToken
// returns the JWT from localStorage
function GetAuthToken () {
  return localStorage.getItem('authToken')
}

// DeleteAuthToken
// removes the JWT from localStorage
function DeleteAuthToken () {
  return localStorage.removeItem('authToken')
}

// DisplaySuccessToast
// shows a toast at the top of the screen with a message and a green background for 8 seconds
function DisplaySuccessToast (message) {
  Toast.open({
    duration: 8 * 1000,
    message: message,
    position: 'is-top',
    type: 'is-success',
    hasIcon: true,
    queue: false
  })
}

// DisplayFailureToast
// shows a toast at the top of the screen with a message and a red background for 8 seconds
function DisplayFailureToast (message) {
  Toast.open({
    duration: 8 * 1000,
    message: message,
    position: 'is-top',
    type: 'is-danger',
    hasIcon: true,
    queue: false
  })
}

// SignoutDialog
// shows a dialog prompt to signout
function SignoutDialog () {
  Dialog.confirm({
    message: 'Are you sure you want to sign out?',
    type: 'is-danger',
    hasIcon: 'true',
    onConfirm: () => {
      const loadingComponent = Loading.open({
        container: null
      })
      setTimeout(() => {
        common.DeleteAuthToken()
        window.location.href = '/login'
      }, 1 * 1000)
    }
  })
}

// TimestampToCalendar
// converts a unix timestamp to a human readable string
function TimestampToCalendar (timestamp) {
  return moment.unix(timestamp).calendar().toLowerCase()
}

export default {
  GetAuthToken,
  DeleteAuthToken,
  DisplaySuccessToast,
  DisplayFailureToast,
  SignoutDialog,
  TimestampToCalendar
}

/*
  common
    commonly used JS functions
*/

import { ToastProgrammatic as Toast } from 'buefy'

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

export default {
  GetAuthToken,
  DeleteAuthToken,
  DisplaySuccessToast,
  DisplayFailureToast
}
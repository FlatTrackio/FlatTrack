/*
  common
    commonly used JS functions
*/

import { ToastProgrammatic as Toast, DialogProgrammatic as Dialog, LoadingProgrammatic as Loading } from 'buefy'
import dayjs from 'dayjs'
import dayjsCalendar from 'dayjs/plugin/calendar'

dayjs.extend(dayjsCalendar)

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
    duration: 4 * 1000,
    message: message,
    position: 'is-bottom',
    type: 'is-success',
    hasIcon: true,
    queue: false
  })
}

// DisplayFailureToast
// shows a toast at the top of the screen with a message and a red background for 8 seconds
function DisplayFailureToast (message) {
  Toast.open({
    duration: 4 * 1000,
    message: message,
    position: 'is-bottom',
    type: 'is-danger',
    hasIcon: true,
    queue: false
  })
}

// SignoutDialog
// shows a dialog prompt to signout
function SignoutDialog (vm) {
  Dialog.confirm({
    message: 'Are you sure you want to sign out?',
    type: 'is-danger',
    hasIcon: 'true',
    onConfirm: () => {
      const loadingComponent = Loading.open({
        container: null
      })
      setTimeout(() => {
        DeleteAuthToken()
        vm.$router.push({ name: 'Login' })
        loadingComponent.close()
      }, 1 * 1000)
    }
  })
}

// TimestampToCalendar
// converts a unix timestamp to a human readable string
function TimestampToCalendar (timestamp) {
  return dayjs(timestamp * 1000).calendar(null, {}).toLowerCase()
}

// DeviceIsMobile
// returns bool if the device is mobile (from screen size)
function DeviceIsMobile () {
  return window.innerWidth <= 870
}

// GetFlatnameFromCache
// returns the flatname
function GetFlatnameFromCache () {
  return localStorage.getItem('flatname')
}

// WriteFlatnameToCache
// sets the flatname in the cache
function WriteFlatnameToCache (name) {
  localStorage.setItem('flatname', name)
}

export default {
  GetAuthToken,
  DeleteAuthToken,
  DisplaySuccessToast,
  DisplayFailureToast,
  SignoutDialog,
  TimestampToCalendar,
  DeviceIsMobile,
  GetFlatnameFromCache,
  WriteFlatnameToCache
}

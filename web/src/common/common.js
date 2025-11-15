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

import dayjs from "dayjs";
import dayjsCalendar from "dayjs/plugin/calendar";
import confetti from "canvas-confetti";
import login from "@/requests/public/login";

dayjs.extend(dayjsCalendar);

// DisplaySuccessToast
// shows a toast at the top of the screen with a message and a green background for 8 seconds
function DisplaySuccessToast(buefy, message) {
  buefy.toast.open({
    duration: 4 * 1000,
    message: message,
    position: "is-bottom",
    type: "is-success",
    hasIcon: true,
    queue: false,
  });
}

// DisplayFailureToast
// shows a toast at the top of the screen with a message and a red background for 8 seconds
function DisplayFailureToast(buefy, message) {
  buefy.toast.open({
    duration: 4 * 1000,
    message: message,
    position: "is-bottom",
    type: "is-danger",
    hasIcon: true,
    queue: false,
  });
}

// SignoutDialog
// shows a dialog prompt to signout
function SignoutDialog(buefy) {
  buefy.dialog.confirm({
    message: "Are you sure you want to sign out?",
    type: "is-danger",
    hasIcon: "true",
    onConfirm: () => {
      const loadingComponent = buefy.loading.open({
        container: null,
      });
      setTimeout(() => {
        login.DeleteUserAuth().then((resp) => {
          window.location.href = "/login";
          loadingComponent.close();
        });
      }, 1 * 1000);
    },
  });
}

// TimestampToCalendar
// converts a unix timestamp to a human readable string
function TimestampToCalendar(timestamp) {
  return dayjs(timestamp * 1000).calendar(null, {
    sameDay: "[today] h:mm A",
    nextDay: "[tomorrow] h:mm A",
    nextWeek: "dddd [at] h:mm A",
    lastDay: "[yesterday] h:mm A",
    lastWeek: "dddd [at] h:mm A",
    sameElse: "DD/MM/YYYY",
  });
}

// DeviceIsMobile
// returns bool if the device is mobile (from screen size)
function DeviceIsMobile() {
  return window.innerWidth <= 870;
}

// GetFlatnameFromCache
// returns the flatname
function GetFlatnameFromCache() {
  return localStorage.getItem("flatname");
}

// WriteFlatnameToCache
// sets the flatname in the cache
function WriteFlatnameToCache(name) {
  localStorage.setItem("flatname", name);
}

// GetHomeLastViewedTimestamp
// returns the timestamp of when the home page was last loaded or cleared
function GetHomeLastViewedTimestamp() {
  return Number(localStorage.getItem("homeLastViewedTimestamp"));
}

// WriteHomeLastViewedTimestamp
// sets the timestamp of when the home page was last loaded or cleared
function WriteHomeLastViewedTimestamp(timestamp) {
  localStorage.setItem("homeLastViewedTimestamp", String(timestamp));
}

// GetEnableAnimations
// returns if animations is enabled
function GetEnableAnimations() {
  var enabled = localStorage.getItem("enableAnimations");
  if (typeof enabled === "undefined") {
    WriteEnableAnimations("false");
    enabled = "false";
  }
  return enabled;
}

// WriteEnableAnimations
// sets animations to enabled or disabled
function WriteEnableAnimations(enable) {
  return localStorage.setItem("enableAnimations", enable);
}

// Hooray
// launch the confetti
function Hooray() {
  confetti({
    particleCount: 100,
    spread: 70,
    origin: { y: 1 },
    zIndex: 19,
  });
}

// GetSetupMessage
// returns a message to display on setup
function GetSetupMessage() {
  return (
    document.head.querySelector("[name~=setupmessage][content]")?.content || ""
  );
}

// GetLoginMessage
// returns a message to display on login
function GetLoginMessage() {
  return (
    document.head.querySelector("[name~=loginmessage][content]")?.content || ""
  );
}
function CopyHrefToClipboard() {
  navigator.clipboard.writeText(window.location.href);
}

export default {
  DisplaySuccessToast,
  DisplayFailureToast,
  SignoutDialog,
  TimestampToCalendar,
  DeviceIsMobile,
  GetFlatnameFromCache,
  WriteFlatnameToCache,
  GetHomeLastViewedTimestamp,
  WriteHomeLastViewedTimestamp,
  GetEnableAnimations,
  WriteEnableAnimations,
  Hooray,
  GetSetupMessage,
  GetLoginMessage,
  CopyHrefToClipboard,
};

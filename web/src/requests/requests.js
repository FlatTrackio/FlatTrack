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

import axios from "axios";
import common from "@/common/common";
import constants from "@/constants/constants";

function redirectToLogin(redirect) {
  if (redirect === false) {
    return;
  }
  if (window.location.pathname !== "/login") {
    let u = new URL(window.location.origin);
    u.pathname = "/login";
    u.searchParams.set("redirect", window.location.pathname);
    window.location = u.toString();
  }
}

function Request(request, redirect = true, publicRoute = false) {
  return new Promise((resolve, reject) => {
    request.headers = {
      Accept: "application/json",
    };
    if (constants.appWebpackHotUpdate) {
      request.baseURL = "http://localhost:8080";
    }
    axios(request)
      .then((resp) => resolve(resp))
      .catch((err) => {
        if (err.response.status === 401) {
          redirectToLogin(redirect);
          reject(err);
        } else if (err.response.status === 503) {
          window.location.href = "/unavailable";
        }
        reject(err);
      });
  });
}

export default Request;

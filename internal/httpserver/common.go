package httpserver

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

// JSONResponse ...
// form generic JSON responses
func JSONResponse(r *http.Request, w http.ResponseWriter, code int, output types.JSONMessageResponse) {
	// simpilify sending a JSON response
	output.Metadata.URL = r.URL.Path
	output.Metadata.Timestamp = time.Now().Unix()
	output.Metadata.Version = common.GetAppBuildVersion()
	response, _ := json.Marshal(output)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

// GetHTTPresponseBodyContents ...
// convert the body of a HTTP response into a JSONMessageResponse
func GetHTTPresponseBodyContents(response *http.Response) (output types.JSONMessageResponse) {
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Failed to read body contents", "error", err)
	}
	if err := json.Unmarshal(responseData, &output); err != nil {
		slog.Error("Failed to unmarshal response body contents", "error", err)
	}
	return output
}

// GetRequestIP ...
// returns r.RemoteAddr unless RealIPHeader is set
func GetRequestIP(r *http.Request) (requestIP string) {
	realIPHeader := common.GetAppRealIPHeader()
	headerValue := r.Header.Get(realIPHeader)
	if realIPHeader == "" || headerValue == "" {
		return r.RemoteAddr
	}
	return headerValue
}

// SetTokenCookie ...
// sets the auth token in the token cookie
func (h *HTTPServer) SetTokenCookie(w http.ResponseWriter, token string) {
	secure := h.instanceURL == nil || h.instanceURL != nil && h.instanceURL.Scheme != "http"
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearTokenCookie ...
// clears the token cookie
func (h *HTTPServer) ClearTokenCookie(w http.ResponseWriter) {
	secure := h.instanceURL == nil || h.instanceURL != nil && h.instanceURL.Scheme != "http"
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

package httpserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

// JSONResponse ...
// form generic JSON responses
func JSONResponse(r *http.Request, w http.ResponseWriter, code int, output types.JSONMessageResponse) {
	// simpilify sending a JSON response
	output.Metadata.URL = r.RequestURI
	output.Metadata.Timestamp = time.Now().Unix()
	output.Metadata.Version = common.GetAppBuildVersion()
	response, _ := json.Marshal(output)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		log.Printf("error: failed to write response; %v\n", err)
	}
}

// GetHTTPresponseBodyContents ...
// convert the body of a HTTP response into a JSONMessageResponse
func GetHTTPresponseBodyContents(response *http.Response) (output types.JSONMessageResponse) {
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(responseData, &output); err != nil {
		log.Printf("error: failed to unmarshal response body contents; %v", err)
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

/*
	handle all types used by API
*/

package types

import (
	"net/http"
)

type User struct {
	TypeMeta `json:",inline"`
	JSONResponseMetadata `json:"metadata"`
	UserSpec `json:"spec"`
}

// UserSpec
// standard user account objects
type UserSpec struct {
	Id                        string `json:"id"`
	Name                      string `json:"name"`
	Email                     string `json:"email"`
	Groups                    string `json:"groups"`
	Password                  string `json:"password,omitempty"`
	PhoneNumber               string `json:"phoneNumber,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	HasSetPassword            bool   `json:"hasSetPassword,omitempty"`
	TaskNotificationFrequency int    `json:"taskNotificationFrequency,omitempty"`
	LastLogin                 string `json:"lastLogin"`
	CreationTimestamp         string `json:"creationTimestamp"`
	ModificationTimestamp     string `json:"modificationTimestamp"`
	DeletionTimestamp         string `json:"deletionTimestamp"`
}

// UserList
// include multiple user accounts
type UserList []User

type Group string

const (
	Flatmember Group = "flatmember"
	Admin      Group = "admin"
)

type JSONResponseMetadata struct {
	URL       string `json:"selfLink"`
	Version   string `json:"version"`
	RequestId string `json:"requestId"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

type JSONMessageResponse struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     interface{}          `json:"spec"`
}

type TypeMeta struct {
	Kind string `json:"kind"`
}

type Endpoints []struct {
	EndpointPath string
	HandlerFunc  http.HandlerFunc
	HttpMethod   string
}

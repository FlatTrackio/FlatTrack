/*
  types
    collection of go types used in FlatTrack's API
*/

package types

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// Group
// request object for a group
type Group struct {
	TypeMeta `json:",inline"`
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     GroupSpec            `json:"spec"`
}

// GroupSpec
// standard values for a group
type GroupSpec struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DefaultGroup bool   `json:"defaultGroup"`
	Description  string `json:"description"`
}

// GroupList
// a list of groups
type GroupList struct {
	TypeMeta `json:",inline"`
	Metadata JSONResponseMetadata `json:"metadata"`
	List     []GroupSpec          `json:"list"`
}

// User
// request object for a user account
type User struct {
	TypeMeta `json:",inline"`
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     UserSpec             `json:"spec"`
}

// UserSpec
// standard user account objects
type UserSpec struct {
	Id                        string   `json:"id"`
	Names                     string   `json:"names"`
	Email                     string   `json:"email"`
	Groups                    []string `json:"groups"`
	Password                  string   `json:"password,omitempty"`
	PhoneNumber               string   `json:"phoneNumber,omitempty"`
	Birthday                  string   `json:"birthday,omitempty"`
	ContractAgreement         bool     `json:"contractAgreement,omitempty"`
	Disabled                  bool     `json:"disabled,omitempty"`
	HasSetPassword            bool     `json:"hasSetPassword,omitempty"`
	TaskNotificationFrequency int      `json:"taskNotificationFrequency,omitempty"`
	LastLogin                 string   `json:"lastLogin"`
	CreationTimestamp         int64    `json:"creationTimestamp"`
	ModificationTimestamp     int64    `json:"modificationTimestamp"`
	DeletionTimestamp         int64    `json:"deletionTimestamp"`
}

// UserList
// include multiple user accounts
type UserList struct {
	TypeMeta `json:",inline"`
	Metadata JSONResponseMetadata `json:"metadata"`
	List     []UserSpec           `json:"list"`
}

// FlatName
// the name of the flat
type FlatName struct {
	FlatName string `json:"flatName"`
}

// Registration
// fields to initialize the instance of FlatTrack
type Registration struct {
	User     UserSpec `json:"user"`
	Timezone string   `json:"timezone"`
	Language string   `json:"language"`
	FlatName string   `json:"flatName"`
}

// JSONResponseMetadata
// values to return in each request
type JSONResponseMetadata struct {
	URL       string `json:"selfLink"`
	Version   string `json:"version"`
	RequestId string `json:"requestId"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

// JSONMessageResponse
// generic JSON response
type JSONMessageResponse struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     interface{}          `json:"spec,omitempty"`
	List     interface{}          `json:"list,omitempty"`
	Data     interface{}          `json:"data,omitempty"`
}

type TypeMeta struct {
	Kind string `json:"kind"`
}

// Endpoints
// all API endpoints stored in an array
type Endpoints []struct {
	EndpointPath string
	HandlerFunc  http.HandlerFunc
	HttpMethod   string
}

// JWTclaim
// contents for JWT token
type JWTclaim struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

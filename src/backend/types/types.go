/*
	handle all types used by API
*/

package types

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type Group struct {
	TypeMeta `json:",inline"`
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     GroupSpec            `json:"spec"`
}

type GroupSpec struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DefaultGroup bool   `json:"defaultGroup"`
}

type GroupList struct {
	TypeMeta `json:",inline"`
	Metadata JSONResponseMetadata `json:"metadata"`
	List     []GroupSpec          `json:"list"`
}

// User
//
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

type JSONResponseMetadata struct {
	URL       string `json:"selfLink"`
	Version   string `json:"version"`
	RequestId string `json:"requestId"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

type JSONMessageResponse struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     interface{}          `json:"spec,omitempty"`
	List     interface{}          `json:"list,omitempty"`
	Data     interface{}          `json:"data,omitempty"`
}

type TypeMeta struct {
	Kind string `json:"kind"`
}

type Endpoints []struct {
	EndpointPath string
	HandlerFunc  http.HandlerFunc
	HttpMethod   string
}

type JWTclaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

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
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	DefaultGroup          bool   `json:"defaultGroup"`
	Description           string `json:"description"`
	CreationTimestamp     int    `json:"creationTimestamp"`
	ModificationTimestamp int    `json:"modificationTimestamp"`
	DeletionTimestamp     int    `json:"deletionTimestamp"`
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
	Birthday                  int      `json:"birthday,omitempty"`
	ContractAgreement         bool     `json:"contractAgreement,omitempty"`
	Disabled                  bool     `json:"disabled,omitempty"`
	Registered                bool     `json:"registered,omitempty"`
	TaskNotificationFrequency int      `json:"taskNotificationFrequency,omitempty"`
	LastLogin                 string   `json:"lastLogin"`
	AuthNonce                 string   `json:"uuthNonce"`
	CreationTimestamp         int      `json:"creationTimestamp"`
	ModificationTimestamp     int      `json:"modificationTimestamp"`
	DeletionTimestamp         int      `json:"deletionTimestamp"`
}

// UserList
// include multiple user accounts
type UserList struct {
	TypeMeta `json:",inline"`
	Metadata JSONResponseMetadata `json:"metadata"`
	List     []UserSpec           `json:"list"`
}

type UserSelector struct {
	Group string `json:"group,omitempty"`
	notId string `json:"notId,omitempty"`
}

type ShoppingListSpec struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	Notes                 string `json:"notes"`
	Author                string `json:"author"`
	AuthorLast            string `json:"authorLast"`
	Completed             bool   `json:"completed"`
	Count                 int    `json:"count,omitempty"`
	CreationTimestamp     int    `json:"creationTimestamp"`
	ModificationTimestamp int    `json:"modificationTimestamp"`
	DeletionTimestamp     int    `json:"deletionTimestamp"`
}

type ShoppingItemSpec struct {
	Id                    string  `json:"id"`
	Name                  string  `json:"name"`
	Price                 float64 `json:"price"`
	Regular               bool    `json:"regular"`
	Notes                 string  `json:"notes"`
	Obtained              string  `json:"obtained"`
	Author                string  `json:"author"`
	AuthorLast            string  `json:"authorLast"`
	CreationTimestamp     int     `json:"creationTimestamp"`
	ModificationTimestamp int     `json:"modificationTimestamp"`
	DeletionTimestamp     int     `json:"deletionTimestamp"`
}

type ShoppingItemSelector struct {
	Regular bool `json:"regular,omitempty"`
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
	Id        string `json:"id"`
	AuthNonce string `json:"authNonce"`
	jwt.StandardClaims
}

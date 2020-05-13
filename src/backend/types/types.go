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
	Id                    string   `json:"id"`
	Names                 string   `json:"names"`
	Email                 string   `json:"email"`
	Groups                []string `json:"groups"`
	Password              string   `json:"password,omitempty"`
	PhoneNumber           string   `json:"phoneNumber,omitempty"`
	Birthday              int64    `json:"birthday,omitempty"`
	ContractAgreement     bool     `json:"contractAgreement,omitempty"`
	Disabled              bool     `json:"disabled,omitempty"`
	Registered            bool     `json:"registered"`
	LastLogin             int      `json:"lastLogin,omitempty"`
	AuthNonce             string   `json:"-"`
	CreationTimestamp     int      `json:"creationTimestamp"`
	ModificationTimestamp int      `json:"modificationTimestamp"`
	DeletionTimestamp     int      `json:"deletionTimestamp"`
}

// UserList
// include multiple user accounts
type UserList struct {
	TypeMeta `json:",inline"`
	Metadata JSONResponseMetadata `json:"metadata"`
	List     []UserSpec           `json:"list"`
}

type UserSelector struct {
	Id      string `json:"id,omitempty"`
	Group   string `json:"group,omitempty"`
	NotId   string `json:"notId,omitempty"`
	NotSelf string `json:"notSelf,omitempty"`
}

type ShoppingListSpec struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	Notes                 string `json:"notes,omitempty"`
	TemplateId            string `json:"templateId,omitempty"`
	Completed             bool   `json:"completed"`
	Count                 int    `json:"count,omitempty"`
	Author                string `json:"author"`
	AuthorLast            string `json:"authorLast"`
	CreationTimestamp     int    `json:"creationTimestamp"`
	ModificationTimestamp int    `json:"modificationTimestamp"`
	DeletionTimestamp     int    `json:"deletionTimestamp"`
}

type ShoppingItemSpec struct {
	Id                    string  `json:"id"`
	ListId                string  `json:"listId"`
	Name                  string  `json:"name"`
	Price                 float64 `json:"price,omitempty"`
	Quantity              int     `json:"quantity"`
	Notes                 string  `json:"notes"`
	Obtained              bool    `json:"obtained"`
	Tag                   string  `json:"tag,omitempty"`
	Author                string  `json:"author"`
	AuthorLast            string  `json:"authorLast"`
	CreationTimestamp     int     `json:"creationTimestamp"`
	ModificationTimestamp int     `json:"modificationTimestamp"`
	DeletionTimestamp     int     `json:"deletionTimestamp"`
}

type ShoppingItemSortType string

const (
	ShoppingItemSortByTag                    = "tag"
	ShoppingItemSortByHighestPrice           = "highestPrice"
	ShoppingItemSortByHighestQuantity        = "highestQuantity"
	ShoppingItemSortByLowestPrice            = "lowestPrice"
	ShoppingItemSortByLowestQuantity         = "lowestQuantity"
	ShoppingItemSortByRecentlyAdded          = "recentlyAdded"
	ShoppingItemSortByRecentlyUpdated        = "recentlyUpdated"
	ShoppingItemSortByLastAdded              = "lastAdded"
	ShoppingItemSortByLastUpdated            = "lastUpdated"
	ShoppingItemSortByAlphabeticalDescending = "alphabeticalDescending"
	ShoppingItemSortByAlphabeticalAscending  = "alphabeticalAscending"
)

type ShoppingItemOptions struct {
	Selector ShoppingItemSelector `json:"Selector"`
	SortBy   string               `json:"sortBy"`
}

type ShoppingItemSelector struct {
	TemplateListItemSelector string `json:"templateListItemSelector"`
}

type ShoppingItemTag struct {
	Name string `json:"name"`
}

type UserCreationSecretSpec struct {
	Id                    string `json:"id"`
	UserId                string `json:"userId"`
	Secret                string `json:"secret"`
	Valid                 bool   `json:"valid"`
	CreationTimestamp     int    `json:"creationTimestamp"`
	ModificationTimestamp int    `json:"modificationTimestamp"`
	DeletionTimestamp     int    `json:"deletionTimestamp"`
}

type UserCreationSecretSelector struct {
	UserId string `json:"userId"`
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

type SystemVersion struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	Mode       string `json:"mode"`
	Date       string `json:"date"`
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

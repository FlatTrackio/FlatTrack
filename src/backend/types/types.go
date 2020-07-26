/*
  types
    collection of go types used in FlatTrack's API
*/

package types

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// Group ...
// request object for a group
type Group struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     GroupSpec            `json:"spec"`
}

// GroupSpec ...
// standard values for a group
type GroupSpec struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	DefaultGroup          bool   `json:"defaultGroup"`
	Description           string `json:"description"`
	CreationTimestamp     int    `json:"creationTimestamp"`
	ModificationTimestamp int    `json:"modificationTimestamp"`
	DeletionTimestamp     int    `json:"deletionTimestamp"`
}

// GroupList ...
// a list of groups
type GroupList struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	List     []GroupSpec          `json:"list"`
}

// User ...
// request object for a user account
type User struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     UserSpec             `json:"spec"`
}

// UserSpec ...
// standard user account objects
// swagger:response userSpec
type UserSpec struct {
	ID                    string   `json:"id"`
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

// UserList ...
// include multiple user accounts
type UserList struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	List     []UserSpec           `json:"list"`
}

// UserSelector ...
// fields for filtering user account lists
type UserSelector struct {
	ID      string `json:"id,omitempty"`
	Group   string `json:"group,omitempty"`
	NotID   string `json:"notId,omitempty"`
	NotSelf string `json:"notSelf,omitempty"`
}

// ShoppingListSpec ...
// fields for a shopping list
type ShoppingListSpec struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Notes                 string `json:"notes,omitempty"`
	TemplateID            string `json:"templateId,omitempty"`
	Completed             bool   `json:"completed"`
	Count                 int    `json:"count,omitempty"`
	Author                string `json:"author"`
	AuthorLast            string `json:"authorLast"`
	CreationTimestamp     int    `json:"creationTimestamp"`
	ModificationTimestamp int    `json:"modificationTimestamp"`
	DeletionTimestamp     int    `json:"deletionTimestamp"`
}

// ShoppingListSortType ...
// ways of sorting shopping lists
type ShoppingListSortType string

// ShoppingListSortTypes ...
// ways of sorting shopping lists
const (
	ShoppingListSortByRecentlyAdded          = "recentlyAdded"
	ShoppingListSortByRecentlyUpdated        = "recentlyUpdated"
	ShoppingListSortByLastAdded              = "lastAdded"
	ShoppingListSortByLastUpdated            = "lastUpdated"
	ShoppingListSortByAlphabeticalDescending = "alphabeticalDescending"
	ShoppingListSortByAlphabeticalAscending  = "alphabeticalAscending"
)

// ShoppingListOptions ...
// options for lists
type ShoppingListOptions struct {
	SortBy   string               `json:"sortBy"`
	Selector ShoppingListSelector `json:"selector"`
}

// ShoppingListSelector ...
// options for creating and selecting lists
type ShoppingListSelector struct {
	Completed string `json:"completed"`
}

// ShoppingItemSpec ...
// fields for a shopping item
type ShoppingItemSpec struct {
	ID                    string  `json:"id"`
	ListID                string  `json:"listId"`
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

// ShoppingItemSortType ...
// ways of sorting shopping list items
type ShoppingItemSortType string

// ShoppingItemSortTypes ...
// ways of sorting shopping list items
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

// ShoppingItemOptions ...
// options for list items
type ShoppingItemOptions struct {
	Selector ShoppingItemSelector `json:"selector"`
	SortBy   string               `json:"sortBy"`
}

// ShoppingItemSelector ...
// options for creating and selecting lists
type ShoppingItemSelector struct {
	TemplateListItemSelector string `json:"templateListItemSelector"`
}

// ShoppingTag ...
// selects a tag
type ShoppingTag struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Author                string `json:"author"`
	AuthorLast            string `json:"authorLast"`
	CreationTimestamp     int    `json:"creationTimestamp"`
	ModificationTimestamp int    `json:"modificationTimestamp"`
	DeletionTimestamp     int    `json:"deletionTimestamp"`
}

// ShoppingTagOptions ...
// options for list items
type ShoppingTagOptions struct {
	SortBy   string               `json:"sortBy"`
}

// ShoppingTagSortTypes ...
// ways of sorting shopping tags
const (
	ShoppingTagSortByRecentlyAdded          = "recentlyAdded"
	ShoppingTagSortByRecentlyUpdated        = "recentlyUpdated"
	ShoppingTagSortByLastAdded              = "lastAdded"
	ShoppingTagSortByLastUpdated            = "lastUpdated"
	ShoppingTagSortByAlphabeticalDescending = "alphabeticalDescending"
	ShoppingTagSortByAlphabeticalAscending  = "alphabeticalAscending"
)

// ShoppingListNotes ...
// notes for shopping lists
type ShoppingListNotes struct {
	Notes string `json:"notes"`
}

// UserCreationSecretSpec ...
// values for a user to confirm their account with
type UserCreationSecretSpec struct {
	ID                    string `json:"id"`
	UserID                string `json:"userId"`
	Secret                string `json:"secret"`
	Valid                 bool   `json:"valid"`
	CreationTimestamp     int    `json:"creationTimestamp"`
	ModificationTimestamp int    `json:"modificationTimestamp"`
	DeletionTimestamp     int    `json:"deletionTimestamp"`
}

// UserCreationSecretSelector ...
// filters the userCreationSecrets
type UserCreationSecretSelector struct {
	UserID string `json:"userId"`
}

// FlatName ...
// the name of the flat
type FlatName struct {
	FlatName string `json:"flatName"`
}

// Registration ...
// fields to initialize the instance of FlatTrack
type Registration struct {
	User     UserSpec `json:"user"`
	Timezone string   `json:"timezone"`
	Language string   `json:"language"`
	FlatName string   `json:"flatName"`
}

// SystemVersion ...
// values for the release of FlatTrack
type SystemVersion struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	Mode       string `json:"mode"`
	Date       string `json:"date"`
}

// JSONResponseMetadata ...
// values to return in each request
type JSONResponseMetadata struct {
	URL       string `json:"selfLink"`
	Version   string `json:"version"`
	RequestID string `json:"requestId"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

// JSONMessageResponse ...
// generic JSON response
type JSONMessageResponse struct {
	Metadata JSONResponseMetadata `json:"metadata"`
	Spec     interface{}          `json:"spec,omitempty"`
	List     interface{}          `json:"list,omitempty"`
	Data     interface{}          `json:"data,omitempty"`
}

// Endpoints ...
// all API endpoints stored in an array
type Endpoints []struct {
	EndpointPath string
	HandlerFunc  http.HandlerFunc
	HTTPMethod   string
}

// JWTclaim ...
// contents for JWT token
type JWTclaim struct {
	ID        string `json:"id"`
	AuthNonce string `json:"authNonce"`
	jwt.StandardClaims
}

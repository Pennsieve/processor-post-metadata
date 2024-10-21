package models

// APIResponse is a minimal type that will only capture the name and id of
// model service responses. Create new types that build on this one if more
// info from a particular response is required.
type APIResponse struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

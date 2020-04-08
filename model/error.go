package model

// Error structure is using for formatting errors according to swagger specification
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

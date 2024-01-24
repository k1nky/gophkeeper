package rest

//go:generate easyjson rest.go
//easyjson:json
type NewSecretRequest struct {
	Extra  string `json:"extra"`
	Secret string `json:"secret"`
}

//go:generate easyjson rest.go
//easyjson:json
type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//go:generate easyjson rest.go
//easyjson:json
type RegisterUserRequest struct {
	LoginUserRequest
}

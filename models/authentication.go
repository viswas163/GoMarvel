package models

// AuthClient : Authentication Client model
type AuthClient struct {
	PublicKey  string
	PrivateKey string
}

// Authenticator : Provides the query params for Client Authentication
type Authenticator struct {
	Timestamp string `url:"ts"`
	PublicKey string `url:"apikey"`
	Hash      string `url:"hash"`
}

// APIError : Error returned by the service
type APIError struct {
	Code    interface{}
	Message string
}

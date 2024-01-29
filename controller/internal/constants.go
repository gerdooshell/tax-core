package internal

type ContextKey string

const (
	APIKey                ContextKey = "apikey"
	MethodContextKey      ContextKey = "method"
	ContactBodyContextKey ContextKey = "method"

	APIKyeNameID string = "x-api-key"
	APIKeyValue  string = "123qwe"
)

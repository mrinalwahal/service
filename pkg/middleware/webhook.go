package middleware

import "net/http"

// X-Webhook-Token is the key used to store the webhook token in the request header.
//
// The webhook token is used to authenticate a webhook request.
const XWebhookToken Key = "X-Webhook-Token"

// Webhook middleware authenticates the request using a unique webhook token.
type WebhookConfig struct {

	// Token is the unique token that will be used to authenticate the request.
	//
	// This field is mandatory.
	Token string
}

// Webhook middleware authenticates the request using a unique webhook token.
func Webhook(config *WebhookConfig) Middleware {

	// Set the default configuration.
	if config == nil {
		config = &WebhookConfig{}
	}

	if config.Token == "" {
		panic("middleware: webhook: token is required")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Extract the token from the request header.
			token := r.Header.Get(string(XWebhookToken))

			// Check if the token is valid.
			if token != config.Token {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

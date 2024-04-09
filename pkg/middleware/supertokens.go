package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SuperTokensWithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			// we add content-type + other headers used by SuperTokens
			w.Header().Set("Access-Control-Allow-Headers",
				strings.Join(append([]string{"Content-Type"},
					supertokens.GetAllCORSHeaders()...), ","))
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Write([]byte(""))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

const UserID Key = "user_id"

func SuperTokensWithSessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userID := sessionContainer.GetUserID()

		// Write the user ID to the request context.
		ctx := context.WithValue(r.Context(), UserID, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
)

func TestJWT(t *testing.T) {

	t.Run("jwt middleware", func(t *testing.T) {

		// Initialize a new router.
		router := http.NewServeMux()

		// Initialize the JWT middleware.
		middleware := JWT(&JWTConfig{
			Key: "secret",
		})

		// Add the middleware to the router.
		router.Handle("/protected", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Read the claims from the request context.
			_, ok := r.Context().Value(XJWTClaims).(jwt.MapClaims)
			if !ok {
				http.Error(w, "failed to parse the claims", http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
		})))

		// Initialize test r and response recorder.
		r := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

		// Attach a dummy JWT to the request.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "3742a2cd-8958-41c1-aba6-ca66c6f3220d",
			"claims": map[string]interface{}{
				"x-user-id": "3742a2cd-8958-41c1-aba6-ca66c6f3220d",
			},
			"aud":   "",
			"iss":   "record",
			"scope": "",
		})
		signed, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Fatal(err)
		}

		r.Header.Add("Authorization", signed)

		// Serve the request.
		router.ServeHTTP(w, r)

		// Validate the status code.
		if status := w.Code; status != http.StatusOK {
			t.Logf("Response: %s", w.Body.String())
			t.Errorf("ServeHTTP() = %v, want %v", status, http.StatusOK)
		}
	})
}

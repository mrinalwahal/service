package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {

	t.Run("jwt middleware", func(t *testing.T) {

		// Initialize a new router.
		router := http.NewServeMux()

		// Initialize the JWT middleware.
		middleware := JWT(&JWTConfig{
			Key: "secret",
		})

		// Attach a dummy JWT to the request.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
			StandardClaims: jwt.StandardClaims{
				Subject: "3742a2cd-8958-41c1-aba6-ca66c6f3220d",
				Issuer:  "record",
			},
			XUserID: uuid.New(),
		})
		signed, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Fatal(err)
		}

		// Add the middleware to the router.
		router.Handle("/protected", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Read the claims from the request context.
			claims, exists := r.Context().Value(XJWTClaims).(JWTClaims)
			if !exists {
				http.Error(w, "failed to parse the claims", http.StatusUnauthorized)
				return
			}

			if claims.XUserID == uuid.Nil {
				t.Logf("Claims received: %v", claims)
				t.Errorf("invalid user_id in claims")
			}

			w.WriteHeader(http.StatusOK)
		})))

		// Initialize test r and response recorder.
		r := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

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

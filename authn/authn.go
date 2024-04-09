package authn

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
}

var database = []User{
	{
		ID:       uuid.New(),
		Email:    "test@gmail.com",
		Password: "password",
	},
}

// GetUserByEmailAndPassword returns the user with the given email and password.
func GetUserByEmailAndPassword(email, password string) (*User, error) {
	for _, user := range database {
		if user.Email == email && user.Password == password {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func GenerateJWT(userID string) string {

	// Generate a new JWT token
	jwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
	}).SignedString([]byte("secret"))
	if err != nil {
		return ""
	}

	return jwt
}

package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Username string `json:"username"`
	Role     string
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(username, role string) (string, error) {
	expirationTime := time.Now().Add(60 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT validates the JWT token
// func ParseToken(tokenStr string) (*Claims, error) {
// 	claims := &Claims{}

// 	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if !token.Valid {
// 		return nil, jwt.NewValidationError("Invalid token", jwt.ValidationErrorSignatureInvalid)
// 	}

//		return claims, nil
//	}
func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}
	//claim := token.Claims.(*Claims)

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Println("Username:", claims.Username)
		fmt.Println("Role:", claims.Role)
		fmt.Println("Issued At:", time.Unix(claims.IssuedAt, 0))
	} else {
		fmt.Println("Invalid token")
	}
	// if claim.Role != "admin" {
	// 	return nil, nil
	// }
	return token, nil
}

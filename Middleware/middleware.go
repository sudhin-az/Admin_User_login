package middleware

import (
	"fmt"
	"net/http"
	auth "user_admin/helpers"
)

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve JWT token from the cookie
		cookie, err := r.Cookie("jwt_admin_token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Missing authorization cookie")
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error retrieving cookie: %v", err)
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization token")
			return
		}

		token, err := auth.ParseToken(tokenString)
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Invalid authorization token: %v", err)
			return
		}

		// claims := token.Claims.(*Claims)
		// role:=claims.Role
		// username:=claims.Username
		// if role=="admin" {}
		// role, val := token.Claims

		next.ServeHTTP(w, r)
	})
}

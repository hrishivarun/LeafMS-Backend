package controller

import (
	"LeafMS-BackEnd/service"
	"context"
	"net/http"
)

// MIDDLEWARES!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func HandleValidateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		claims, err := service.VerifyToken(authHeader)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HandleValidateAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggedUsername, ok := r.Context().Value("username").(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if loggedUsername != service.Admin {
			http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

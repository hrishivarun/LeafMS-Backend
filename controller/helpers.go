package controller

import (
	db "LeafMS-BackEnd/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

var database = db.ConnectDB()

// generate JWT token
func generateJWT(username string) (string, error) {
	secretKey := []byte("jiteshmc" + username)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string, username string) error {
	secretKey := []byte("jiteshmc" + username)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// function to validate the db.user
func validateCred(userToAuthorize db.User) (db.User, db.UserLogin) {
	var login db.UserLogin
	data, err := database.FindOne("employees", bson.D{
		{Key: "username", Value: userToAuthorize.Username},
		{Key: "password", Value: userToAuthorize.Password}})
	if err != nil {
		login.Login = false
		log.Fatal("Failed authentication. Error:- \n\t", err)
		return db.User{}, login
	}

	var user db.User
	err = bson.Unmarshal(data, &user)
	if err != nil {
		log.Fatal("Couldn't unwrap the user data recieved from mongoDB.\nError:-\n\n", err)
	}

	if user.Username == "" {
		login.Login = false
		return user, login
	} else {
		login.Username = user.Username
		login.Login = true
	}
	return user, login
}

// MIDDLEWARES!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func HandleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		jwtToken := r.Header.Get("Authorization")
		sessionId := r.Header.Get("Session-Id")

		err := verifyToken(jwtToken, sessionId)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

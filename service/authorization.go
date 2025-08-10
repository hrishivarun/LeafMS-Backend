package service

import (
	"LeafMS-BackEnd/database"
	"LeafMS-BackEnd/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecretKey)
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return JWTSecretKey, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// function to validate the models.user
func ValidateAdminCred(userToAuthorize models.Employee) (models.Employee, models.LoginInfo) {

	loginInfo := models.LoginInfo{
		Username: userToAuthorize.Username,
	}
	data, err := database.DbConn.FindOne("employees", bson.D{
		{Key: "username", Value: Admin},
		{Key: "password", Value: userToAuthorize.Password}}, nil)
	if err != nil {
		loginInfo.Status = http.StatusUnauthorized
		log.SetPrefix("WARNING: ")
		log.Println("Failed authentication. Error:- \n\t", err)
		return models.Employee{}, loginInfo
	}

	var user models.Employee
	err = bson.Unmarshal(data, &user)
	if err != nil {
		loginInfo.Status = http.StatusNotFound
		log.SetPrefix("WARNING: ")
		log.Println("Couldn't unwrap the user data recieved from mongoDB.\nError:-\n\n", err)
	}

	if user.Username == "" {
		loginInfo.Status = http.StatusNotFound
		return user, loginInfo
	}

	token, err := GenerateJWT(userToAuthorize.Username)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Failed to generate token")
		loginInfo.Status = http.StatusInternalServerError
		return models.Employee{}, loginInfo
	}

	loginInfo.Token = token
	return userToAuthorize, loginInfo
}

// function to validate the models.user
func ValidateCred(userToAuthorize models.LoginReq) models.LoginInfo {
	loginInfo := models.LoginInfo{
		Username: userToAuthorize.Username,
	}
	data, err := database.DbConn.FindOne("employees", bson.D{
		{Key: "username", Value: userToAuthorize.Username},
		{Key: "password", Value: userToAuthorize.Password}}, nil)
	if err != nil {
		loginInfo.Status = http.StatusUnauthorized
		log.SetPrefix("WARNING: ")
		log.Println("Failed authentication. Error:- \n\t", err)
		return loginInfo
	}

	var user models.LoginReq
	err = bson.Unmarshal(data, &user)
	if err != nil {
		loginInfo.Status = http.StatusNotFound
		log.SetPrefix("WARNING: ")
		log.Println("Couldn't unwrap the user data recieved from mongoDB.\nError:-\n\n", err)
	}

	if user.Username == "" {
		loginInfo.Status = http.StatusNotFound
		return loginInfo
	}

	token, err := GenerateJWT(userToAuthorize.Username)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Failed to generate token")
		loginInfo.Status = http.StatusInternalServerError
		return loginInfo
	}

	loginInfo.Status = http.StatusOK
	loginInfo.Token = token
	return loginInfo
}

func ValidateApprover(employeeUsername *string, approverUsername string) bool {
	if employeeUsername == nil {
		return true
	}
	applierInfoRaw, err := database.DbConn.FindOne("employees", bson.D{{Key: "username", Value: *employeeUsername}}, nil)
	if err != nil {
		return false
	}
	var applierInfo models.Employee
	if err = bson.Unmarshal(applierInfoRaw, &applierInfo); err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("what the fuck? why is there a problem deserializing employee from database entry?")
		return false
	}
	return applierInfo.ApproverUserName == approverUsername
}

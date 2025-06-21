package service

import (
	"LeafMS-BackEnd/database"
	"LeafMS-BackEnd/models"
	"LeafMS-BackEnd/utils"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

// generate JWT token
func GenerateJWT(username string) (string, error) {
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

func VerifyToken(tokenString string, username string) error {
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

// function to validate the models.user
func ValidateCred(userToAuthorize models.Employee) (models.Employee, models.UserLogin) {
	var login models.UserLogin
	data, err := database.DbConn.FindOne("employees", bson.D{
		{Key: "username", Value: userToAuthorize.Username},
		{Key: "password", Value: userToAuthorize.Password}})
	if err != nil {
		login.Login = false
		log.Fatal("Failed authentication. Error:- \n\t", err)
		return models.Employee{}, login
	}

	var user models.Employee
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

func FilterHolidaysFromLeaveRequest(leaveApplication models.MetaLeaveInfo) (models.MetaLeaveInfo, error) {
	var splitLeaves []models.LeaveInfo
	for _, leave := range leaveApplication.Leaves {
		leaveSlices, err := utils.RemoveHolidayFromLeaveData(leave)
		if err != nil {
			log.Println("Could not remove the holidays from the leave applied. Err : ", err)
			return models.MetaLeaveInfo{}, err
		}

		splitLeaves = append(splitLeaves, leaveSlices...)
	}

	var leavesLackingWeekend []models.LeaveInfo
	for _, leave := range splitLeaves {
		leaveSlices, err := utils.RemoveWeekendsFromLeaveData(leave)
		if err != nil {
			log.Fatalln("There was an error while removing weekends from the applied leave. Err : ", err)
		}

		leavesLackingWeekend = append(leavesLackingWeekend, leaveSlices...)
	}

	leaveApplication.Leaves = leavesLackingWeekend
	return leaveApplication, nil
}

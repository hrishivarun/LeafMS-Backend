package utils

import (
	models "LeafMS-BackEnd/models"
	"fmt"
	"strings"
)

func InterFaceToUser(user interface{}) models.Employee {
	var result models.Employee
	str := fmt.Sprintf("%v", user)
	for len(str) > 4 {
		i1 := strings.Index(str, "{")
		i2 := strings.Index(str, "}")
		temp := str[i1+1 : i2]
		str = str[i2+1:]
		temp2 := strings.Split(temp, " ")
		result = setUserVal(temp2, result)
	}
	return result
}

func setUserVal(temp2 []string, user models.Employee) models.Employee {
	switch temp2[0] {
	case "username":
		user.Username = temp2[1]
	case "password":
		user.Password = temp2[1]
	case "name":
		user.FirstName = temp2[1]
	case "team":
		user.Team = temp2[1]
	case "designation":
		user.Designation = temp2[1]
	case "approver":
		user.Approver = temp2[1]
	}
	return user
}

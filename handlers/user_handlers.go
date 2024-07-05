package handlers

import (
	"encoding/json"
	"http-server/database"
	"http-server/models"
	"http-server/network"
	"net"
	"strconv"
	"strings"
)

var dao = database.GetDao()

func registerUserHandlers() {
	registerHandler(POST, "/user/create", createUser)
	registerHandler(GET, "/user/{id}", getUser)
}

func getUser(conn net.Conn, http models.HttpRequest, pathVars map[string]string) {
	// dao := database.GetDao()

	id, err := strconv.Atoi(pathVars["id"])

	if err != nil {

	}

	data := models.User{Id: id}
	user, err := dao.GetUserById(data.Id)
	var response string

	if err != nil {
		response = network.RESPONSE_BAD_REQUEST + network.CRLF + "not working"
	} else {
		userJson, _ := json.Marshal(user)
		response = network.RESPONSE_OK + network.CRLF + string(userJson)
	}

	network.SendData(response, conn)
}

func createUser(conn net.Conn, http models.HttpRequest, _ map[string]string) {
	// dao := database.GetDao()
	data := new(models.User)
	json.Unmarshal([]byte(http.Body), &data)
	var response string

	if err := dao.CreateUser(data.Username, data.Password); err != nil {
		var msg string
		usernameTaken := "Username already exists."
		passwordTooShort := "Password must be 6 or more characters."

		switch {
		case strings.HasPrefix(err.Error(), "UNIQUE"):
			msg = usernameTaken
		case strings.HasPrefix(err.Error(), "CHECK"):
			msg = passwordTooShort
		default:
			msg = network.RESPONSE_BAD_REQUEST
		}

		response = network.RESPONSE_BAD_REQUEST + network.CRLF + msg
	} else {
		response = network.RESPONSE_OK + network.CRLF + "Created user"
	}

	network.SendData(response, conn)
}

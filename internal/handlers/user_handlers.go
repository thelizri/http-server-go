package handlers

import (
	"encoding/json"
	"fmt"
	userrepository "http-server/internal/data/repositories/user"
	"http-server/internal/models"
	"http-server/internal/network"
	"net"
	"strconv"
)

var userRepository = userrepository.NewUserRepository()

func registerUserHandlers() {
	registerHandler(POST, "/users/create", createUser)
	registerHandler(GET, "/users/{id}", getUserByIdAsPathVariable)
	registerHandler(GET, "/users", getUserByIdAsQuery)
}

func getUserByIdAsQuery(conn net.Conn, http models.HttpRequest) {
	key := "id"
	id, err := strconv.Atoi(http.Query[key])
	var response string

	if err != nil {
		response = network.RESPONSE_BAD_REQUEST + network.CRLF + fmt.Sprintf("missing query key: %s", key)
		network.SendData(response, conn)
		return
	}

	data := models.User{Id: id}
	user, err := userRepository.GetUserById(data.Id)

	if err != nil {
		response = network.RESPONSE_BAD_REQUEST + network.CRLF + err.Error()
	} else {
		userJson, _ := json.Marshal(user)
		response = network.RESPONSE_OK + network.CRLF + string(userJson)
	}

	network.SendData(response, conn)
}

func getUserByIdAsPathVariable(conn net.Conn, http models.HttpRequest) {
	key := "id"
	id, err := strconv.Atoi(http.PathVariables[key])
	var response string

	if err != nil {
		response = network.RESPONSE_BAD_REQUEST + network.CRLF + fmt.Sprintf("missing path variable: %s", key)
		network.SendData(response, conn)
		return
	}

	data := models.User{Id: id}
	user, err := userRepository.GetUserById(data.Id)

	if err != nil {
		response = network.RESPONSE_BAD_REQUEST + network.CRLF + "not working"
	} else {
		userJson, _ := json.Marshal(user)
		response = network.RESPONSE_OK + network.CRLF + string(userJson)
	}

	network.SendData(response, conn)
}

func createUser(conn net.Conn, http models.HttpRequest) {
	// dao := database.GetDao()
	data := new(models.User)
	json.Unmarshal([]byte(http.Body), &data)
	var response string

	if err := userRepository.CreateUser(data.Username, data.Password); err != nil {

		response = network.RESPONSE_BAD_REQUEST + network.CRLF + err.Error()
	} else {
		response = network.RESPONSE_OK + network.CRLF + "Created user"
	}

	network.SendData(response, conn)
}

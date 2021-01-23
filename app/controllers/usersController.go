package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kenshin579/analyzing-restful-api-golang-jwt-mysql/app/models"
	"github.com/kenshin579/analyzing-restful-api-golang-jwt-mysql/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Get one user by id
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		utils.Respond(w, utils.Message(false, "There was an error in your request"))
		return
	}

	user := models.GetUser(id)
	if user == nil {
		utils.Respond(w, utils.Message(false, "User not found"))
		return
	}

	resp := utils.Message(true, "success")
	resp["data"] = user
	utils.Respond(w, resp)
	return
}

// Get all the users in the users table
func GetUsers(w http.ResponseWriter, r *http.Request) {
	resp := utils.Message(true, "success")
	users := models.GetUsers()
	if users == nil {
		utils.Respond(w, utils.Message(false, "No users found"))
		return
	}
	resp["data"] = users
	utils.Respond(w, resp)
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	fmt.Println("user.Name", user.Name)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	defer r.Body.Close()

	resp := user.Create()
	utils.Respond(w, resp)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user models.User
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.Respond(w, utils.Message(false, "There was an error in your request"))
		return
	}

	err = models.GetUserForUpdateOrDelete(id, &user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "User not found"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	user.ID = uint(id)
	user.UpdatedAt = time.Now().Local()
	defer r.Body.Close()

	// Update user here
	err = models.UpdateUser(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Could not update the record"))
		return
	}
	resp := utils.Message(true, "Updated successfully")
	resp["data"] = user
	utils.Respond(w, resp)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user models.User
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "There was an error in your request"))
		return
	}

	err = models.GetUserForUpdateOrDelete(id, &user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "User not found"))
		return
	}

	err = models.DeleteUser(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Could not delete the record"))
		return
	}
	utils.Respond(w, utils.Message(true, "User has been deleted successfully"))
	return
}

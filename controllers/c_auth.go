package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dimaskiddo/frame-go/utils"
)

// GetAuth Function to Get Authorization Token
func GetAuth(w http.ResponseWriter, r *http.Request) {
	var creds utils.BasicCredentials

	// Decode JSON from Request Body to User Data
	// Use _ As Temporary Variable
	_ = json.NewDecoder(r.Body).Decode(&creds)

	// Make Sure Username and Password is Not Empty
	if len(creds.Username) == 0 || len(creds.Password) == 0 {
		utils.ResponseBadRequest(w, "Invalid authorization")
		log.Println("Invalid authorization")
		return
	}

	// Some Business Logic Here to Match The Username and Password
	if creds.Username == "user" && creds.Password == "password" {
		// Get JWT Token From Pre-Defined Function
		token, err := utils.GetJWTToken(creds.Username)
		if err != nil {
			utils.ResponseInternalError(w, err.Error())
			log.Println(err.Error())
		} else {
			var response utils.JWTResponse

			response.Status = true
			response.Code = http.StatusOK
			response.Token = token

			utils.ResponseWrite(w, response.Code, response)
		}
	} else {
		utils.ResponseUnauthorized(w)
	}
}

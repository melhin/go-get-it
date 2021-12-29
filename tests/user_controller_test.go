package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gorm.io/gorm"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	RouteDbTransactionWrapper(t, func(t *testing.T, tx *gorm.DB, testApp *TestRequestRoute) {

		// Signup a new user on a fresh database
		var signup_request = []byte(`{"username":"test_username1", "name": "somename", "password": "password1"}`)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(signup_request))
		req.Header.Set("Content-Type", "application/json")
		response := testApp.ExecuteRequest(req, tx)
		assert.Equal(t, http.StatusCreated, response.Code)

		// Login the user and get user token
		var login_request_body = []byte(`{"username":"test_username1", "password": "password1"}`)
		login_req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(login_request_body))
		req.Header.Set("Content-Type", "application/json")
		login_response := testApp.ExecuteRequest(login_req, tx)
		assert.Equal(t, http.StatusOK, login_response.Code)
		var body map[string]interface{}
		json.Unmarshal(login_response.Body.Bytes(), &body)
		assert.NotNil(t, body["token"])

		// Use user token to get some information
		secret_req, _ := http.NewRequest("GET", "/secret", bytes.NewBuffer([]byte(``)))
		secret_req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cast.ToString(body["token"])))
		secret_response := testApp.ExecuteRequest(secret_req, tx)
		assert.Equal(t, http.StatusOK, secret_response.Code)
		var secret_body map[string]interface{}
		json.Unmarshal(secret_response.Body.Bytes(), &secret_body)
		assert.Equal(t, secret_body["user_id"], float64(1))
		assert.Equal(t, secret_body["authorized"], true)
		assert.IsType(t, secret_body["expires"], float64(1))

	})

}

func TestLoginFailure(t *testing.T) {
	RouteDbTransactionWrapper(t, func(t *testing.T, tx *gorm.DB, testApp *TestRequestRoute) {
		createMultipleUsers(t, tx)

		// Wrong password check
		var mismatch = []byte(`{"username":"test_username1", "password": "mismatch_password"}`)
		mis_req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(mismatch))
		mis_req.Header.Set("Content-Type", "application/json")
		mis_response := testApp.ExecuteRequest(mis_req, tx)
		assert.Equal(t, http.StatusBadRequest, mis_response.Code)
		var mismatch_body map[string]interface{}
		json.Unmarshal(mis_response.Body.Bytes(), &mismatch_body)
		assert.Equal(t, mismatch_body["message"], "Invalid username or password")

		// Invalid Login user
		var invalid = []byte(`{"username":"unknown_user", "password": "mismatch_password"}`)
		invalid_req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(invalid))
		invalid_req.Header.Set("Content-Type", "application/json")
		invalid_response := testApp.ExecuteRequest(invalid_req, tx)
		assert.Equal(t, http.StatusBadRequest, invalid_response.Code)
		var invalid_response_body map[string]interface{}
		json.Unmarshal(mis_response.Body.Bytes(), &invalid_response_body)
		assert.Equal(t, invalid_response_body["message"], "Invalid username or password")

		// Password not given
		missing, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username":"unknown_user"}`)))
		missing_response := testApp.ExecuteRequest(missing, tx)
		assert.Equal(t, http.StatusBadRequest, missing_response.Code)
	})

}

func TestGetSecretFailure(t *testing.T) {
	RouteDbTransactionWrapper(t, func(t *testing.T, tx *gorm.DB, testApp *TestRequestRoute) {
		createMultipleUsers(t, tx)

		// No Auth token
		req, _ := http.NewRequest("GET", "/secret", nil)
		response := testApp.ExecuteRequest(req, tx)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		var body map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &body)
		assert.Equal(t, "Unauthorized", body["message"])

		// Wrong Bearer Token
		wrong_req, _ := http.NewRequest("GET", "/secret", nil)
		wrong_req.Header.Set("Authorization", "Bearer Wrong Token")
		wrong_response := testApp.ExecuteRequest(wrong_req, tx)
		assert.Equal(t, http.StatusUnauthorized, wrong_response.Code)
		var wrong_body map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &wrong_body)
		assert.Equal(t, "Unauthorized", wrong_body["message"])
	})

}

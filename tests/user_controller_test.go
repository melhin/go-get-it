package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	RouteDbTransactionWrapper(t, func(t *testing.T, tx *gorm.DB, testApp *TestRequestRoute) {
		var signup_request = []byte(`{"username":"test_username1", "name": "somename", "password": "password1"}`)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(signup_request))
		req.Header.Set("Content-Type", "application/json")
		response := testApp.ExecuteRequest(req, tx)
		assert.Equal(t, http.StatusCreated, response.Code)
		var login_request_body = []byte(`{"username":"test_username1", "password": "password1"}`)
		login_req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(login_request_body))
		req.Header.Set("Content-Type", "application/json")
		login_response := testApp.ExecuteRequest(login_req, tx)
		assert.Equal(t, http.StatusOK, login_response.Code)
		var body map[string]interface{}
		json.Unmarshal(login_response.Body.Bytes(), &body)
		assert.NotNil(t, body["token"])
	})

}

func TestLogin(t *testing.T) {
	RouteDbTransactionWrapper(t, func(t *testing.T, tx *gorm.DB, testApp *TestRequestRoute) {
		createMultipleUsers(t, tx)
		var correct = []byte(`{"username":"test_username1", "password": "password1"}`)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(correct))
		req.Header.Set("Content-Type", "application/json")
		response := testApp.ExecuteRequest(req, tx)
		assert.Equal(t, http.StatusOK, response.Code)

		var mismatch = []byte(`{"username":"test_username1", "password": "mismatch_password"}`)
		mis_req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(mismatch))
		mis_req.Header.Set("Content-Type", "application/json")
		mis_response := testApp.ExecuteRequest(mis_req, tx)
		assert.Equal(t, http.StatusBadRequest, mis_response.Code)
		var mismatch_body map[string]interface{}
		json.Unmarshal(mis_response.Body.Bytes(), &mismatch_body)
		assert.Equal(t, mismatch_body["message"], "Invalid username or password")
	})

}

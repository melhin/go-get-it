package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestLogin(t *testing.T) {
	db := Connect()
	tx := db.Begin()
	defer tx.Rollback()
	defer CloseDB(db)
	createMultipleUsers(t, tx)

	var correct = []byte(`{"username":"test_username1", "password": "password1"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(correct))
	req.Header.Set("Content-Type", "application/json")
	response := ExecuteRequest(req, tx)

	CheckResponseCode(t, http.StatusOK, response.Code)

	var mismatch = []byte(`{"username":"test_username1", "password": "mismatch_password"}`)
	mis_req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(mismatch))
	mis_req.Header.Set("Content-Type", "application/json")
	mis_response := ExecuteRequest(mis_req, tx)
	CheckResponseCode(t, http.StatusBadRequest, mis_response.Code)
	var mismatch_body map[string]interface{}
	json.Unmarshal(mis_response.Body.Bytes(), &mismatch_body)

	assert.Equal(t, mismatch_body["message"], "Invalid username or password")

}

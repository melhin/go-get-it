package tests

import (
	"go-get-it/app"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var MainApp = &app.App{}

func init() {
	MainApp.Router = mux.NewRouter()
	MainApp.InitializeRoutes()
}

func ExecuteRequest(req *http.Request, tx *gorm.DB) *httptest.ResponseRecorder {
	MainApp.DB = tx
	rr := httptest.NewRecorder()
	MainApp.Router.ServeHTTP(rr, req)

	return rr
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

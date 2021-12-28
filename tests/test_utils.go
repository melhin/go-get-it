package tests

import (
	"go-get-it/app"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TestRequestRoute struct {
	App app.App
}

func (Ta *TestRequestRoute) Initialize() {
	Ta.App.Router = mux.NewRouter()
	Ta.App.InitializeRoutes()
}

func (Ta *TestRequestRoute) ExecuteRequest(req *http.Request, tx *gorm.DB) *httptest.ResponseRecorder {
	Ta.App.DB = tx
	rr := httptest.NewRecorder()
	Ta.App.Router.ServeHTTP(rr, req)
	return rr
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func DbTransactionWrapper(t *testing.T, wrapper_func func(*testing.T, *gorm.DB)) {
	db := Connect()
	tx := db.Begin()
	defer tx.Rollback()
	defer CloseDB(db)
	wrapper_func(t, tx)

}

func RouteDbTransactionWrapper(t *testing.T, wrapper_func func(*testing.T, *gorm.DB, *TestRequestRoute)) {
	db := Connect()
	tx := db.Begin()
	defer tx.Rollback()
	defer CloseDB(db)
	testApp := &TestRequestRoute{}
	testApp.Initialize()
	wrapper_func(t, tx, testApp)
}

package sj

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/mgo.v2"
)

const (
	TestSessionsColl = "Session"
	TestDBURL        = "mongodb://127.0.0.1:27017"
	TestDBName       = "testing-db"
)

type (
	TestRequest struct {
		Body    string
		Handler http.Handler
		Header  http.Header
	}
)

func DialTestDB(t *testing.T) (*mgo.Session, *mgo.Database) {
	session, err := mgo.Dial(TestDBURL)
	if err != nil {
		t.Fatal(err)

	}

	db := session.DB(TestDBName)

	return session, db

}

func CleanTestDB(s *mgo.Session, db *mgo.Database, t *testing.T) {
	err := db.DropDatabase()
	if err != nil {
		t.Fatal(err)
	}
	s.Close()
}

func (t *TestRequest) SendWithToken(method, path, token string) *httptest.ResponseRecorder {
	reqData := *t
	body := bytes.NewBufferString(reqData.Body)
	reqData.Header.Add("X-XSRF-TOKEN", token)

	req, _ := http.NewRequest(method, path, body)
	req.Header = reqData.Header
	w := httptest.NewRecorder()
	reqData.Handler.ServeHTTP(w, req)
	*t = reqData
	return w
}

func (t *TestRequest) Send(method, path string) *httptest.ResponseRecorder {
	reqData := *t
	body := bytes.NewBufferString(reqData.Body)

	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	reqData.Handler.ServeHTTP(w, req)
	*t = reqData
	return w
}

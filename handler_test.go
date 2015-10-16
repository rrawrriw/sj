package sj

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/rrawrriw/angular-sauth-handler"
)

func NewTestSession(user, token string, db *mgo.Database, t *testing.T) aauth.Session {
	coll := db.C(TestSessionsColl)
	expires := time.Now().AddDate(0, 0, 1)
	session := aauth.Session{
		Token:   token,
		UserID:  user,
		Expires: expires,
	}
	err := coll.Insert(session)
	if err != nil {
		t.Fatal(err)
	}

	return session
}

// Erzeuge standard Benutzer und Serien in der Datenbank für Testzwecke
func NewTestDBEnv(t *testing.T, db *mgo.Database) (User, aauth.Session, SeriesList) {
	userName := "greatLover99"

	series1 := Series{
		Title: "Narcos",
		Image: Resource{
			"kinox.to",
			"http://kinox.to/Stream/Narcos.html",
		},
		Episodes: Resource{
			"imdb.com",
			"http://www.imdb.com/title/tt2707408",
		},
		Desc: Resource{
			"imdb.com",
			"http://www.imdb.com/title/tt2707408",
		},
		Portal: Resource{
			"kinox.to",
			"http://kinox.to/Stream/Narcos.html",
		},
	}

	id1, err := NewSeries(db, series1)
	if err != nil {
		t.Fatal(err)
	}
	series1.ID = id1

	series2 := Series{
		Title: "Mr. Robot",
		Image: Resource{
			"imdb.com",
			"http://www.imdb.com/title/tt4158110/",
		},
		Episodes: Resource{
			"serienjunkie.de",
			"http://www.serienjunkies.de/mr-robot/",
		},
		Desc: Resource{
			"serienjunkie.de",
			"http://www.serienjunkies.de/mr-robot/",
		},
		Portal: Resource{
			"kinox.to",
			"http://kinox.to/Stream/Mr-Robot.html",
		},
	}

	id2, err := NewSeries(db, series2)
	if err != nil {
		t.Fatal(err)
	}
	series2.ID = id2

	ids := []bson.ObjectId{
		id1,
		id2,
	}

	user := User{
		Name:   userName,
		Series: ids,
	}

	uID, err := NewUser(db, user)
	if err != nil {
		t.Fatal(err)
	}
	user.Id = uID

	userToken := "123"
	session := NewTestSession(string(uID.Hex()), userToken, db, t)

	sList := SeriesList{
		series1,
		series2,
	}

	return user, session, sList
}

func ParseSuccessResponse(b *bytes.Buffer) (SuccessResponse, error) {
	resp := SuccessResponse{}

	err := json.Unmarshal(b.Bytes(), &resp)
	if err != nil {
		return SuccessResponse{}, err
	}

	return resp, nil
}

func EqualSuccessResponse(r1 SuccessResponse, b *bytes.Buffer, cFun func(SuccessResponse, SuccessResponse) bool) bool {

	r2, err := ParseSuccessResponse(b)
	if err != nil {
		return false
	}

	if r1.Status == r2.Status && cFun(r1, r2) {
		return true
	}

	return false
}

func ParseSeriesList(r SuccessResponse) (SeriesList, error) {
	cErr := errors.New("type convert error")
	tmp, ok := r.Data.([]interface{})
	if !ok {
		return SeriesList{}, cErr
	}

	sL2 := SeriesList{}
	for _, e := range tmp {
		tmp, ok := e.(map[string]interface{})
		if !ok {
			return SeriesList{}, cErr
		}

		fName := []string{"Image", "Desc", "Episodes", "Portal"}
		fields := map[string]Resource{}

		for _, name := range fName {
			v, ok := tmp[name].(map[string]interface{})
			if !ok {
				return SeriesList{}, cErr
			}
			r := Resource{
				Name: v["Name"].(string),
				URL:  v["URL"].(string),
			}
			fields[name] = r
		}

		series := Series{
			Title:    tmp["Title"].(string),
			Image:    fields["Image"],
			Desc:     fields["Desc"],
			Episodes: fields["Episodes"],
			Portal:   fields["Portal"],
		}
		sL2 = append(sL2, series)
	}

	return sL2, nil
}

func EqualSeriesList(r1, r2 SuccessResponse) bool {
	sL1, ok := r1.Data.(SeriesList)
	if !ok {
		return false
	}

	sL2, err := ParseSeriesList(r2)
	if err != nil {
		return false
	}

	sort.Sort(sL1)
	sort.Sort(sL2)

	for i, s := range sL1 {
		if !EqualSeries(s, sL2[i]) {
			return false
		}
	}

	return true
}

func EqualUserBasic(u1, u2 User) error {
	if u1.Name != u2.Name {
		m := fmt.Sprintf("Expect %v was %v", u1.Name, u2.Name)
		return errors.New(m)
	}

	if u1.Pass != u2.Pass {
		m := fmt.Sprintf("Expect %v was %v", u1.Pass, u2.Pass)
		return errors.New(m)
	}

	return nil

}

func EqualFailResponse(r1 *bytes.Buffer, r2 FailResponse) error {
	fr := FailResponse{}
	err := json.Unmarshal(r1.Bytes(), &fr)
	if err != nil {
		return err
	}

	if fr.Status != r2.Status {
		m := fmt.Sprintf("Expect fail response")
		return errors.New(m)
	}

	if fr.Err != r2.Err {
		m := fmt.Sprintf("Expect %v was %v", r2.Err, fr.Err)
		return errors.New(m)
	}

	return nil
}

func ExistsIDField(r1, r2 SuccessResponse) bool {
	id2, ok := r2.Data.(map[string]interface{})
	if !ok {
		return false
	}

	_, ok = id2["ID"]
	if !ok {
		return false
	}

	return true
}

func NewTestApp(t *testing.T) AppCtx {
	s, _ := DialTestDB(t)
	specs := Specs{
		DBName: TestDBName,
		DBURL:  TestDBURL,
	}

	ctx := AppCtx{
		MgoSession: s,
		Specs:      specs,
		Mutex:      &sync.Mutex{},
	}

	return ctx
}

func Test_GET_SeriesOfUser_OK(t *testing.T) {
	app := NewTestApp(t)
	db := app.DB()
	defer CleanTestDB(app.MgoSession, db, t)

	user, session, sList := NewTestDBEnv(t, db)
	auth := aauth.AngularAuth(db, TestSessionsColl)

	handler := gin.New()
	req := TestRequest{
		Body:    "",
		Header:  http.Header{},
		Handler: handler,
	}

	h := NewAppHandler(ReadSeriesOfUserHandler, app)
	handler.GET("/:id", auth, h)

	url := fmt.Sprintf("/%v", user.ID)
	resp := req.SendWithToken("GET", url, session.Token)

	if resp.Code != http.StatusOK {
		t.Fatal("Expect http-status", http.StatusOK, "was", resp.Code)
	}

	expectResult := NewSuccessResponse(sList)

	r := EqualSuccessResponse(expectResult, resp.Body, EqualSeriesList)
	if !r {
		t.Fatal("Expect", expectResult, "was", resp.Body)
	}

}

func Test_POST_Series_OK(t *testing.T) {
	app := NewTestApp(t)
	db := app.DB()
	defer CleanTestDB(app.MgoSession, db, t)

	_, session, _ := NewTestDBEnv(t, db)

	auth := aauth.AngularAuth(db, TestSessionsColl)

	body := `
	{
		"Data": {
			"Title": "Elementary",
			"Image": {
				"Name": "kinox.to",
				"URL": "http://kinox.to/1"
			},
			"Desc": {
				"Name": "kinox.to",
				"URL": "http://kinox.to/1"
			},
			"Episodes": {
				"Name": "kinox.to",
				"URL": "http://kinox.to/1"
			},
			"Portal": {
				"Name": "kinox.to",
				"URL": "http://kinox.to/1"
			}
		}
	}`

	handler := gin.New()
	req := TestRequest{
		Body:    body,
		Header:  http.Header{},
		Handler: handler,
	}

	h := NewAppHandler(NewSeriesHandler, app)
	handler.POST("/", auth, h)

	resp := req.SendWithToken("POST", "/", session.Token)

	if resp.Code != http.StatusOK {
		t.Fatal("Expect", http.StatusOK, "was", resp.Code)
	}

	expectResult := NewSuccessResponse(nil)
	r := EqualSuccessResponse(expectResult, resp.Body, ExistsIDField)
	if !r {
		t.Fatal("Expect", expectResult, "was", resp.Body)
	}
}

func Test_ParseNewSeriesRequest_FailMissingFields(t *testing.T) {
	data := `
	{
		"Data": {
			"Title": "Title"
		}
	}`

	body := bytes.NewReader([]byte(data))
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	ginCtx := gin.Context{}
	ginCtx.Request = req

	_, err = ParseNewSeriesRequest(&ginCtx)
	expect := NewMissingFieldError("Image")
	if err.Error() != expect.Error() {
		t.Fatal("Expect", expect, "was", err)
	}
}

func Test_POST_NewUser_OK(t *testing.T) {
	app := NewTestApp(t)
	db := app.DB()
	defer CleanTestDB(app.MgoSession, db, t)

	handler := gin.New()
	newUserBody := `
	{
		"Data": {
			"Name": "machine_XXX",
			"Password": "love!"
		}
	}
	`
	req := TestRequest{
		Body:    newUserBody,
		Header:  http.Header{},
		Handler: handler,
	}

	h := NewAppHandler(NewUserHandler, app)
	handler.POST("/", h)

	resp := req.Send("POST", "/")

	expectResult := NewSuccessResponse(nil)
	r := EqualSuccessResponse(expectResult, resp.Body, ExistsIDField)
	if !r {
		t.Fatal("Expect success response with id field was", resp.Body)
	}

}

func Test_POST_NewUser_FailUserExists(t *testing.T) {
	app := NewTestApp(t)
	db := app.DB()
	defer CleanTestDB(app.MgoSession, db, t)

	user := User{
		Name: "machine_XXX",
		Pass: "love!",
	}
	_, err := NewUser(db, user)
	if err != nil {
		t.Fatal(err)
	}

	handler := gin.New()
	newUserBody := `
	{
		"Data": {
			"Name": "machine_XXX",
			"Password": "love!"
		}
	}
	`
	req := TestRequest{
		Body:    newUserBody,
		Header:  http.Header{},
		Handler: handler,
	}

	h := NewAppHandler(NewUserHandler, app)
	handler.POST("/", h)

	resp := req.Send("POST", "/")

	expectResp := NewFailResponse(UserExistsError)
	err = EqualFailResponse(resp.Body, expectResp)
	if err != nil {
		t.Fatal(err)
	}

}

func Test_ParseNewUserRequest_OK(t *testing.T) {
	data := `
	{
		"Data": {
			"Name": "A",
			"Password": "a"
		}
	}`

	body := bytes.NewReader([]byte(data))
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	expectUser := User{
		Name: "A",
		Pass: "a",
	}

	user, err := ParseNewUserRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	err = EqualUserBasic(expectUser, user)
	if err != nil {
		t.Fatal(err)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rrawrriw/angular-auth"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Specs struct {
		DBName string
		DBURL  string
	}

	SuccessResponse struct {
		Status string
		Data   interface{}
	}

	FailResponse struct {
		Status string
		Err    string
	}

	JSONRequest struct {
		Data interface{}
	}

	AppContext interface {
		DB() *mgo.Database
	}

	AppCtx struct {
		Mutex      *sync.Mutex
		MgoSession *mgo.Session
		Specs      Specs
	}

	AppHandler func(*gin.Context, AppContext) error

	IDData struct {
		ID string
	}
)

func (app AppCtx) DB() *mgo.Database {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	sCopy := app.MgoSession.Copy()

	return sCopy.DB(app.Specs.DBName)
}

func NewApp() (AppCtx, error) {
	specs := Specs{}
	err := envconfig.Process(AppNamePrefix, &specs)
	if err != nil {
		return AppCtx{}, err
	}

	url := specs.DBURL
	session, err := mgo.Dial(url)
	if err != nil {
		return AppCtx{}, err
	}

	ctx := AppCtx{
		MgoSession: session,
		Specs:      specs,
	}

	return ctx, nil
}

func NewSuccessResponse(d interface{}) SuccessResponse {
	resp := SuccessResponse{
		Status: "success",
		Data:   d,
	}

	return resp
}

func NewFailResponse(err error) FailResponse {
	resp := FailResponse{
		Status: "fail",
		Err:    fmt.Sprintf("%v", err),
	}

	return resp
}

func ExistsFields(s map[string]interface{}, f []string) bool {
	for _, v := range f {
		_, ok := s[v]
		if !ok {
			return false
		}
	}

	return true
}

func ExportResource(s map[string]interface{}, key string) (Resource, bool) {
	v, ok := s[key].(map[string]interface{})
	if !ok {
		return Resource{}, false
	}

	name, ok := v["Name"].(string)
	if !ok {
		return Resource{}, false
	}

	url, ok := v["URL"].(string)
	if !ok {
		return Resource{}, false
	}

	r := Resource{
		Name: name,
		URL:  url,
	}

	return r, true

}

func ParseSeriesRequest(c *gin.Context) (Series, error) {
	reqErr := errors.New("Request error")

	buf := bytes.NewBuffer([]byte{})
	_, err := buf.ReadFrom(c.Request.Body)

	req := JSONRequest{}
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		return Series{}, err
	}

	data, ok := req.Data.(interface{})
	if !ok {
		return Series{}, reqErr
	}

	m, ok := data.(map[string]interface{})
	if !ok {
		return Series{}, reqErr
	}

	fields := []string{
		"Title",
		"Image",
		"Desc",
		"Episodes",
		"Portal",
	}
	ok = ExistsFields(m, fields)
	if !ok {
		return Series{}, reqErr
	}

	title, ok := m["Title"].(string)
	if !ok {
		return Series{}, errors.New("Title is missing")
	}

	resources := map[string]Resource{}
	resourceNames := []string{
		"Image",
		"Desc",
		"Episodes",
		"Portal",
	}
	for _, key := range resourceNames {
		v, ok := ExportResource(m, key)
		if !ok {
			errMsg := fmt.Sprintf("%v is missing", key)
			return Series{}, errors.New(errMsg)
		}
		resources[key] = v
	}

	series := Series{
		Title:    title,
		Image:    resources["Image"],
		Desc:     resources["Desc"],
		Episodes: resources["Episodes"],
		Portal:   resources["Portal"],
	}

	return series, nil
}

func ContextErrorDeco(h AppHandler, app AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c, app)
		if err != nil {
			resp := NewFailResponse(err)
			c.JSON(http.StatusOK, resp)
			return
		}
	}
}

func NewAppHandler(h AppHandler, app AppContext) gin.HandlerFunc {
	wrap := ContextErrorDeco(h, app)
	return wrap
}

func EmptyResource(r Resource) bool {
	if r.Name == "" && r.URL == "" {
		return true
	}

	return false
}

func EmptySeries(s Series) bool {
	if s.Title == "" &&
		EmptyResource(s.Image) &&
		EmptyResource(s.Desc) &&
		EmptyResource(s.Episodes) &&
		EmptyResource(s.Portal) {
		return true
	}

	return false
}

func NewSeriesHandler(c *gin.Context, app AppContext) error {
	series, err := ParseSeriesRequest(c)
	if err != nil {
		return err
	}

	if EmptySeries(series) {
		return errors.New("Wrong request")
	}

	id, err := NewSeries(app.DB(), series)
	if err != nil {
		return err
	}

	data := IDData{
		ID: string(id.Hex()),
	}
	c.JSON(http.StatusOK, NewSuccessResponse(data))
	return nil
}

func ReadSeriesOfUserHandler(c *gin.Context, app AppContext) error {
	tmp, err := c.Get("Session")
	if err != nil {
		return err
	}

	s, ok := tmp.(aauth.Session)
	if !ok {
		errors.New("Cannot find session")
	}

	userID := c.Params.ByName("id")
	if s.UserID == userID {
		errors.New("Wrong user")
	}

	db := app.DB()
	sList, err := ReadSeriesOfUser(db, bson.ObjectIdHex(s.UserID))

	if err != nil {
		return err
	}

	resp := NewSuccessResponse(sList)
	c.JSON(http.StatusOK, resp)

	return nil
}
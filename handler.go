package main

import (
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

	AppContext interface {
		DB() *mgo.Database
	}

	AppCtx struct {
		Mutex      *sync.Mutex
		MgoSession *mgo.Session
		Specs      Specs
	}

	AppHandler func(*gin.Context, AppContext) error
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

package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Resource struct {
		Name string
		URL  string
	}

	Series struct {
		Title    string
		Image    Resource
		Episodes Resource
		Desc     Resource
		Portal   Resource
	}

	User struct {
		Name   string
		Series []bson.ObjectId
	}

	ChangeUser struct {
		Name   string
		Series interface{}
	}

	// In der Zukunft ist es  möglich das man zum Beispiel
	// für das Feld Portal mehrer Resources hinterlegt kann.
	// Dann gibt es eine undefinierte Situationen beim
	// updaten soll nun ein Eintrag an ein Array angehängt
	// werden oder nicht. Daher dachte ich mir gibt es in
	// Zukunft ein Append Type. Je nachdem ob ein Feld ein
	// Append Type ist oder nicht kann das update wie gewünscht
	// durchgeführt werden.
	ChangeSeries struct {
		Title    string
		Image    Resource
		Episodes Resource
		Desc     Resource
		Portal   Resource
	}

	Episode struct {
		SeriesID bson.ObjectId
		Title    string
		Session  int
		Episode  int
		Watched  bool
	}

	AppendList []interface{}
	RemoveList []interface{}
)

func NewSeries(db *mgo.Database, series Series) (bson.ObjectId, error) {
	return bson.ObjectId(""), nil
}

func ReadSeries(db *mgo.Database, id bson.ObjectId) (Series, error) {
	return Series{}, nil
}

func ReadAllSeries(db *mgo.Database, sList []bson.ObjectId) ([]Series, error) {
	return []Series{}, nil
}

func UpdateSeries(db *mgo.Database, id bson.ObjectId, change ChangeSeries) error {
	return nil
}

func RemoveSeries(db *mgo.Database, id bson.ObjectId) error {
	return nil
}

func NewUser(db *mgo.Database, user User) (bson.ObjectId, error) {
	return bson.ObjectId(""), nil
}

func ReadUser(db *mgo.Database, id bson.ObjectId) (User, error) {
	return User{}, nil
}

func ReadSeriesFromUser(db *mgo.Database, user User) ([]Series, error) {
	return []Series{}, nil
}

func UpdateUser(db *mgo.Database, id bson.ObjectId, change ChangeUser) error {
	return nil
}

func RemoveUser(db *mgo.Database, id bson.ObjectId) error {
	return nil
}

func NewEpisode(db *mgo.Database, episode Episode) (bson.ObjectId, error) {
	return bson.ObjectId(""), nil
}

func NewEpisodeBatch(db *mgo.Database, episodes []Episode) ([]bson.ObjectId, error) {
	return []bson.ObjectId{}, nil
}

func ReadEpisode(db *mgo.Database, id bson.ObjectId) (Episode, error) {
	return Episode{}, nil
}

func ReadEpisodes(db *mgo.Database, id bson.ObjectId) ([]Episode, error) {
	return []Episode{}, nil
}

func WatchEpisode(db *mgo.Database, id bson.ObjectId) error {
	return nil
}

func ReadWatchedEpisodes(db *mgo.Database, id bson.ObjectId) ([]Episode, error) {
	return []Episode{}, nil
}

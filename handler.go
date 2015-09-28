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
)

func NewSeries(db *mgo.Database, series Series) (bson.ObjectId, error) {
	return bson.ObjectId(""), nil
}

func ReadSeries(db *mgo.Database, id bson.ObjectId) (Series, error) {
	return Series{}, nil
}

func UpdateSeries(db *mgo.Database, id bson.ObjectId, change ChangeSeries) error {
	return nil
}

func RemoveSeries(db *mgo.Database, id bson.ObjectId) error {
	return nil
}

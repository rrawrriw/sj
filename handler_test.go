package main

import (
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	TestDBURL  = "mongodb://127.0.0.1:2701"
	TestDBName = "testing-db"
)

func DialTest(t *testing.T) (*mgo.Session, *mgo.Database) {
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

func EqualResource(r1 Resource, r2 Resource) bool {
	if r1.Name == r2.Name && r1.URL == r2.URL {
		return true
	}

	return false
}

func EqualSeries(s1 Series, s2 Series) bool {
	if s1.Title != s2.Title &&
		!EqualResource(s1.Image, s2.Image) &&
		!EqualResource(s1.Episodes, s2.Episodes) &&
		!EqualResource(s1.Desc, s2.Desc) &&
		!EqualResource(s1.Portal, s2.Portal) {
		return true
	}

	return false
}

func Test_CRUDFuncSeries_OK(t *testing.T) {
	session, db := DialTest(t)
	defer CleanTestDB(session, db, t)

	series := Series{
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

	id, err := NewSeries(db, series)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ReadSeries(db, id)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualSeries(series, result) {
		t.Fatal("Expect", series, "was", result)
	}

	change := ChangeSeries{
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

	updatedSeries := Series{
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

	err = UpdateSeries(db, id, change)
	if err != nil {
		t.Fatal(err)
	}

	result, err = ReadSeries(db, id)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualSeries(updatedSeries, result) {
		t.Fatal("Expect", updatedSeries, "was", result)
	}

	err = RemoveSeries(db, id)
	if err != nil {
		t.Fatal(err)
	}

	result, err = ReadSeries(db, id)
	if err != mgo.ErrNotFound {
		t.Fatal(err)
	}

}

func Test_ReadAllSeries_OK(t *testing.T) {
}

func Test_CRUDFuncEpisode_OK(t *testing.T) {
	session, db := DialTest(t)
	defer CleanTestDB(session, db, t)

	seriesID := bson.ObjectId("123")

	episode := Episode{
		SeriesID: seriesID,
		Session:  1,
		Episode:  1,
		Title:    "Title",
		watched:  false,
	}

	id, err := NewEpisode(db, episode)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ReadEpisode(db, id)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualEpisode(episode, result) {
		t.Fatal("Expect", episode, "was", result)
	}

	err := WatchEpisode(db, seriesID, episode.Session, episode.Episode)
	if err != nil {
		t.Fatal(err)
	}

	updatedEpisode := Episode{
		SeriesID: seriesID,
		Session:  1,
		Episode:  1,
		Title:    "Title",
		Watched:  true,
	}

	result, err := ReadEpisode(db, id)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualEpisode(updatedEpisode, result) {
		t.Fatal("Expect", episode, "was", result)
	}

	episode2 := Episode{
		SeriesID: seriesID,
		Session:  1,
		Episode:  2,
		Title:    "Title 2",
		Watched:  true,
	}

	id2, err := NewEpisode(db, episode2)
	if err != nil {
		t.Fatal(err)
	}

	watchedEpisodes := []Episode{
		episode,
		episode2,
	}

	allWatchedEpisodes, err := ReadWatchedEpisodes(db)
	if err != nil {
		t.Fatal(err)
	}

	for i, e := range episodes {
		if !EqualEpisode(e, allWatchedEpisodes[i]) {
			t.Fatal("Expect", e, "was", allWatchedEpisodes[i])
		}
	}

}

func Test_NewBatch_OK(t *testing.T) {
}

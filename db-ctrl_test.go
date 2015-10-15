package sj

import (
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func EqualResource(r1 Resource, r2 Resource) bool {
	if r1.Name == r2.Name && r1.URL == r2.URL {
		return true
	}

	return false
}

func EqualEpisode(e1 Episode, e2 Episode) bool {
	if e1.Episode == e2.Episode &&
		e1.SeriesID == e2.SeriesID &&
		e1.Session == e2.Session &&
		e1.Title == e2.Title &&
		e1.Watched == e1.Watched {
		return true
	}

	return false
}

func EqualSeries(s1 Series, s2 Series) bool {
	if s1.Title == s2.Title &&
		EqualResource(s1.Image, s2.Image) &&
		EqualResource(s1.Episodes, s2.Episodes) &&
		EqualResource(s1.Desc, s2.Desc) &&
		EqualResource(s1.Portal, s2.Portal) {
		return true
	}

	return false
}

func ExistsID(ids []bson.ObjectId, id bson.ObjectId) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}

	return false
}

func EqualUser(u1 User, u2 User) bool {
	if u1.Name == u2.Name {
		for _, s := range u1.Series {
			if !ExistsID(u2.Series, s) {
				return false
			}
		}

		return true
	}

	return false
}

func Test_CRUDFuncSeries_OK(t *testing.T) {
	session, db := DialTestDB(t)
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

	change = ChangeSeries{
		Title: "True Detective",
	}

	updatedSeries = Series{
		Title: "True Detective",
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

	_, err = ReadSeries(db, id)
	if err != mgo.ErrNotFound {
		t.Fatal(err)
	}

}

func Test_ReadAllSeries_OK(t *testing.T) {
	session, db := DialTestDB(t)
	defer CleanTestDB(session, db, t)

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

	ids := []bson.ObjectId{
		id1,
		id2,
	}

	sList, err := ReadAllSeries(db, ids)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualSeries(series2, sList[0]) {
		t.Fatal("Expect", series2, "was", sList[0])
	}

	if !EqualSeries(series1, sList[1]) {
		t.Fatal("Expect", series1, "was", sList[1])
	}

	user := User{
		Name: "pimmel",
		Series: []bson.ObjectId{
			id1,
			id2,
		},
	}

	uID, err := NewUser(db, user)
	if err != nil {
		t.Fatal(err)
	}

	sList, err = ReadSeriesOfUser(db, uID)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualSeries(series2, sList[0]) {
		t.Fatal("Expect", series2, "was", sList[0])
	}

	if !EqualSeries(series1, sList[1]) {
		t.Fatal("Expect", series1, "was", sList[1])
	}

}

func Test_CRUDFuncUser_OK(t *testing.T) {
	session, db := DialTestDB(t)
	defer CleanTestDB(session, db, t)

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

	sID1, err := NewSeries(db, series1)

	user := User{
		Name: "Nase",
		Series: []bson.ObjectId{
			sID1,
		},
	}

	uID, err := NewUser(db, user)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ReadUser(db, uID)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualUser(user, result) {
		t.Fatal("Expect", user, "was", result)
	}

	result, err = FindUser(db, user.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualUser(user, result) {
		t.Fatal("Expect", user, "was", result)
	}

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

	sID2, err := NewSeries(db, series2)
	if err != nil {
		t.Fatal(err)
	}

	change := ChangeUser{
		Name: "Lang Nase",
		Series: AppendIDItems{
			sID2,
		},
	}

	err = UpdateUser(db, uID, change)
	if err != nil {
		t.Fatal(err)
	}

	updatedUser := User{
		Name: "Lang Nase",
		Series: []bson.ObjectId{
			sID2,
			sID1,
		},
	}

	result, err = ReadUser(db, uID)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualUser(updatedUser, result) {
		t.Fatal("Expect", updatedUser, "was", result)
	}

	change = ChangeUser{
		Series: RemoveIDItems{
			sID1,
		},
	}

	err = UpdateUser(db, uID, change)
	if err != nil {
		t.Fatal(err)
	}

	updatedUser = User{
		Name: "Lang Nase",
		Series: []bson.ObjectId{
			sID2,
		},
	}

	if !EqualUser(updatedUser, result) {
		t.Fatal("Expect", updatedUser, "was", result)
	}

	err = RemoveUser(db, uID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ReadUser(db, uID)
	if err != mgo.ErrNotFound {
		t.Fatal(err)
	}

}

func Test_CRUDFuncEpisode_OK(t *testing.T) {
	session, db := DialTestDB(t)
	defer CleanTestDB(session, db, t)

	seriesID := bson.NewObjectId()

	episode := Episode{
		SeriesID: seriesID,
		Session:  1,
		Episode:  1,
		Title:    "Title",
		Watched:  false,
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

	err = WatchEpisode(db, id)
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

	result, err = ReadEpisode(db, id)
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

	_, err = NewEpisode(db, episode2)
	if err != nil {
		t.Fatal(err)
	}

	watchedEpisodes := []Episode{
		episode,
		episode2,
	}

	allWatchedEpisodes, err := ReadWatchedEpisodes(db, seriesID)
	if err != nil {
		t.Fatal(err)
	}

	for i, e := range watchedEpisodes {
		if !EqualEpisode(e, allWatchedEpisodes[i]) {
			t.Fatal("Expect", e, "was", allWatchedEpisodes[i])
		}
	}

}

func Test_NewEpisodeBatch_OK(t *testing.T) {
	session, db := DialTestDB(t)
	defer CleanTestDB(session, db, t)

	seriesID := bson.NewObjectId()
	episode := Episode{
		SeriesID: seriesID,
		Session:  1,
		Episode:  1,
		Title:    "Title",
		Watched:  false,
	}

	episode2 := Episode{
		SeriesID: seriesID,
		Session:  1,
		Episode:  2,
		Title:    "Title 2",
		Watched:  true,
	}

	episodes := []Episode{
		episode,
		episode2,
	}
	_, err := NewEpisodeBatch(db, episodes)
	if err != nil {
		t.Fatal(err)
	}

	eResult, err := ReadEpisodes(db, seriesID)
	if err != nil {
		t.Fatal(err)
	}

	if !EqualEpisode(episode, eResult[0]) {
		t.Fatal("Expect", episode, "was", eResult[0])
	}

	if !EqualEpisode(episode2, eResult[1]) {
		t.Fatal("Expect", episode2, "was", eResult[2])
	}

}

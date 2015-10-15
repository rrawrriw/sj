package sj

import (
	"errors"
	"sort"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	SeriesColl  = "Series"
	UserColl    = "Users"
	EpisodeColl = "Episodes"
)

type (
	Resource struct {
		Name string `bson:"Name"`
		URL  string `bson:"URL"`
	}

	Series struct {
		ID       bson.ObjectId `bson:"_id,omitempty"`
		Title    string        `bson:"Title"`
		Image    Resource      `bson:"Image"`
		Episodes Resource      `bson:"Episodes"`
		Desc     Resource      `bson:"Desc"`
		Portal   Resource      `bson:"Portal"`
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

	SeriesList []Series

	User struct {
		// Noraml ID konflikt mit AngularSingIn API
		// Anfangsbuchstaben müssen groß sein da ansonst
		// die mgo Funktionen nicht mehr funktionieren
		// diese benötigen eine Öffentliche API sprich
		// Großbuchstaben
		Id     bson.ObjectId   `bson:"_id,omitempty"`
		Name   string          `bson:"Name"`
		Pass   string          `bson:"Password"`
		Series []bson.ObjectId `bson:"Series"`
	}

	ChangeUser struct {
		Name   string
		Series interface{}
	}

	Episode struct {
		ID       bson.ObjectId `bson:"_id,omitempty"`
		SeriesID bson.ObjectId `bson:"SeriesID"`
		Title    string        `bson:"Title"`
		Session  int           `bson:"Session"`
		Episode  int           `bson:"Episode"`
		Watched  bool          `bson:"Watched"`
	}

	Episodes []Episode

	AppendIDItems []bson.ObjectId
	RemoveIDItems []bson.ObjectId
)

// Define Sort List
func (l SeriesList) Len() int {
	return len(l)
}

func (l SeriesList) Less(x, y int) bool {
	return l[x].Title < l[y].Title
}

func (l SeriesList) Swap(x, y int) {
	l[x], l[y] = l[y], l[x]
}

func (l Episodes) Len() int {
	return len(l)
}

func (l Episodes) Less(x, y int) bool {
	if l[x].Session < l[y].Session {
		return true
	}

	if l[x].Session == l[y].Session {
		if l[x].Episode < l[y].Session {
			return true
		} else if l[x].Episode > l[y].Session {
			return false
		} else {
			return true
		}
	}

	return false
}

func (l Episodes) Swap(x, y int) {
	l[x], l[y] = l[y], l[x]
}

// Setup API for SiginHandler function
func (u User) ID() bson.ObjectId {
	return u.Id
}

func (u User) Password() string {
	return u.Pass
}

func NewSeries(db *mgo.Database, series Series) (bson.ObjectId, error) {
	coll := db.C(SeriesColl)

	id := bson.NewObjectId()
	series.ID = id
	err := coll.Insert(series)
	if err != nil {
		return bson.ObjectId(""), err
	}

	return id, nil
}

func ReadSeries(db *mgo.Database, id bson.ObjectId) (Series, error) {
	coll := db.C(SeriesColl)

	series := Series{}

	err := coll.FindId(id).One(&series)
	if err != nil {
		return Series{}, err
	}

	return series, nil
}

func ReadAllSeries(db *mgo.Database, sList []bson.ObjectId) ([]Series, error) {
	coll := db.C(SeriesColl)

	resultList := SeriesList{}

	err := coll.Find(bson.M{
		"_id": bson.M{
			"$in": sList,
		},
	}).All(&resultList)
	if err != nil {
		return []Series{}, err
	}
	sort.Sort(resultList)

	return resultList, nil
}

func ResourceEmpty(r Resource) bool {
	if r.Name == "" && r.URL == "" {
		return true
	}

	return false
}

func UpdateSeries(db *mgo.Database, id bson.ObjectId, change ChangeSeries) error {
	coll := db.C(SeriesColl)

	set := bson.M{}

	if change.Title != "" {
		set["Title"] = change.Title
	}

	if !ResourceEmpty(change.Image) {
		set["Image"] = change.Image
	}

	if !ResourceEmpty(change.Desc) {
		set["Desc"] = change.Desc
	}

	if !ResourceEmpty(change.Episodes) {
		set["Episodes"] = change.Episodes
	}

	if !ResourceEmpty(change.Portal) {
		set["Portal"] = change.Portal
	}

	update := bson.M{
		"$set": set,
	}

	mgoChange := mgo.Change{
		Update:    update,
		ReturnNew: false,
	}

	changeInfo, err := coll.FindId(id).Apply(mgoChange, nil)
	if err != nil {
		return err
	}

	if changeInfo.Updated != 1 {
		return errors.New("update error")
	}

	return nil
}

func RemoveSeries(db *mgo.Database, id bson.ObjectId) error {
	coll := db.C(SeriesColl)

	err := coll.RemoveId(id)
	if err != nil {
		return err
	}

	return nil
}

func NewUser(db *mgo.Database, user User) (bson.ObjectId, error) {
	coll := db.C(UserColl)

	id := bson.NewObjectId()
	user.Id = id

	err := coll.Insert(user)
	if err != nil {
		return bson.ObjectId(""), nil
	}

	return id, nil
}

func ReadUser(db *mgo.Database, id bson.ObjectId) (User, error) {
	coll := db.C(UserColl)

	user := User{}

	err := coll.FindId(id).One(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindUser(db *mgo.Database, name string) (User, error) {
	coll := db.C(UserColl)

	user := User{}

	query := bson.M{"Name": name}
	err := coll.Find(query).One(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func ReadSeriesOfUser(db *mgo.Database, id bson.ObjectId) ([]Series, error) {
	coll := db.C(UserColl)

	user := User{}
	err := coll.FindId(id).One(&user)
	if err != nil {
		return []Series{}, err
	}

	sList, err := ReadAllSeries(db, user.Series)
	if err != nil {
		return []Series{}, err
	}

	return sList, nil
}

func UpdateUser(db *mgo.Database, id bson.ObjectId, change ChangeUser) error {
	coll := db.C(UserColl)

	update := bson.M{}
	set := bson.M{}
	push := bson.M{}
	pull := bson.M{}

	if change.Name != "" {
		set["Name"] = change.Name
	}

	switch change.Series.(type) {
	case AppendIDItems:
		push["Series"] = bson.M{
			"$each": change.Series.(AppendIDItems),
		}
	case RemoveIDItems:
		pull["Series"] = bson.M{
			"$in": change.Series.(RemoveIDItems),
		}
	case []bson.ObjectId:
		set["Series"] = change.Series.([]bson.ObjectId)
	}

	if len(set) > 0 {
		update["$set"] = set
	}

	if len(push) > 0 {
		update["$push"] = push
	}

	if len(pull) > 0 {
		update["$pull"] = pull
	}

	mgoChange := mgo.Change{
		Update: update,
	}

	changeInfo, err := coll.FindId(id).Apply(mgoChange, nil)
	if err != nil {
		return err
	}

	if changeInfo.Updated != 1 {
		errors.New("update error")
	}

	return nil
}

func RemoveUser(db *mgo.Database, id bson.ObjectId) error {
	coll := db.C(UserColl)

	err := coll.RemoveId(id)
	if err != nil {
		return err
	}

	return nil
}

func NewEpisode(db *mgo.Database, episode Episode) (bson.ObjectId, error) {
	coll := db.C(EpisodeColl)

	id := bson.NewObjectId()
	episode.ID = id

	err := coll.Insert(episode)
	if err != nil {
		return bson.ObjectId(""), err
	}

	return id, nil
}

func NewEpisodeBatch(db *mgo.Database, episodes []Episode) ([]bson.ObjectId, error) {
	coll := db.C(EpisodeColl)

	ids := []bson.ObjectId{}
	inserts := []interface{}{}
	for _, e := range episodes {
		id := bson.NewObjectId()
		e.ID = id
		ids = append(ids, id)
		inserts = append(inserts, e)
	}

	err := coll.Insert(inserts...)
	if err != nil {
		return []bson.ObjectId{}, err
	}

	return ids, nil
}

func ReadEpisode(db *mgo.Database, id bson.ObjectId) (Episode, error) {
	coll := db.C(EpisodeColl)

	episode := Episode{}
	err := coll.FindId(id).One(&episode)
	if err != nil {
		return Episode{}, err
	}

	return episode, nil
}

func ReadEpisodes(db *mgo.Database, id bson.ObjectId) ([]Episode, error) {
	coll := db.C(EpisodeColl)

	result := []Episode{}
	query := bson.M{
		"SeriesID": id,
	}
	find := coll.Find(query)
	err := find.All(&result)
	if err != nil {
		return []Episode{}, err
	}

	return result, nil
}

func WatchEpisode(db *mgo.Database, id bson.ObjectId) error {
	coll := db.C(EpisodeColl)

	update := bson.M{
		"$set": bson.M{
			"Watched": true,
		},
	}
	mgoChange := mgo.Change{
		Update:    update,
		ReturnNew: false,
	}

	changeInfo, err := coll.FindId(id).Apply(mgoChange, nil)
	if err != nil {
		return err
	}

	if changeInfo.Updated != 1 {
		return errors.New("update error")
	}

	return nil
}

func ReadWatchedEpisodes(db *mgo.Database, id bson.ObjectId) (Episodes, error) {
	coll := db.C(EpisodeColl)

	result := Episodes{}
	query := bson.M{
		"SeriesID": id,
		"Watched":  true,
	}
	find := coll.Find(query)
	err := find.All(&result)
	if err != nil {
		return Episodes{}, err
	}

	sort.Sort(result)

	return result, nil
}

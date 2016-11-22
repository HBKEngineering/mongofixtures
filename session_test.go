package mongofixtures

import (
	"strconv"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type document struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
}

func checkCount(t *testing.T, collection *mgo.Collection, search string, count int, message string) {
	c, err := collection.Find(bson.M{"title": search}).Count()

	if err != nil {
		t.Fatal(err)
	}

	if c != count {
		t.Fatal(message + " : " + strconv.Itoa(c))
	}
}

func TestLoader(t *testing.T) {

	mongoSession, err := mgo.Dial("localhost")
	collection := mongoSession.DB("sample").C("collection1")

	session, err := Begin("localhost", "sample")
	defer session.End()

	if err != nil {
		t.Fatal(err)
	}

	type Test struct {
		Bla []int
	}

	err = session.Push("collection2", Test{Bla: []int{1, 2, 3, 4}})

	err = session.Push("collection1", document{Id: bson.NewObjectId(), Title: "This is a demo"})

	if err != nil {
		t.Fatal(err)
	}

	checkCount(t, collection, "This is a demo", 1, "Wrong count after inserting a document")

	err = session.Push("collection1", document{Id: bson.NewObjectId(), Title: "This is a demo 2"})

	if err != nil {
		t.Fatal(err)
	}

	checkCount(t, collection, "This is a demo 2", 1, "Wrong count after inserting a document")

	session.ImportYamlFile("test.yml")

	count, _ := mongoSession.DB("sample").C("employees").Find(nil).Count()
	if count != 3 {
		t.Fatal("Wrong count after inserting employees via yml")
	}

	// @todo check that Moss is friend with Roy
	// @todo check that Roy is friend with Moss
	// @todo check that Moss & Roy belongs to the basement

	err = session.Clean("collection1")
	err = session.Clean("employees")
	err = session.Clean("locations")

	if err != nil {
		t.Fatal(err)
	}

	checkCount(t, collection, "This is a demo", 0, "Wrong count after freeing the collection")
}

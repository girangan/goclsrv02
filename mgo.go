// modified sample from https://gist.github.com/border/3489566

package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
	"time"
	"os"
)

type Person struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Phone     string
	Timestamp time.Time
}

var (
	IsDrop = true
	trackcount = 0
)

func runTestDB() ([]Person, error) {
	trackcount++
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: []string{os.Getenv("MONGO_URL")},
		Username: "root",
		Password: os.Getenv("MONGODB_PASS"),
	})
	if err != nil {
		return nil, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Drop Database
	if IsDrop {
		err = session.DB("test").DropDatabase()
		if err != nil {
			return nil, err
		}
	}

	// Collection People
	c := session.DB("test").C("people")

	// Index
	index := mgo.Index{
		Key:        []string{"name", "phone"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = c.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	// Insert Datas
	err = c.Insert(&Person{Name: "Ale"+string(trackcount), Phone: "+55 53 1234 4321", Timestamp: time.Now()},
		&Person{Name: "Cla", Phone: "+66 33 1234 5678", Timestamp: time.Now()})

	if err != nil {
		return nil, err
	}

	// Query One
	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).Select(bson.M{"phone": 0}).One(&result)
	if err != nil {
		return nil, err
	}
	fmt.Println("Phone", result)

	// Query All
	var results []Person
	err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)

	if err != nil {
		return nil, err
	}
	fmt.Println("Results All: ", results)

	// Update
	colQuerier := bson.M{"name": "Ale"}
	change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": time.Now()}}
	err = c.Update(colQuerier, change)
	if err != nil {
		return nil, err
	}

	// Query All
	err = c.Find(nil).Sort("-timestamp").All(&results)

	if err != nil {
		return nil, err
	}
	fmt.Println("Results All: ", results)
	return results, nil
}
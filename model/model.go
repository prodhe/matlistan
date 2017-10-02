package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id       bson.ObjectId `bson:"_id"`
	Pid      bson.ObjectId // Profile: _id
	Username string
	Password string
}

type Profile struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string
}

type Session struct {
	Id            bson.ObjectId `bson:"_id"`
	Pid           bson.ObjectId // Profile: _id
	LastSeen      time.Time
	Authenticated bool
}

type Recipe struct {
	Id          bson.ObjectId `bson:"_id"`
	Userid      bson.ObjectId
	Title       string
	Category    string
	Ingredients []string
	Description string
}
package model

import "gopkg.in/mgo.v2/bson"

type Recipe struct {
	Id          bson.ObjectId `bson:"_id"`
	Userid      bson.ObjectId
	Title       string
	Category    string
	Ingredients []string
	Description string
}
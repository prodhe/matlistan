package model

import (
	"strconv"
	"strings"
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
	Name string        `json:"name"`
}

type Session struct {
	Id            bson.ObjectId `bson:"_id"`
	Pid           bson.ObjectId // Profile: _id
	LastSeen      time.Time
	Authenticated bool
}

type ingredient struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
	Name   string `json:"name"`
}

type Recipe struct {
	Id          bson.ObjectId `bson:"_id"`
	Pid         bson.ObjectId // Profile: _id
	Title       string        `json:"title"`
	Categories  []string      `json:"categories"`
	Ingredients []ingredient  `json:"ingredients"`
	Description string        `json:"description"`
}

// BreakIngredient will separate a string of input into amount, unit and name, if possible.
// Zero values otherwise and the input returned as name.
// Example: "1 dl water" -> amount: 1, unit: "dl", name: "water"
func (r *Recipe) BreakIngredients(s string) []ingredient {
	var ings = make([]ingredient, 0)

	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if !strings.ContainsAny(line, " ") {
			ings = append(ings, ingredient{0, "", line})
			continue
		}

		split := strings.SplitN(line, " ", 3)

		var amount int

		amount, err := strconv.Atoi(split[0])
		if err != nil {
			amount = 0
		}

		if len(split) != 3 && amount == 0 {
			ings = append(ings, ingredient{amount, "", line})
			continue
		}

		if len(split) != 3 {
			ings = append(ings, ingredient{amount, "", split[1]})
			continue
		}

		var unit string
		switch split[1] {
		case "tsk":
			fallthrough
		case "msk":
			fallthrough
		case "dl":
			fallthrough
		case "l":
			fallthrough
		case "g":
			fallthrough
		case "kg":
			fallthrough
		case "st":
			fallthrough
		case "pkt":
			unit = split[1]
		default:
			unit = ""
		}

		ings = append(ings, ingredient{amount, unit, split[2]})
	}

	return ings
}

// BreakTitle returns the main part of the title and a slice of every #hashtag found, without the leading # mark.
func (r *Recipe) BreakTitle(s string) (string, []string) {
	if !strings.ContainsAny(s, "#") {
		return s, nil
	}

	// index 0 is considered to be the title
	split := strings.Split(s, " #")

	return strings.TrimSpace(split[0]), split[1:]
}

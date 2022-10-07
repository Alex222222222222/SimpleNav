package main

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Link struct {
	// TODO: add links tags
	ID             int    // the unique id of a single link
	Name           string // the name of the link
	URL            string // the url
	Description    string // description set by user or fetch by the client
	IMG            string // link image provided by path
	FatherCategory int    // category
	Priority       int    // priority used in the ordering of links
}

type Category struct {
	ID             int    // unique id of a category
	Name           string // name of the category
	FatherCategory int    // father category, -1 implies no father category
	SubCategory    []int  // sub category
	Links          []int  // all links
	Description    string // description
	Hidden         bool   // hidden without the cookies?
	Priority       int    // priority used in the ordering of links
}

var Links map[int]*Link
var Categories []int
var AllCategories map[int]*Category

func LoadData() (err error) {

	var id int
	var name string
	var ft int
	var links string
	var description string
	var hidden bool
	var priority int
	var url string
	var img string
	var sb string

	// load category
	rows, err := db.Query(
		"SELECT * FROM Category",
	)
	if err != nil {
		err1 := rows.Close()
		if err1 != nil {
			return err1
		}

		return err
	}

	for rows.Next() {
		err = rows.Scan(&id, &name, &ft, &sb, &links, &description, &hidden, &priority)
		if err != nil {
			return err
		}

		var linkss []int
		err = json.Unmarshal([]byte(links), &linkss)
		if err != nil {
			return err
		}

		var sbs []int
		err = json.Unmarshal([]byte(sb), &sbs)
		if err != nil {
			return err
		}

		AllCategories[id] = &Category{
			ID:             id,
			Name:           name,
			FatherCategory: ft,
			SubCategory:    sbs,
			Links:          linkss,
			Description:    description,
			Hidden:         hidden,
			Priority:       priority,
		}

		Categories = append(Categories, id)
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	// load links
	rows, err = db.Query(
		"SELECT * FROM Links",
	)
	if err != nil {
		err1 := rows.Close()
		if err1 != nil {
			return err1
		}

		return err
	}

	for rows.Next() {
		err = rows.Scan(&id, &name, &url, &description, &img, &ft, &priority)
		if err != nil {
			return err
		}

		Links[id] = &Link{
			ID:             id,
			Name:           name,
			URL:            url,
			Description:    description,
			IMG:            img,
			FatherCategory: ft,
			Priority:       priority,
		}
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	return nil

}

func DBInit() (err error) {
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS Category (
			ID INT PRIMARY KEY,
			Name TEXT,
			FatherCategory INT,
			SubCategory TEXT,
			Links TEXT,
			Description TEXT,
			Hidden BOOLEAN,
			Priority int
		  );`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS Links (
			ID INT PRIMARY KEY,
			Name TEXT,
			URL TEXT,
			Description TEXT,
			IMG TEXT,
			FatherCategory INT,
			Priority int
		  );`,
	)
	if err != nil {
		return err
	}

	err = LoadData()
	if err != nil {
		return err
	}

	return nil
}

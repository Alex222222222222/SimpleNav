package main

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Link struct {
	ID             int      // the unique id of a single link
	Name           string   // the name of the link
	URL            string   // the url
	Description    string   // description set by user or fetch by the client
	IMG            string   // link image provided by path
	FatherCategory []int    // category
	Priority       int      // priority used in the ordering of links
	Tags           []string // tags of the links
}

type Category struct {
	ID             int    // unique id of a category
	Name           string // name of the category
	FatherCategory []int  // father category, -1 implies no father category
	SubCategory    []int  // sub category
	Links          []int  // all links
	Description    string // description
	Hidden         bool   // hidden without the cookies?
	Priority       int    // priority used in the ordering of links
}

// store all links
var Links map[int]*Link

// the order of the category for render
var Categories []int

// store all category
var AllCategories map[int]*Category

// load all category and links to ram for render purpose
func LoadData() (err error) {

	var id int
	var name string
	var ft []int
	var fts []byte
	var links []byte
	var description string
	var hidden bool
	var priority int
	var url string
	var img string
	var sb []byte
	var tagss []byte

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
		err = rows.Scan(&id, &name, &fts, &sb, &links, &description, &hidden, &priority)
		if err != nil {
			return err
		}

		var linkss []int
		err = json.Unmarshal(links, &linkss)
		if err != nil {
			return err
		}

		var sbs []int
		err = json.Unmarshal(sb, &sbs)
		if err != nil {
			return err
		}

		err = json.Unmarshal(fts, &ft)

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

	// TODO order category with priority at here
	// TODO the subcategory of category with a subcategory here

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
		err = rows.Scan(&id, &name, &url, &description, &img, &fts, &priority, &tagss)
		if err != nil {
			return err
		}

		var tags []string
		err = json.Unmarshal(tagss, &tags)
		if err != nil {
			return err
		}

		err = json.Unmarshal(fts, &ft)
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
			Tags:           tags,
		}
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	// TODO order sublinks or each category here

	return nil

}

// init the database
// create the *sql.DB
// create the table if the table does not exist
func DBInit() (err error) {
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS Category (
			ID INT PRIMARY KEY,
			Name TEXT,
			FatherCategory TEXT,
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
			FatherCategory TEXT,
			Priority int,
			Tags TEXT
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

// TODO add category
// insert a new category to database
// also insert itself to the list of subcategories of its father category
//
// fc means father category
// sb means sub category
// it contain a list of id of links
// d means description
func AddCategory(name string, fc []string, d string, hidden bool, priority int) (err error)

// TODO add links
// insert a new link to database and copy the img file to static/img/front/id.png
// also, automatically run the FetchLinkInformation func to fill missing information
// also insert itself to the list of links of its father category
//
// d means description
// img represent the given path of the img, if the img does not exist or the target path is not a img, then return a error
func AddLink(name string, url string, d string, img string, fc []string, priority int, tags []string) (err error)

// TODO
// fetch the target html and detect informations including names from the target html
// if imgRequired is true this will download the front img of the website and save it to the static/img/front/id.png
// id is used to store the img file, if the target file exist, this will replace it
func FetchLinkInformation(url string, id int, imgRequired bool) (l *Link, err error)

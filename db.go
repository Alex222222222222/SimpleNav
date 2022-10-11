package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

/*
CREATE TABLE IF NOT EXISTS Links (
	ID INT PRIMARY KEY,
	Name TEXT,
	URL TEXT,
	Description TEXT,
	IMG TEXT,
	FatherCategory TEXT,
	Priority int,
	Tags TEXT
);
*/
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

/*
CREATE TABLE IF NOT EXISTS Category (
	ID INT PRIMARY KEY,
	Name TEXT,
	FatherCategory TEXT,
	SubCategory TEXT,
	Links TEXT,
	Description TEXT,
	Hidden BOOLEAN,
	Priority int
);
*/
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

/*
CREATE TABLE IF NOT EXISTS Tags (
	ID INT PRIMARY KEY,
	Name TEXT,
	Links TEXT,
	Description TEXT,
	Hidden BOOLEAN,
	Priority int
);
*/
type Tag struct {
	ID          int    // unique id of a tag
	Name        string // name of the tag
	Links       []int  // all links
	Description string // description
	Hidden      bool   // hidden without the cookies?
	Priority    int    // priority used in the ordering of tags
}

// store all links
var Links map[int]*Link

// the order of the category for render
var Categories []int

// store all category
var AllCategories map[int]*Category

// TODO load all tags at the LoadData() func
// store a tags
var Tags map[int]*Link

// TODO have bug in json.UnMarshal, maybe caused by empty data in the database, need to deal with the empty data
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
		`CREATE TABLE IF NOT EXISTS Tags (
			ID INT PRIMARY KEY,
			Name TEXT,
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

	return nil
}

// add category
// insert a new category to database
// also insert itself to the list of subcategories of its father category
//
// fc means father category
// sb means sub category
// it contain a list of id of links
// d means description
func AddCategory(name string, fc []string, d string, hidden bool, priority int) (err error) {

	var id int

	err = db.QueryRow("SELECT ID FROM Category ORDER BY ID DESC").Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		id = 1
	} else if err != nil {
		return err
	}

	id += 1
	fcl := make([]int64, 0, len(fc))
	var fci int64
	for i := 0; i < len(fc); i += 1 {
		fct := strings.Split(fc[i], ":")
		for j := 0; j < len(fct); j += 1 {
			if len(fct[j]) > 0 {
				fci, err = strconv.ParseInt(fct[j], 10, 64)
				fcl = append(fcl, fci)
			}
		}
	}

	fcs, err := json.Marshal(fcl)

	_, err = db.Exec(
		`INSERT INTO CATEGORY (ID, Name, FatherCategory,Description, Hidden, Priority)
		VALUES (?, ?, ?, ?, ?, ?)`,
		id,
		name,
		fcs,
		d,
		hidden,
		priority,
	)
	if err != nil {
		return err
	}

	var sb []int64
	var sbs []byte
	var idj int64
	for i := 0; i < len(fcl); i += 1 {

		idj = fcl[i]
		err = db.QueryRow("SELECT SubCategory FROM Category WHERE ID = ?", idj).Scan(&sbs)
		if errors.Is(err, sql.ErrNoRows) {

		} else if err != nil {
			return err
		} else {

			if len(sbs) == 0 {
				sb = make([]int64, 0, 1)
			} else {
				err = json.Unmarshal(sbs, &sb)
				if err != nil {
					return err
				}
			}

			sb = append(sb, int64(id))
			sbs, err = json.Marshal(sb)
			if err != nil {
				return err
			}

			_, err = db.Exec("UPDATE Category SET SubCategory = ? WHERE ID = ?", sbs, idj)
			if err != nil {
				return err
			}

		}

	}

	return nil
}

// TODO need to be test
// add links
// insert a new link to database and copy the img file to static/img/front/id.png
// also, automatically run the FetchLinkInformation func to fill missing information
// also insert itself to the list of links of its father category
//
// d means description
// img represent the given path of the img, if the img does not exist or the target path is not a img, then return a error
func AddLink(name string, url string, d string, img string, fc []string, priority int, tags []string) (err error) {

	// get the id
	var id int
	err = db.QueryRow("SELECT ID FROM Links ORDER BY ID DESC").Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		id = 1
	} else if err != nil {
		return err
	}
	id += 1

	// parse the father categories
	fcl := make([]int64, 0, len(fc))
	var fci int64
	var fct []string
	for i := 0; i < len(fc); i += 1 {
		fct = strings.Split(fc[i], ":")
		for j := 0; j < len(fct); j += 1 {
			if len(fct[j]) > 0 {
				fci, err = strconv.ParseInt(fct[j], 10, 64)
				fcl = append(fcl, fci)
			}
		}
	}
	if len(fcl) == 0 {
		return errors.New("at least on father category must be provided for a link")
	}
	fcs, err := json.Marshal(fcl)

	// parse the tags
	tl := make([]int64, 0, len(tags))
	var tt []string
	var ti int64
	for i := 0; i < len(tags); i += 1 {
		tt = strings.Split(tags[i], ":")
		for j := 0; j < len(tt); j += 1 {
			if len(tt[j]) > 0 {
				ti, err = strconv.ParseInt(tt[j], 10, 64)
				tl = append(tl, ti)
			}
		}
	}
	ts, err := json.Marshal(tl)

	// fetch and fill the missing information
	ir := false
	if img == "" {
		ir = true
	}
	l, err := FetchLinkInformation(url, id, ir)
	if name == "" {
		name = l.Name
	}
	if d == "" {
		d = l.Description
	}
	if ir {
		img = l.IMG
	}

	// copy the img file
	if !ir {
		exist, err := FileExist("./static/img/front/")
		if err != nil {
			return err
		}
		if !exist {
			err = os.Mkdir("./static/img/front/", os.ModeDir)
			if err != nil {
				return err
			}
		}

		if img == "" {
			return errors.New("the img file is required")
		}
		tails := strings.Split(img, "/")
		if tails[len(tails)-1] == "" {
			return errors.New("the given img file is a dir")
		}
		tails = strings.Split(tails[len(tails)-1], ".")
		tail := tails[len(tails)-1]
		err = Copy(img, "./static/img/front/"+fmt.Sprint(id)+"."+tail)
		img = "./static/img/front/" + fmt.Sprint(id) + "." + tail
	}

	// insert the data to db
	_, err = db.Exec(
		`INSERT INTO CATEGORY (ID, Name, URL, Description,IMG ,FatherCategory, Priority, Tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		name,
		url,
		d,
		img,
		fcs,
		priority,
		ts,
	)
	if err != nil {
		return err
	}

	// update the father category
	var link []int64
	var ls []byte
	var idj int64
	for i := 0; i < len(fcl); i += 1 {

		idj = fcl[i]
		err = db.QueryRow("SELECT Links FROM Category WHERE ID = ?", idj).Scan(&ls)
		if errors.Is(err, sql.ErrNoRows) {

		} else if err != nil {
			return err
		} else {

			if len(ls) == 0 {
				link = make([]int64, 0, 1)
			} else {
				err = json.Unmarshal(ls, &link)
				if err != nil {
					return err
				}
			}

			link = append(link, int64(id))
			ls, err = json.Marshal(link)
			if err != nil {
				return err
			}

			_, err = db.Exec("UPDATE Category SET Links = ? WHERE ID = ?", ls, idj)
			if err != nil {
				return err
			}

		}

	}

	return nil
}

// TODO FetchLinkInformation
// fetch the target html and detect informations including names from the target html
// if imgRequired is true this will download the front img of the website and save it to the static/img/front/id.png
// id is used to store the img file, if the target file exist, this will replace it
func FetchLinkInformation(url string, id int, imgRequired bool) (l *Link, err error) {
	return nil, nil
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = in.Close()
	if err != nil {
		return err
	}

	return out.Close()
}

// test if file exist
func FileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

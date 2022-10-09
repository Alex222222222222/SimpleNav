package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var Templates map[string]*Template

var TemplateVariableRegexp *regexp.Regexp
var TemplateRequiredVariableRegexp *regexp.Regexp

type Template struct {
	Name             string
	RequiredVariable []string
	HTML             []byte
	Replace          [][]int
}

func newTemplate(path string, name string) (err error) {

	var t = &Template{}

	// check if the name already exist in all the templates
	if _, ok := Templates[name]; ok {
		return errors.New(
			fmt.Sprintf("Template with required name: %s already exist.", name),
		)
	}

	t.HTML, err = os.ReadFile(path)
	if err != nil {
		return err
	}

	// Load the basic html
	t.Replace = TemplateVariableRegexp.FindAllIndex(t.HTML, -1)

	// Load all required variables
	RequiredVariable := TemplateRequiredVariableRegexp.FindSubmatch(t.HTML)
	t.RequiredVariable = strings.Split(string(RequiredVariable[0]), ",")

	t.Name = name

	Templates[name] = t

	return nil

}

func InitTemplate() (err error) {

	Templates = make(map[string]*Template, 0)

	TemplateVariableRegexp = regexp.MustCompile("{{[A-Za-z0-9_]+}}")
	TemplateRequiredVariableRegexp = regexp.MustCompile("<!--\\nRequiredVariable:([A-Za-z0-9_,]+)\\n-->")

	// TODO config this in the config.toml
	templatesPath := map[string]string{
		"index":                   "./template/index.html",
		"category":                "./template/category.html",
		"categoryCard":            "./template/categoryCard.html",
		"categorySubCategory":     "./template/categorySubCategory.html",
		"categoryWithSubCategory": "./template/categoryWithSubCategory.html",
		"linkTags":                "./template/linkTags.html",
		"sidebarSingleItem":       "./template/sidebarSingleItem.html",
		"sidebarSubItem":          "./template/sidebarSubItem.html",
		"sidebarWithSubItem":      "./template/sidebarWithSubItem.html",
	}

	for name, path := range templatesPath {
		err := newTemplate(path, name)
		if err != nil {
			return err
		}
	}

	return nil

}

func renderSidebarSubitem(c *Category, hidden bool) (res []byte, err error) {
	if !hidden || !c.Hidden {
		data := make(map[string][]byte, 0)
		t := Templates["sidebarSubItem"]

		data["TERM_ID"] = []byte(
			"term" + fmt.Sprint(c.ID),
		)

		data["TERM_NAME"] = []byte(c.Name)

		return t.Render(data)
	} else {
		return make([]byte, 0), nil
	}
}

func renderSideBarWithSubitem(c *Category, hidden bool) (res []byte, err error) {

	if !hidden || !c.Hidden {
		data := make(map[string][]byte, 0)
		t := Templates["sidebarWithSubItem"]

		data["TERM_ID"] = []byte(
			"term" + fmt.Sprint(c.ID),
		)

		// TODO: change icon name to specific icon name
		data["ICON_NAME"] = []byte("icon-book-mark")

		data["TERM_NAME"] = []byte(c.Name)

		var subitems = make([]byte, 0)

		for i := 0; i < len(c.SubCategory); i += 1 {
			cs := AllCategories[c.SubCategory[i]]
			if len(cs.SubCategory) > 0 {
				return nil, errors.New("Subcategory with subcategory has not been supported yet")
			}

			html, err := renderSidebarSubitem(cs, hidden)
			if err != nil {
				return nil, err
			}

			subitems = append(subitems, html...)
		}

		data["SIDEBAR_SUBITEM"] = subitems

		return t.Render(data)
	} else {
		return nil, nil
	}

}

func renderSideBarSingleItem(c *Category, hidden bool) (res []byte, err error) {
	if !hidden || !c.Hidden {
		data := make(map[string][]byte, 0)
		t := Templates["sidebarSingleItem"]

		data["TERM_ID"] = []byte(
			"term" + fmt.Sprint(c.ID),
		)

		// TODO: change icon name to specific icon name
		data["ICON_NAME"] = []byte("icon-book-mark")

		data["TERM_NAME"] = []byte(c.Name)

		return t.Render(data)
	} else {
		return make([]byte, 0), nil
	}

}

func renderSideBar(hidden bool) (res []byte, err error) {
	for i := 0; i < len(Categories); i += 1 {
		c := AllCategories[Categories[i]]
		if len(c.SubCategory) > 0 && c.FatherCategory >= 0 {
			html, err := renderSideBarWithSubitem(c, hidden)
			if err != nil {
				return make([]byte, 0), err
			}
			res = append(res, html...)
		} else if c.FatherCategory >= 0 {
			html, err := renderSideBarSingleItem(c, hidden)
			if err != nil {
				return make([]byte, 0), err
			}
			res = append(res, html...)
		}
	}

	return res, nil
}

// TODO function to render main content
func renderIndexCategory(hidden bool) (res []byte, err error) {
	return nil, nil
}

func (t *Template) Render(data map[string][]byte) (html []byte, err error) {
	html = make([]byte, 0, len(t.HTML))
	html = append(html, t.HTML...)

	for cnt := len(t.Replace) - 1; cnt >= 0; cnt -= 1 {
		variableName := string(html[t.Replace[cnt][0]+2 : t.Replace[cnt][1]-2])
		dataT, ok := data[variableName]
		if !ok {
			return make([]byte, 0), errors.New(("Required variable does not provided."))
		}

		fmt.Println(string(html[t.Replace[cnt][0]:t.Replace[cnt][1]]))
		html = append(html[:t.Replace[cnt][0]], append(dataT, html[t.Replace[cnt][1]:]...)...)

	}

	return html, nil
}

func renderIndex(hidden bool, path string) (err error) {

	requiredData := make(map[string][]byte, 0)
	for key, data := range Config {
		requiredData[key] = []byte(data)
	}
	requiredData["MAIN_CONTENT"], err = renderIndexCategory(hidden)
	if err != nil {
		return err
	}
	requiredData["SIDEBAR_MENU"], err = renderSideBar(hidden)
	if err != nil {
		return err
	}

	html, err := Templates["index"].Render(requiredData)
	if err != nil {
		return err
	}

	// write html to disk
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(html)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil

}

func Render(page string, hidden bool, path string) (err error) {
	switch page {
	case "Index":
		return renderIndex(hidden, path)
	default:
		return nil
	}
}

func RenderAll() (err error) {

	err = InitWebDir()
	if err != nil {
		return err
	}

	err = Render("Index", true, "./web/index.html")
	if err != nil {
		return err
	}

	return nil
}

package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"example.com/purandixit07/go-course/pkg/config"
	"example.com/purandixit07/go-course/pkg/models"
)

func AddTemplateData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// get the template cache from app config
	// create a template cache
	var tc map[string]*template.Template
	if app1.UseCache {
		tc = app1.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddTemplateData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}
	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

var app1 *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app1 = a
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named .page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		layoutMatches, err := filepath.Glob("./templates/*.layout.html")

		if err != nil {
			return myCache, err
		}

		if len(layoutMatches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

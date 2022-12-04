package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/tranvanloc412/bookings/pkg/config"
	"github.com/tranvanloc412/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplate(a *config.AppConfig) {
	app = a
}

func AppDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	// create template cache
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get template from cache
	buf := new(bytes.Buffer)
	ts, ok := tc[tmpl]

	if !ok {
		log.Fatal("Can't get template cache \n")
	}

	td = AppDefaultData(td)
	_ = ts.Execute(buf, td)

	// parse teamplate
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Fatal("Can't parse template %s \n", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	tc := make(map[string]*template.Template)

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		log.Printf("Can't get pages %s \n", err)
		return tc, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles(fmt.Sprintf("./templates/%s", name))

		if err != nil {
			log.Printf("Can't get layouts %s \n", err)
			return tc, err
		}

		tmpl, err = tmpl.ParseGlob("./templates/*.layout.html")
		if err != nil {
			log.Printf("Can't parse layouts %s \n", err)
			return tc, err
		}

		tc[name] = tmpl
	}
	fmt.Printf("tmpl: %+v", tc)
	return tc, nil
}

// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	_, inMap := tc[t]
// 	if !inMap {
// 		log.Println("Create new template")
// 		err = createTemplate(t)
// 		if err != nil {
// 			fmt.Println("error parsing template:", err)
// 		}
// 	} else {
// 		log.Println("use cache template")
// 	}

// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println("error parsing template:", err)
// 	}
// 	bytes.Buffer
// }

// func createTemplate(t string) error {
// 	var tmpl *template.Template
// 	var err error

// 	tmpls := []string{fmt.Sprintf("./templates/%s", t), "./templates/base.layout.html"}

// 	tmpl, err = template.ParseFiles(tmpls...)
// 	if err != nil {
// 		fmt.Println("error parsing template:", err)
// 	}
// 	tc[t] = tmpl
// 	return nil
// }

package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/tranvanloc412/bookings/pkg/config"
	"github.com/tranvanloc412/bookings/pkg/models"
)

var functions = template.FuncMap{}

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
	ts, ok := tc[tmpl]

	if !ok {
		log.Fatal("Can't get template cache \n")
	}

	buf := new(bytes.Buffer)

	td = AppDefaultData(td)
	_ = ts.Execute(buf, td)

	// parse teamplate
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Fatal("Can't parse template %s \n", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	tc := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		log.Printf("Can't get pages %s \n", err)
		return tc, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			log.Printf("Can't get layouts %s \n", err)
			return tc, err
		}

		tmpl, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			log.Printf("Can't parse layouts %s \n", err)
			return tc, err
		}
		if len(tmpl) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return tc, err
			}
		}

		tc[name] = ts
	}
	return tc, nil
}

// myCache := map[string]*template.Template{}

// pages, err := filepath.Glob("./templates/*.page.tmpl")
// if err != nil {
// 	return myCache, err
// }

// for _, page := range pages {
// 	name := filepath.Base(page)
// 	ts, err := template.New(name).Funcs(functions).ParseFiles(page)
// 	if err != nil {
// 		return myCache, err
// 	}

// 	matches, err := filepath.Glob("./templates/*.layout.tmpl")
// 	if err != nil {
// 		return myCache, err
// 	}

// 	if len(matches) > 0 {
// 		ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
// 		if err != nil {
// 			return myCache, err
// 		}
// 	}

// 	myCache[name] = ts
// }

// return myCache, nil
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

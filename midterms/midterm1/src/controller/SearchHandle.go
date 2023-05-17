package controller

import (
	"html/template"
	"midterm1/models"
	"net/http"
)

func SearchPage(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	r.ParseForm()

	name := r.FormValue("productName")
	products := GetProductsByName(name)

	toInData := models.ToIndexData{
		Products: products,
	}
	tpl.ExecuteTemplate(w, "index.html", toInData)
}

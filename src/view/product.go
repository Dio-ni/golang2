package view

import (
	"html/template"
	"log"
	"midterm1/controller"
	"midterm1/models"
	"net/http"
)

func ProductsHandle(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	products := controller.GetProductsByName("")

	toInData := models.ToIndexData{
		Products: products,
	}
	tpl.ExecuteTemplate(w, "productList.html", toInData)

	// if nerr != nil {
	// 	log.Println(nerr)
	// }
}

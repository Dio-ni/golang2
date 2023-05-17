package controller

import (
	"fmt"
	"html/template"
	"log"
	"midterm1/config"
	"midterm1/models"
	"net/http"

	"github.com/gorilla/mux"
)

func RenderProductPage(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	params := mux.Vars(r)
	productId := params["id"]
	userId := r.FormValue("userId")

	p, err := GetProduct(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
		return
	}

	comments, err := GetComments(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}
	var rate string
	db, err := config.LoadDB()
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	stmt := "SELECT rate FROM ratings WHERE productId = ? and userId=?"
	row := db.QueryRow(stmt, productId, userId)
	err = row.Scan(&rate)

	fmt.Println(rate)
	data := models.ProductPage{
		Product:  p,
		Comments: comments,
		Rate:     rate,
	}

	tpl.ExecuteTemplate(w, "product.html", data)
}

func GetProductsByName(name string) []models.Product {
	db, err := config.LoadDB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM products WHERE car_name LIKE ? ORDER by rating DESC;;", "%"+name+"%")
	if err != nil {
		log.Println(err)
	}

	products := []models.Product{}
	for result.Next() {
		var p models.Product
		err = result.Scan(&p.Id, &p.Car_name, &p.Details, &p.Price, &p.Rate)
		if err != nil {
			log.Println(err)
		}
		products = append(products, p)
	}

	return products
}

func GetProduct(productId string) (models.Product, error) {
	db, err := config.LoadDB()
	if err != nil {
		return models.Product{}, err
	}
	defer db.Close()

	result := db.QueryRow("SELECT * FROM product WHERE product_id = ?", productId)
	var p models.Product
	err = result.Scan(&p.Id, &p.Car_name, &p.Details, &p.Price, &p.Rate)
	if err != nil {
		return models.Product{}, err
	}

	return p, nil
}

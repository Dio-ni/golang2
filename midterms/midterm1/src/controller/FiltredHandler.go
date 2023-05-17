package controller

import (
	"database/sql"
	"html/template"
	"log"
	"midterm1/config"
	"midterm1/models"
	"net/http"
	"strconv"
)

func FiltredProduct(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	r.ParseForm()

	minPriceStr := r.FormValue("minPrice")
	minPrice, err := strconv.Atoi(minPriceStr)
	if err != nil {
		log.Println(err)
	}

	maxPriceStr := r.FormValue("maxPrice")
	maxPrice, err := strconv.Atoi(maxPriceStr)
	if err != nil {
		log.Println(err)
	}
	database, err := config.LoadDB()
	products, err := GetFilteredProducts(database, minPrice, maxPrice)
	if err != nil {
		log.Println(err)
	}

	toInData := models.ToIndexData{
		Products: products,
	}
	tpl.ExecuteTemplate(w, "index.html", toInData)
}

func GetFilteredProducts(db *sql.DB, minPrice, maxPrice int) ([]models.Product, error) {
	// Prepare the SQL query
	query := "SELECT id, car_name,details, price,rating FROM products  WHERE price >= ? AND price <= ? ORDER by rating DESC;"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(minPrice, maxPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the results
	products := []models.Product{}

	// Iterate over the rows
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.Id, &p.Car_name, &p.Details, &p.Price, &p.Rate); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

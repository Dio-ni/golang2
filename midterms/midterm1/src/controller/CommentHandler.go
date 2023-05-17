package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"midterm1/config"
	"midterm1/models"
	"net/http"
)

func SendComment(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	fmt.Println("inside")
	var result string
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "product.html", nil)
		return
	}

	commentText := r.FormValue("commentText")
	fmt.Println(commentText)
	if len(commentText) == 0 {
		result = "write some comment"
		tpl.ExecuteTemplate(w, "product.html", result)
		return
	} else {

		userId := r.FormValue("userId")
		productId := r.FormValue("productId")
		var insertStmt *sql.Stmt
		database, err := config.LoadDB()
		insertStmt, err2 := database.Prepare("INSERT INTO comments (productid,userid, comment) VALUES (?, ?,?);")
		// fmt.Println(userId)
		if err2 != nil {
			fmt.Println("error preparing statement:", err2)
			tpl.ExecuteTemplate(w, "index.html", "there was a problem registering account")
			return
		}
		defer insertStmt.Close()
		var result sql.Result

		result, err2 = insertStmt.Exec(productId, userId, commentText)
		lastIns, _ := result.LastInsertId()
		fmt.Println("lastIns comment:", lastIns)
		if err2 != nil {
			fmt.Println("error inserting new user")
			tpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
			return
		}
		p, err := GetProduct(productId)
		fmt.Println(productId)
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

		data := models.ProductPage{
			Product:  p,
			Comments: comments,
		}

		http.Redirect(w, r, "/product:"+productId, http.StatusSeeOther)
		tpl.ExecuteTemplate(w, "product.html", data)
	}
}

func GetComments(productId string) ([]models.Comment, error) {
	db, err := config.LoadDB()
	if err != nil {
		return []models.Comment{}, err
	}
	defer db.Close()

	res, err := db.Query("SELECT u.name, c.comment FROM comments c join users u on u.userid=c.userid WHERE productId = ?", productId)
	if err != nil {
		return []models.Comment{}, err
	}

	comments := []models.Comment{}
	for res.Next() {
		var c models.Comment
		err = res.Scan(&c.Name, &c.Comment)
		if err != nil {
			return []models.Comment{}, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

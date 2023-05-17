package view

import (
	"database/sql"
	"fmt"
	"html/template"
	"midterm1/config"
	"midterm1/controller"
	"midterm1/models"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	var userid int64
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "registration.html", "")
		return
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	// fmt.Println("password:", password, "\npswdLength:", len(password))

	stmt := "SELECT userid FROM users WHERE name = ?"
	db, err := config.LoadDB()
	row := db.QueryRow(stmt, username)

	err = row.Scan(&userid)

	if err != sql.ErrNoRows {
		fmt.Println("username already exists, err:", err)
		tpl.ExecuteTemplate(w, "registration.html", "username already taken")
		return
	}
	// func (db *DB) Prepare(query string) (*Stmt, error)
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO users (name, password) VALUES (?, ?);")

	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()
	var result sql.Result

	result, err = insertStmt.Exec(username, password)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
		return
	}
	userid = lastIns
	products := controller.GetProductsByName("")

	data := models.ToIndexData{
		Username: username,
		UserId:   userid,
		Products: products}

	tpl.ExecuteTemplate(w, "index.html", data)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

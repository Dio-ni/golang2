package main

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"midterm1/config"
	"midterm1/models"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
)

var database *sql.DB
var toInData models.ToIndexData
var savUsername string
var tmpl *template.Template
var store = sessions.NewCookieStore([]byte("4hq2HxUvD0r6hJ5y"))

type Comment struct {
	Name    string
	Comment string
}
type ProductPage struct {
	Product  Product
	Comments []Comment
}

type sendData struct {
	Username string
	Products []Product
	Comments []Comment
}

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	CategoryID  int64   `json:"category_id"`
	Rating      float64 "json:rating"
}

type User struct {
	Username    string "json:username"
	Email       string "json:email"
	Password    string "json:password"
	FirstName   string "json:first_name"
	LastName    string "json:last_name"
	Address     string "json:address"
	PhoneNumber string "json:phone_number"
}

// MAIN PAGE || WELCOMING PAGE
func MainPage(w http.ResponseWriter, r *http.Request) {

	type Username struct {
		Username string
	}
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	namefromsession, ok := session.Values["username"].(string)
	if !ok {

		tpl.ExecuteTemplate(w, "index.html", Username{Username: ""})
	} else {
		tpl.ExecuteTemplate(w, "index.html", Username{Username: namefromsession})
	}

}

// SIGN IN AND SIGN UP
// signup serves form for registring new users
func signup(w http.ResponseWriter, r *http.Request) {
	var user_id int64
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "signup.html", "")
		return
	}

	//COLLECTING DATA FROM FRONT RESPONSE
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// CHECKING USERNAME IS FREE
	stmt := "SELECT user_id FROM user WHERE username = ?"
	db, err := config.LoadDB()
	row := db.QueryRow(stmt, user.Username)

	err = row.Scan(&user_id)

	if err != sql.ErrNoRows {
		fmt.Println("username already exists, err:", err)
		http.Error(w, "username already taken", http.StatusInternalServerError)
		return
	}

	// INSERRING NEW USER
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO user(username, email, password, first_name, last_name, address, phone_number) VALUES(?,?,?,?,?,?,?);")

	if err != nil {
		fmt.Println("error preparing statement:", err)
		http.Error(w, "There was a problem registering the account", http.StatusInternalServerError)
		return
	}
	defer insertStmt.Close()
	var result sql.Result
	result, err = insertStmt.Exec(user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.Address, user.PhoneNumber)
	if err != nil {
		fmt.Println("error inserting new user", err)
		http.Error(w, "There was a problem registering the account", http.StatusInternalServerError)
		return
	}

	//GET LAST INSERTED USERS ID
	lastIns, _ := result.LastInsertId()
	fmt.Println("lastIns:", lastIns)
	//GET SESSION CLEAN IT AND SAVE USERS DATA (ID, NAME)
	session, err := store.Get(r, "session")
	if err != nil {
		fmt.Print(err)
	}
	delete(session.Values, "userId")
	delete(session.Values, "username")
	session.Save(r, w)
	session.Values["userId"] = lastIns
	session.Values["username"] = user.Username
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("session user saving error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//REDIRECT TO MAIN PAGE
	fmt.Println("User created successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func signin(w http.ResponseWriter, r *http.Request) {

	db, err := config.LoadDB()
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "login.html", "")
		return
	}
	//GET DATA FROM FORM
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	//SELECT USERS DATA
	var pass string
	var userid int64
	stmt := "SELECT user_id, password FROM user WHERE username = ?"
	row := db.QueryRow(stmt, username)
	err = row.Scan(&userid, &pass)
	if err != nil {
		fmt.Println("error selecting  Username")
		tpl.ExecuteTemplate(w, "login.html", "incorrect password or username")
		return
	}
	//IF PASSWORD CORRECT SAVE DATA(ID,NAME) TO SESSION AND REDIRET TO MAIN PAGE
	if password == pass {
		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Print(err)
		}
		delete(session.Values, "userId")
		delete(session.Values, "username")
		session.Save(r, w)

		session.Values["userId"] = userid
		session.Values["username"] = username
		err = session.Save(r, w)
		if err != nil {
			fmt.Println("session user saving error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		// IF NOT SEND ERROR MESSAGE TO FRONT
		fmt.Println("error selecting")
		tpl.ExecuteTemplate(w, "login.html", "incorrect password or username")
		return
	}

	fmt.Println("User signED in ")

}

// PRODUCTS LIST PAGE WITH SEARCHING, FILTERING FUNC-S
func ProductList(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	products, _ := GetProductsByName("")
	sendData := sendData{
		Products: products,
	}

	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	namefromsession, ok := session.Values["username"].(string)
	if !ok {
		sendData.Username = ""
	} else {
		sendData.Username = namefromsession
	}

	tpl.ExecuteTemplate(w, "productList.html", sendData)

}
func GetProductsByName(product_name string) ([]Product, error) {
	db, err := config.LoadDB()
	if err != nil {
		return nil, err
	}

	query := "SELECT product_id, name, description, price, image_url, category_id,rating FROM product where name LIKE ?"
	rows, err := db.Query(query, "%"+product_name+"%")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.CategoryID, &p.Rating)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	r.ParseForm()

	name := r.FormValue("productName")
	products, _ := GetProductsByName(name)

	sendData := sendData{
		Products: products,
	}
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	namefromsession, ok := session.Values["username"].(string)
	if !ok {
		sendData.Username = ""
	} else {
		sendData.Username = namefromsession
	}

	tpl.ExecuteTemplate(w, "productList.html", sendData)

}

func filtredProduct(w http.ResponseWriter, r *http.Request) {
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

	db, err := config.LoadDB()
	// Prepare the SQL query
	query := "SELECT * FROM product  WHERE price >= ? AND price <= ? ORDER by rating DESC;"
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(minPrice, maxPrice)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	// Create a slice to hold the results
	products := []Product{}

	// Iterate over the rows
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.CategoryID, &p.Rating); err != nil {
			fmt.Println(err)
		}
		products = append(products, p)
	}
	var sendData sendData
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	namefromsession, ok := session.Values["username"].(string)
	if !ok {
		sendData.Username = ""
	} else {
		sendData.Username = namefromsession
	}
	sendData.Products = products
	tpl.ExecuteTemplate(w, "productList.html", sendData)

}

func renderProductPage(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	productId := params["id"]

	p, err := GetProductbyId(productId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
		return
	}

	comments, err := GetComments(productId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}
	db, err := config.LoadDB()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	data := ProductPage{
		Product:  p,
		Comments: comments,
	}
	fmt.Println(data)
	tpl.ExecuteTemplate(w, "product.html", data)
}
func GetComments(productId string) ([]Comment, error) {
	db, err := config.LoadDB()
	if err != nil {
		return []Comment{}, err
	}
	defer db.Close()

	res, err := db.Query("SELECT u.username, c.comment_text FROM comments c join user u on u.user_id=c.user_id WHERE product_id = ?", productId)
	if err != nil {
		return []Comment{}, err
	}

	comments := []Comment{}
	for res.Next() {
		var c Comment
		err = res.Scan(&c.Name, &c.Comment)
		if err != nil {
			return []Comment{}, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func GetProductbyId(productId string) (Product, error) {
	db, err := config.LoadDB()
	if err != nil {
		return Product{}, err
	}
	defer db.Close()

	result := db.QueryRow("SELECT * FROM product WHERE product_id = ?", productId)
	var p Product
	err = result.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.CategoryID, &p.Rating)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func sendComment(w http.ResponseWriter, r *http.Request) {
	var result string
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "product.html", nil)
		return
	}

	commentText := r.FormValue("commentText")
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	namefromsession, _ := session.Values["userId"]
	userId := namefromsession
	productId := r.FormValue("productId")

	fmt.Println(commentText)
	if len(commentText) == 0 {
		result = "write some comment"
		tpl.ExecuteTemplate(w, "product.html", result)
		return
	} else {

		var insertStmt *sql.Stmt
		database, _ := config.LoadDB()

		insertStmt, err2 := database.Prepare("INSERT INTO comments (product_id,user_id, comment_text) VALUES (?, ?,?);")
		// INSERT INTO `comments`(`id`, `product_id`, `user_id`, `comment_text`, `created_at`) VALUES
		// fmt.Println(userId)
		if err2 != nil {
			fmt.Println("error preparing statement:", err2)
			tpl.ExecuteTemplate(w, "index.html", "there was a problem registering account")
			return
		}
		defer insertStmt.Close()
		var result sql.Result
		fmt.Println(productId, userId, commentText)
		result, err2 = insertStmt.Exec(productId, userId, commentText)
		if err2 != nil {
			fmt.Println(err2)
		}
		lastIns, _ := result.LastInsertId()
		fmt.Println("lastIns comment:", lastIns)
		if err2 != nil {
			fmt.Println("error inserting new user")
			tpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
			return
		}

	}
	http.Redirect(w, r, "/product:"+productId, http.StatusSeeOther)

}

func sendRating(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // Parses the request body
	rating := r.Form.Get("rating")
	productId := r.Form.Get("productId")
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := session.Values["userId"]

	db, err := config.LoadDB()

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO ratings (rate, product_Id,user_Id) VALUES (?, ?,?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	defer insertStmt.Close()
	var result sql.Result
	result, err = insertStmt.Exec(rating, productId, userId)
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	lastIns, _ := result.LastInsertId()
	fmt.Print(lastIns)
	http.Redirect(w, r, "/product:"+productId, http.StatusSeeOther)
}
func profile(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "profile.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// uploadHandler handles the file upload request
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(10 << 20) // Max file size of 10MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file data
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file data", http.StatusInternalServerError)
		return
	}
	// folderPath := "uploads"
	// err = os.MkdirAll(folderPath, 0755)
	// if err != nil {
	// 	http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
	// 	return
	// }

	folderPath := "static/img"
	filePath := filepath.Join(folderPath, handler.Filename)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, "Failed to get absolute file path", http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(absPath, fileBytes, 0644)
	if err != nil {
		http.Error(w, "Failed to save file to disk", http.StatusInternalServerError)
		return
	}

	// Save the file path to the database
	err = saveFilePathToDB(filePath)
	if err != nil {
		http.Error(w, "Failed to save file path to database", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully!")
}

// saveFilePathToDB saves the file path to the database
func saveFilePathToDB(filePath string) error {
	// Prepare the SQL statement
	db, err := config.LoadDB()

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO files (file_path) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(filePath)
	if err != nil {
		return err
	}

	return nil
}

var tpl *template.Template

func main() {

	var err error
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	db, err := config.LoadDB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	database = db

	routes := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{path: "/", handler: MainPage},
		{path: "/productList", handler: ProductList},
		{path: "/search", handler: SearchPage},
		{path: "/signup", handler: signup},
		{path: "/signin", handler: signin},
		{path: "/profile", handler: profile},
		{path: "/filtred", handler: filtredProduct},
		{path: "/product:{id}", handler: renderProductPage},
		{path: "/sendComment", handler: sendComment},
		{path: "/ratings", handler: sendRating},
		{path: "/add-product", handler: uploadHandler},
	}
	r := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	for _, route := range routes {
		r.HandleFunc(route.path, route.handler)
	}

	log.Fatal(http.ListenAndServe(":8080", r))

	fmt.Println("Server is listening...")
	// http.ListenAndServe(":8080", nil)
}

func init() {
	// Register User type with gob
	gob.Register(User{})
}

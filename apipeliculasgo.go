package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type Genres struct {
	Gen_id    int    `json:"gen_id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Gen_title string `json:"gen_title"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}
func createNewBooking(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(reqBody)
	var genres Genres
	json.Unmarshal(reqBody, &genres)
	db.Create(&genres)
	fmt.Println("Creating a new Genre")
	json.NewEncoder(w).Encode(genres)
}
func returnAllBookings(w http.ResponseWriter, r *http.Request) {
	genres := []Genres{}
	db.Find(&genres)
	fmt.Println("Return all Genres")
	json.NewEncoder(w).Encode(genres)
}
func returnSinbleBookingdos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	genres := []Genres{}
	s, err := strconv.Atoi(key)
	if err == nil {
		db.Where("gen_id = ?", s).Find(&genres)
		json.NewEncoder(w).Encode(genres)
	}

}
func returnSingleBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	genres := []Genres{}
	db.Find(&genres)
	for _, genres := range genres {
		// string to int
		s, err := strconv.Atoi(key)
		if err == nil {
			if genres.Gen_id == s {
				fmt.Println(genres)
				fmt.Println("Genre", key)
				json.NewEncoder(w).Encode(genres)
			}
		}
	}
}
func deleteBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	s, err := strconv.Atoi(key)
	if err == nil {
		db.Delete(&Genres{Gen_id: s})
		fmt.Println("Delete")
	}

}
func updateBooking(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data Genres
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	vars := mux.Vars(r)
	key := vars["id"]
	s, er := strconv.Atoi(key)

	if er == nil {
		data.Gen_id = s
		db.Save(&data)
		fmt.Println("Actualizado", data)
	}

}
func handleRequests() {
	log.Println("Starting development server at http://127.0.0.1:3000/")
	log.Println("Quit the server with CONTROL-C.")
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"x-Requested-with", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/genres", createNewBooking).Methods("POST")
	myRouter.HandleFunc("/genres", returnAllBookings).Methods("GET")
	myRouter.HandleFunc("/genres/{id}", returnSinbleBookingdos).Methods("GET")
	myRouter.HandleFunc("/genres/{id}", deleteBooking).Methods("DELETE")
	myRouter.HandleFunc("/genres/{id}", updateBooking).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(myRouter)))

}

func main() {
	// Please define your user name and password for my sql.
	db, err = gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/apigopeliculas?charset=utf8&parseTime=True")
	// NOTE: See we’re using = to assign the global var
	// instead of := which would assign it only in this function

	if err != nil {
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}

	db.AutoMigrate(&Genres{})
	handleRequests()

}

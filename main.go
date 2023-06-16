package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name   string `json:"name"`
	Rollno string `json:"rollno"`
	City   string `json:"city"`
}

var DB *gorm.DB
var err error

const DSN = "saiteja:S@iteja555@tcp(127.0.0.1:3306)/saiDB?charset=utf8mb4&parseTime=True&loc=Local"

func InitialMigration() {
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	DB.AutoMigrate(&Student{})
}

func allStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student
	DB.Find(&students)
	fmt.Println("{}", students)
	json.NewEncoder(w).Encode(students)
}

func singleStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var student Student
	DB.First(&student, params["id"])
	json.NewEncoder(w).Encode(student)
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	json.NewDecoder(r.Body).Decode(&student)
	DB.Create(&student)
	json.NewEncoder(w).Encode(student)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var student Student
	DB.First(&student, params["id"])
	json.NewDecoder(r.Body).Decode(&student)
	DB.Save(&student)
	json.NewEncoder(w).Encode(student)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var student Student
	DB.Delete(&student, params["id"])
	json.NewEncoder(w).Encode(student)
}

func createStudent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]
	rollno := vars["rollno"]
	city := vars["city"]

	fmt.Println(name)
	fmt.Println(rollno)
	fmt.Println(city)

	DB.Create(&Student{Name: name, Rollno: rollno, City: city})
	fmt.Fprintf(w, "New Student Successfully Created")
}

func handleRequests() {
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/students", allStudents).Methods("GET")
	myRouter.HandleFunc("/students/{id}", singleStudent).Methods("GET")
	myRouter.HandleFunc("/students", CreateStudent).Methods("POST")
	myRouter.HandleFunc("/students/{id}", UpdateStudent).Methods("PUT")
	myRouter.HandleFunc("/students/{id}", DeleteStudent).Methods("DELETE")
	myRouter.HandleFunc("/students/{name}/{rollno}/{city}", createStudent).Methods("POST")

	log.Fatal(http.ListenAndServe(":3003", myRouter))
}

func main() {
	fmt.Println("Go ORM")
	InitialMigration()
	handleRequests()
}

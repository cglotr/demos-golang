package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	var employee employee
	log.Println(r.Body)
	_ = json.NewDecoder(r.Body).Decode(&employee)
	st, err := db.Prepare("INSERT INTO employee (name, city) VALUES (?,?)")

	if err != nil {
		panic(err.Error())
	}

	st.Exec(employee.Name, employee.City)
	log.Printf("INSERT employee name: %s, city: %s\n", employee.Name, employee.City)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	params := mux.Vars(r)
	st, err := db.Prepare("DELETE FROM employee WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}

	employeeID, err := strconv.Atoi(params["id"])

	if err != nil {
		panic(err.Error())
	}

	st.Exec(employeeID)
	log.Printf("DELETE employee id: %d\n", employeeID)
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbName := "goblog"
	dbPass := "T3$t!992"
	dbUser := "test"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM employee ORDER BY id ASC")

	if err != nil {
		panic(err.Error())
	}

	employees := []employee{}

	for rows.Next() {
		var id int
		var name, city string

		err := rows.Scan(&id, &name, &city)

		if err != nil {
			panic(err.Error())
		}

		employee := employee{id, name, city}
		employees = append(employees, employee)
	}

	json.NewEncoder(w).Encode(employees)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/employees", getEmployees).Methods("GET")
	router.HandleFunc("/employees", createEmployee).Methods("POST")
	router.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}

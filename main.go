// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Schedulerjob struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

var db *sql.DB
var err error
var Jobs []Schedulerjob

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnSingleScheduler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	for _, job := range Jobs {
		if job.Id == key {
			json.NewEncoder(w).Encode(job)
		}
	}
}

func createNewScheduler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job Schedulerjob
	json.Unmarshal(reqBody, &job)
	Jobs = append(Jobs, job)
	json.NewEncoder(w).Encode(job)
}

func deleteScheduler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for index, job := range Jobs {
		if job.Id == id {
			Jobs = append(Jobs[:index], Jobs[index+1:]...)
		}
	}
}

func returnAllSchedulerJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: schedulers")
	w.Header().Set("Content-Type", "application/json")
	var jobs []Schedulerjob
	result, err := db.Query("SELECT Id,Name from schedulerstate order by Id")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var job Schedulerjob
		err := result.Scan(&job.Id, &job.Name)
		if err != nil {
			panic(err.Error())
		}
		jobs = append(jobs, job)
	}
	json.NewEncoder(w).Encode(jobs)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/schedulers", returnAllSchedulerJobs)
	myRouter.HandleFunc("/scheduler", createNewScheduler).Methods("POST")
	myRouter.HandleFunc("/scheduler/{id}", deleteScheduler).Methods("DELETE")
	myRouter.HandleFunc("/scheduler/{id}", returnSingleScheduler)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/scheduler")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	handleRequests()
}

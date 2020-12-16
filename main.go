package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Order represents the model for an order
// Default table name will be `orders`
type Machine struct {
	// gorm.Model
	MachineID      uint   `json:"orderId" gorm:"primary_key;auto_increment;not_null"`
	MachineName string    `json:"machineName"`
	CreateAt    time.Time `json:"orderedAt"`
	//Items        []Item    `json:"items" gorm:"foreignkey:OrderID"`
}

// Item represents the model for an item in the order
/**
type Item struct {
	// gorm.Model
	LineItemID  uint   `json:"lineItemId" gorm:"primary_key;auto_increment;not_null"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"-"`
}
*/

var db *gorm.DB
var token *string

func initDB(){
	var err error
	dataSourceName := "root:@tcp(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// Create the database. This is a one-time step.
	// Comment out if running multiple times - You may see an error otherwise
	db.Exec("CREATE DATABASE machine_db")
	db.Exec("USE machine_db")

	// Migration to create tables for Machine schema
	db.AutoMigrate(&Machine{})
}

func createMachine(w http.ResponseWriter, r *http.Request){

	//verified token
	mytoken := r.Header.Get("token")

	if(mytoken != *token){
		fmt.Fprintln(w, "you can not pass")
		return
	}
		
	var machine Machine
	json.NewDecoder(r.Body).Decode(&machine)

	//Create machine object into DB
	tNow := time.Now()
	tUnix := tNow.Unix()
	timeT := time.Unix(tUnix, 0)
	machine.CreateAt = timeT
	db.Create(&machine)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(machine)
}

func main(){
	var t string
	fmt.Println("Type token：")
	fmt.Scan(&t)
	token = &t

	var port string
	fmt.Println("Type API port：")
	fmt.Scan(&port)
	
	router := mux.NewRouter();

	//Create
	router.HandleFunc("/machine", createMachine).Methods("POST")

	// Initialize db connection
	initDB()

	log.Fatal(http.ListenAndServe(":"+port, router))
}
package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alikurb12/stocks_api_go/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file. %v", err)
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		log.Fatal("Error connecting to the database. %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database. %v", err)
	}
	fmt.Println("Connected to the database successfully!")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)
	fmt.Println("Creating stock:", stock)
	if err != nil {
		log.Fatal("Error decoding request body. %v", err)
	}

	insertID := InsertStock(stock)

	res := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}
	fmt.Println("Response:", res)
	json.NewEncoder(w).Encode(res)
	fmt.Println("Stock created successfully:", stock)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	fmt.Println("Fetching stock with ID:", id)
	if err != nil {
		log.Fatal("Error converting ID to integer. %v", err)
	}

	stock, err := getStock(int64(id))
	fmt.Println("Fetched stock:", stock)
	if err != nil {
		log.Fatal("Error fetching stock from database. %v", err)
	}
	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()
	for _, stock := range stocks {
		fmt.Println("Stock:", stock)
	}
	if err != nil {
		log.Fatal("Error fetching all stocks from database. %v", err)
	}
	json.NewEncoder(w).Encode(stocks)
	fmt.Println("All stocks fetched successfully")
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	fmt.Println("Updating stock with ID:", id)
	if err != nil {
		log.Fatal("Error converting ID to integer. %v", err)
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatal("Error decoding request body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Stock with ID %d updated successfully, %d rows affected", id, updatedRows)
	fmt.Println(msg)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
	fmt.Println("Stock updated successfully:", stock)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atio(params["id"])
	fmt.Println("Deleting stock with ID:", id)
	if err != nil {
		log.Fatal("Error converting ID to integer. %v", err)
	}

	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock with ID %d deleted successfully, %d rows affected", id, deletedRows)
	fmt.Println(msg)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
	fmt.Println("Stock deleted successfully with ID:", id)
}

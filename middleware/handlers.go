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

	insertID := insertStock(stock)

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
	id, err := strconv.Atoi(params["id"])
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

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()
	sqlStatement := "INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING id"
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatal("Error inserting stock into database. %v", err)
	}
	fmt.Println("Inserted stock with ID:", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()
	var stock models.Stock
	sqlStatement := "SELECT stockid, name, price, company FROM stocks WHERE stockid=$1"
	err := db.QueryRow(sqlStatement, id).Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Company)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Stock{}, fmt.Errorf("no stock found with ID %d", id)
		}
		log.Fatal("Error fetching stock from database. %v", err)
	}
	fmt.Println("Fetched stock with ID:", id)
	return stock, nil
}

func getAllStocks() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()
	sqlStatement := "SELECT * FROM stocks"
	var stocks []models.Stock
	fmt.Println("Fetching all stocks from database")
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatal("Error fetching all stocks from database. %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatal("Error scanning row. %v", err)
		}
		stocks = append(stocks, stock)
		fmt.Println("Fetched stock:", stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()
	sqlStatement := "UPDATE stocks SET name=$1, price=$2, company=$3 WHERE stockid=$4"
	res, err := db.Exec(sqlStatement, stock.Name, stock.Price, stock.Company, id)
	if err != nil {
		log.Fatal("Error updating stock in database. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error getting rows affected. %v", err)
	}
	fmt.Println("Updated stock with ID:", id, "Rows affected:", rowsAffected)
	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := CreateConnection()
	defer db.Close()
	sqlStatement := "DELETE FROM stocks WHERE stockid=$1"
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal("Error deleting stock from database. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error getting rows affected. %v", err)
	}
	fmt.Println("Deleted stock with ID:", id, "Rows affected:", rowsAffected)
	return rowsAffected
}

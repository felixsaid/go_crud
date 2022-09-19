package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"v1/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Response struct {
	ID      int64  `json: "id, omitempty"`
	Message string `json: "message, omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres...")

	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertId := insertStock(stock)

	res := Response{
		ID:      insertId,
		Message: "Stock successfully created.",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to int . %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get all the stock. %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode request body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected %v", updatedRows)

	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}

	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("Stock deleted successfully. Totla rows/records affected %v", deletedRows)

	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlQuerry := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlQuerry, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to process the query %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stock models.Stock

	sqlQuerry := `SELECT * FROM stocks WHERE stockid = $1`

	row := db.QueryRow(sqlQuerry, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned.")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan rows %v", err)
	}

	return stock, err
}

func getAllStocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stocks []models.Stock

	sqlQuery := `SELECT * FROM stocks`
	rows, err := db.Query(sqlQuery)

	if err != nil {
		log.Fatalf("Unable to execute querry %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock

		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Unable to scan row %v", err)
		}
		stocks = append(stocks, stock)
	}

	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlQuerry := `UPDATE stocks set name = $2, price = $3, company = $4 WHERE stockid = $1`
	res, err := db.Exec(sqlQuerry, id, &stock.Name, &stock.Price, &stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute querry %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error occured when checking affected rows %v", err)
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)

	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlQuery := `DELETE FROM stocks WHERE stockid = $1`

	res, err := db.Exec(sqlQuery, id)

	if err != nil {
		log.Fatalf("Could not execute querry %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error occured when checking the effected rows %v", err)
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)

	return rowsAffected
}

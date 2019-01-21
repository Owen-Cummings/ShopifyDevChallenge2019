package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//struct for managing database connection handles
type Db struct {
	*sql.DB
}

//formatting for products and cartEntries to be stored and retrieved
type product struct{
	ID    			int
	Title 			string
	Price 			float64
	InventoryCount	int
}

type cartEntry struct{
	CartID			int
	ProductID    	int
	Title 			string
	Price 			float64
}

//function for creating a database handle using specified connection information
// return connection or error encountered
func CreateConnection(connInfo string) (*Db, error){

	db, err := sql.Open("mysql", connInfo)
	if err != nil {
		return nil, err
	}

	//check connection to db, exit if error occurs
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Db{db}, err
}

//function to format arguments for db connection returning them as a string
func ConnInfo(username string, password string, dbName string) string{
	return fmt.Sprintf(
		"%s:%s@/%s", username, password, dbName,
		)
}

//Query to return
func (db *Db) GetProducts(input string) []product {
	stmt, err := db.Prepare(input)
	errorCheck(err, "GetProducts Prepare")

	results, err := stmt.Query()

	var p product
	products := []product{}

	for results.Next() {
		err = results.Scan(
			&p.ID,
			&p.Title,
			&p.Price,
			&p.InventoryCount,
		)
		errorCheck(err, "GetProducts Scan")
	products = append(products, p)
	}
	return products
}

//function to modify products
func (db *Db) ModifyProducts(query string) *sql.Rows{
	stmt, err := db.Prepare(query)
	errorCheck(err, "ModifyProducts Prepare")

	results, err := stmt.Query()

	errorCheck(err, "ModifyProducts Query")
	return results
}

//function to return cart details by parsing entries and returning them
func (db *Db) GetCart(input string) ([]cartEntry, float64) {
	stmt, err := db.Prepare(input)

	errorCheck(err, "GetCart Prepare")

	results, err := stmt.Query()

	var p cartEntry
	var totalCost float64
	products := []cartEntry{}

	for results.Next() {
		err = results.Scan(
			&p.CartID,
			&p.ProductID,
			&p.Title,
			&p.Price,
		)
		errorCheck(err,"GetCart Scan")

		products = append(products, p)
		totalCost += p.Price
	}
	return products, totalCost
}

// function to modify cart details executing a statement and returning nothing
func (db *Db) ModifyCart(query string) {
	stmt, err := db.Prepare(query)
	errorCheck(err, "ModifyCart Prepare")

	_, err = stmt.Exec()
	errorCheck(err, "ModifyCart Execution")
}

//Checkout operations for the databse
func (db *Db) CheckOut(query string) (int, float64){
	//variables to store calculated totals to be returned
	var totalItems int
	var totalPrice float64

	//execute passed argument to fetch cartEntries for matching cartID
	stmt, err := db.Prepare(query)
	errorCheck(err, "Checkout Prepare")

	//execute and assign output to results
	results, err := stmt.Query()

	var c cartEntry
	cartEntries := []cartEntry{}

	//iterate through cart contents stored in results
	for results.Next() {
		err = results.Scan(
			&c.CartID,
			&c.ProductID,
			&c.Title,
			&c.Price,
		)
		request := fmt.Sprintf("UPDATE PRODUCTS SET INVENTORY_COUNT = (INVENTORY_COUNT - 1)  WHERE PRODUCT_ID = %d", c.ProductID)
		_, err := db.Exec(request)
		errorCheck(err, "GetProducts Scan/removal")
		totalItems++
		totalPrice += c.Price
		cartEntries = append(cartEntries, c)
	}

	return totalItems, totalPrice
}

//
func errorCheck(err error, details string) {
	if err != nil{
		fmt.Printf("Error in %s: ", details)
		fmt.Println(err)
	}
}
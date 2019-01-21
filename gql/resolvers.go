package gql

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/Owen-Cummings/ShopifyDevChallenge2019/mysql"
)

// Resolver struct for accessing go-gql-api/mysql database
type Resolver struct {
	db *mysql.Db
}

// variables to temporarily store product details
var count int
var title string
var price float64

//simple resolver for database query, takes resolver, and graphQL parameters
func (resolver *Resolver) ProductResolver(params graphql.ResolveParams) (interface{}, error) {
	minStock := 0

	if params.Args["InStock"] == true{
		minStock = 1
	}

	//switch case for handling Product query elements
	if params.Args["Title"] != nil{
		//check if input matches column data requirements
		Title, ok := params.Args["Title"].(string)
		if ok {
			request := fmt.Sprintf("SELECT * FROM PRODUCTS WHERE TITLE = \"%s\" AND INVENTORY_COUNT >= %d", Title, minStock)
			products := resolver.db.GetProducts(request)
			return products, nil
		}
	} else {
		request := fmt.Sprintf("SELECT * FROM PRODUCTS WHERE INVENTORY_COUNT >= %d", minStock)
		products := resolver.db.GetProducts(request)
		return products, nil
	}

	//if input is invalid it will return nothing indicating failure to calling function
	return nil, nil
}

//function for validating stock and modifying products
func (resolver *Resolver) ProductMutator(params graphql.ResolveParams) (interface{}, error){

	Id, ok := params.Args["ID"].(int)
	if ok {
		//validate that there is more than one item available, if not return
		query := fmt.Sprintf("SELECT INVENTORY_COUNT, TITLE FROM PRODUCTS WHERE PRODUCT_ID = %d", Id)
		results := resolver.db.ModifyProducts(query)
		results.Next()
		err := results.Scan(&count, &title)
		if err != nil {
			return nil, err
		}
		//Check if item has enough stock to be purchased
		if count > 0 {
			//validate passed parameters
				query := fmt.Sprintf("UPDATE PRODUCTS SET INVENTORY_COUNT = (INVENTORY_COUNT - 1)  WHERE PRODUCT_ID = %d", Id)
				resolver.db.ModifyProducts(query)
				return fmt.Sprintf("Product: '%s' purchased successfully ", title), nil
			} else {
				return fmt.Sprintf("Product: '%s' not in stock", title), nil
			}
	}
	return nil, nil
}

// Resolver for getting cart entries for a specified cart and returning them
func (resolver *Resolver) CartResolver(params graphql.ResolveParams) (interface{}, error) {

	Id, ok := params.Args["CartID"].(int)
	if ok {
		request := fmt.Sprintf("SELECT * FROM CART WHERE CART_ID = %d", Id)
		cartEntries, totalCost := resolver.db.GetCart(request)
		fmt.Printf("Total Cost of Cart%d: $%.2f", Id, totalCost)
		return cartEntries, nil
	}

	return nil, nil
}

// resolver for altering cart content
func (resolver *Resolver) CartMutator(params graphql.ResolveParams) (interface{}, error){
	cartID, ok := params.Args["CartID"].(int)
	if ok {
		//create cart table with specified id to store products
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS CART (CART_ID INT NOT NULL, PRODUCT_ID INT NOT NULL, TITLE VARCHAR(50) NOT NULL,PRICE DECIMAL(10,2) NOT NULL);")
		resolver.db.ModifyCart(query)
		itemId, ok := params.Args["ItemID"].(int)
		if ok {
			query = fmt.Sprintf("SELECT * FROM PRODUCTS WHERE PRODUCT_ID = %d", itemId)
			results := resolver.db.ModifyProducts(query)
			results.Next()
			err := results.Scan(&itemId, &title, &price, &count)
			if err != nil {
				return nil, err
			}
			fmt.Printf("(%d, %d, %s, %.2f)", cartID, itemId, title, price)
			//Check if item has enough stock to be purchased
			if count > 0 {
				//validate passed parameters
				query := fmt.Sprintf("INSERT INTO CART VALUES (%d, %d, \"%s\", %.2f)", cartID, itemId, title, price)
				resolver.db.ModifyCart(query)
				return fmt.Sprintf("Product: '%s' added successfully ", title), nil
			} else {
				return fmt.Sprintf("Product: '%s' not in stock", title), nil
			}
		}
	}
	return nil, nil
}
// Resolver for checking out
func (resolver *Resolver) CheckOut(params graphql.ResolveParams) (interface{}, error) {
	Id, ok := params.Args["CartID"].(int)
	if ok {
		request := fmt.Sprintf("SELECT * FROM CART WHERE CART_ID = %d", Id)
		totalItems, totalPrice := resolver.db.CheckOut(request)
		request = fmt.Sprintf("DELETE FROM CART WHERE CART_ID = %d", Id)
		resolver.db.ModifyCart(request)
		return fmt.Sprintf("Successfully checked out cart %d with %d items with a total cost of %.2f", Id, totalItems, totalPrice), nil
	}
	return "Checkout failure", nil
}
package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/Owen-Cummings/ShopifyDevChallenge2019/mysql"
)

// struct to interact with graphQL objects
type Query struct{
	Query *graphql.Object
}

//creation of root query containing all query operations
func CreateRootQuery(db *mysql.Db) (*Query) {
	//resolver function from gql/resolvers.go to access the database
	resolver := Resolver{db: db}

	QueryType := Query{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{
					//Query to return products by either title, whether it's in stock, or both
					"GetProducts": &graphql.Field{
						Type: graphql.NewList(Product),
						Args: graphql.FieldConfigArgument{
							"InStock": &graphql.ArgumentConfig{
								Type: graphql.Boolean,
							},
							"Title": &graphql.ArgumentConfig{
								Type: graphql.String,
							},
						},
						Resolve: resolver.ProductResolver,
					},
					//Query to fetch cartentries with matching CartID
					"GetCartContents": &graphql.Field{
						Type: graphql.NewList(CartEntry),
						Args: graphql.FieldConfigArgument{
							"CartID": &graphql.ArgumentConfig{
								Type: graphql.Int,
							},
						},
						Resolve: resolver.CartResolver,
					},
				},
			},

		),

	}
	return &QueryType
}

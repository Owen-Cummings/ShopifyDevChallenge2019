package gql

import "github.com/graphql-go/graphql"

// configuring Product object for graphql to take sql product data
var Product = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
			"inventorycount": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var CartEntry = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CartEntry",
		Fields: graphql.Fields{
			"cartid": &graphql.Field{
				Type: graphql.Int,
			},
			"productid": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
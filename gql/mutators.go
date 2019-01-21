package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/Owen-Cummings/ShopifyDevChallenge2019/mysql"
)

type Mutation struct{
	Mutation *graphql.Object
}

func CreateRootMutator(db *mysql.Db) (*Mutation) {
	//resolver function from gql/resolvers.go to access the database
	resolver := Resolver{db: db}

	MutatorType := Mutation{
		Mutation: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Mutation",
				Fields: graphql.Fields{
					//AddToCart mutator using CartMutator resolver and taking a CartID or ItemID
					"AddToCart": &graphql.Field{
						Type: graphql.String,
						Args: graphql.FieldConfigArgument{
							"CartID": &graphql.ArgumentConfig{
								Type: graphql.Int,
							},
							"ItemID": &graphql.ArgumentConfig{
								Type: graphql.Int,
							},
						},
						Resolve: resolver.CartMutator,
					},
					//Checkout mutator using CheckOut resolver and taking a CartID then processing purchases
					"CheckOut": &graphql.Field{
						Type: graphql.String,
						Args: graphql.FieldConfigArgument{
							"CartID": &graphql.ArgumentConfig{
								Type: graphql.Int,
							},
						},
						Resolve: resolver.CheckOut,
					},

				},
			},
		),
	}
	return &MutatorType
}
package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
	"github.com/Owen-Cummings/ShopifyDevChallenge2019/gql"
	"github.com/Owen-Cummings/ShopifyDevChallenge2019/mysql"
	"github.com/Owen-Cummings/ShopifyDevChallenge2019/server"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	//create server connection and defer close to ensure it is done after all operations
	router, db := serverStart()

	defer db.Close()

	//Listen on port 4000 and log errors
	log.Fatal(http.ListenAndServeTLS(":4000", "server.crt", "server.key", router))
}

var tokenAuth *jwtauth.JWTAuth

func serverStart() (*chi.Mux, *mysql.Db){

	//hardcoded connection to mySQL server using database/sql for implementation
	db, err := mysql.CreateConnection(
		mysql.ConnInfo(
			"username", "password", "dbName",
			),
	)
	// if error is encountered log error, print, and shut down with os.Exit(1)
	if err != nil {
		log.Fatal(err)
	}

	//create a router to handle traffic with chi
	router := chi.NewRouter()

	//define the root query for graphQL using the database connection defined earlier
	rootQuery := gql.CreateRootQuery(db)
	rootMutator := gql.CreateRootMutator(db)

	//define graphql schema with root query
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: rootQuery.Query,
			Mutation: rootMutator.Mutation},
	)

	if err != nil{
		fmt.Println("Error creating schema: ", err)
	}

	srvr := server.Server{
		GqlSchema: &schema,
	}

	//middleware
	router.Use(
		render.SetContentType(render.ContentTypeJSON),	// content-type headers for API
		middleware.Throttle(4),					//limit actively processed requests to API to 4
		middleware.Logger,								// log api calls to API
		middleware.DefaultCompress,						// compress results
		middleware.StripSlashes,						// remove path trailing slash and continue MUX routing
		middleware.Recoverer,							//recover panics without server downtime
	)

	// private route for graphql server using JWT for authentication
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)
		router.Post("/go-graphql", srvr.GraphQL())
	})

	//create public route to generate JWT token
	router.Group(func(router chi.Router) {
		router.Get("/generate-jwt", func(w http.ResponseWriter, r *http.Request) {
			tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
			userID := rand.Intn(1000)
			_, newToken, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": userID})
			_, _ = w.Write([]byte(fmt.Sprintf("JWT_Token: %s", newToken)))
		})
	})

	//return router and database variables to be used for starting server in main
	return router, db
}

//sample token for testing API
func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": 12345})
	fmt.Printf("Sample JWT: %s\n\n", tokenString)
}


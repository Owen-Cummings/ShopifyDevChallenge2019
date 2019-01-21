package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
	"net/http"
)

// server struct holds db connection and handlers
type Server struct {
	GqlSchema *graphql.Schema
}

type requestBody struct {
	Request string `json:"request"`
}

func (server *Server) GraphQL() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Body == nil {
			http.Error(writer, "Must provide gql query in request body", 400)
			return
		}

		//declare request body struct to store request decoded in json format
		var reqBody requestBody
		err := json.NewDecoder(request.Body).Decode(&reqBody)
		if err != nil {
			http.Error(writer, "Error parsing JSON request body", 400)
		}

		//call and execute graphQL query
		result := Execute(reqBody.Request, *server.GqlSchema)

		// use chi/render function to marshal HTML data to JSON content
		render.JSON(writer, request, result)
	}
}

//function to execute graphQL queries
func Execute(query string, schema graphql.Schema) *graphql.Result {

	result := graphql.Do(graphql.Params{
		Schema: schema,
		RequestString: query,
	})
	//check errors associated with result
	if len(result.Errors)>0 {
		fmt.Printf("Unexpected errors inside Execute: %v", result.Errors)
	}

	return result
}
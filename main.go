package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"./db"
	"./api"
	"net/url"
	"io/ioutil"
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}


func main() {
	database := db.Connect()
	database.Begin()
	db.CreateUserTable(database)
	db.CreateStarTable(database)
	schema := api.StartAPI(database)

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Received Query: ")
		fmt.Println(r.URL.String())
		//fmt.Println(r.URL.Query().Get("query"))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
		fmt.Print("Result: ")
		fmt.Println(fmt.Sprint(result))
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Received Auth: ")
		code := r.URL.Query().Get("code")
		fmt.Println(code)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Exchange Code for Access-Token from Github
		resp, err := http.PostForm("https://github.com/login/oauth/access_token",
			url.Values{
				"client_id": {"8830bb5e86c0911377e0"},
				"client_secret": {"f1c4678365dee99553a18ddda9fa419d22d683e1"},
				"code": {code},
				"redirect_uri": {"http://localhost:3000/auth"},
		})
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		params, err := url.ParseQuery(string(body))
		accessToken := params.Get("access_token")
		fmt.Print("token: ")
		fmt.Println(accessToken)

		// select id or add user with token
		result := executeQuery("{user(token:\""+accessToken+"\"){id}}", schema)
		json.NewEncoder(w).Encode(result)
		fmt.Print("Result: ")
		fmt.Println(fmt.Sprint(result))
	})

	fmt.Println("Now server is running on port 4000")
	http.ListenAndServe(":4000", nil)
	db.Close(database)
}

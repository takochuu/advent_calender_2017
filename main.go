package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

var q graphql.ObjectConfig = graphql.ObjectConfig{
	Name: "query",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: resolveID,
		},
		"name": &graphql.Field{
			Type:    graphql.String,
			Resolve: resolveName,
		},
	},
}

var m graphql.ObjectConfig = graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "Params",
				Fields: graphql.Fields{
					"id": &graphql.Field{
						Type: graphql.Int,
					},
					"address": &graphql.Field{
						Type: graphql.NewObject(graphql.ObjectConfig{
							Name: "state",
							Fields: graphql.Fields{
								"state": &graphql.Field{
									Type: graphql.String,
								},
								"city": &graphql.Field{
									Type: graphql.String,
								},
							},
						}),
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// ここで更新処理をする
				return User{
					Id: 10000,
					Address: Address{
						State: "三宿",
						City:  "世田谷区",
					},
				}, nil
			},
		},
	},
}

type User struct {
	Id      int64   `json:"id"`
	Address Address `json:"address"`
}

type Address struct {
	State string `json:"state"`
	City  string `json:"city"`
}

var schemaConfig graphql.SchemaConfig = graphql.SchemaConfig{
	Query:    graphql.NewObject(q),
	Mutation: graphql.NewObject(m),
}
var schema, _ = graphql.NewSchema(schemaConfig)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	r := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(r.Errors) > 0 {
		fmt.Printf("エラーがあるよ: %v", r.Errors)
	}

	j, _ := json.Marshal(r)
	fmt.Printf("%s \n", j)

	return r
}

func handler(w http.ResponseWriter, r *http.Request) {
	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(r.Body)
	query := bufBody.String()

	result := executeQuery(query, schema)
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func resolveID(p graphql.ResolveParams) (interface{}, error) {
	return p.Args["id"], nil
}

func resolveName(p graphql.ResolveParams) (interface{}, error) {
	return "hoge", nil
}

package api

import (
	"github.com/graphql-go/graphql"
	"github.com/go-pg/pg"
	"math/rand"
	"time"
)


func StartAPI(db *pg.DB) graphql.Schema {
	database := db
	// Star object type with fields "id", "name", "user_id", "img" and "link"
	starType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "star",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"user_id": &graphql.Field{
					Type: graphql.Int,
				},
				"img": &graphql.Field{
					Type: graphql.String,
				},
				"link": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	// User object type with fields "id", "token", and it's related "star" objects
	userType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "user",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"token": &graphql.Field{
					Type: graphql.String,
				},
				"stars": &graphql.Field{
					Type: graphql.NewList(starType),
					Description: "List of posts by the user",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.Int,
							DefaultValue: 0,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						user, hasUser := p.Source.(User)
						if hasUser {
							return selectStarsByUser(user.ID, database), nil
						}
						return nil, nil
					},
				},
			},
		},
	)


	/*
	   	Query object type with fields "user", "star", "stars"(returns all of a user's stars), and
		"all_stars"(returns every users' stars)
	*/
	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type: userType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"token": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						userID, hasID := p.Args["id"].(int)
						token, hasToken := p.Args["token"].(string)
						if hasID {
							return selectUserByID(userID, database), nil
						} else if hasToken {
							return selectOrInsertUserByToken(token, database), nil
						} else {
							return selectAllStars(database), nil
						}
						return nil, nil
					},
				},
				"star": &graphql.Field{
					Type: starType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						starID, isOK := p.Args["id"].(int)
						if isOK {
							return selectStarByID(starID, database), nil
						}
						return nil, nil
					},
				},
				"stars": &graphql.Field{
					Type: starType,
					Args: graphql.FieldConfigArgument{
						"user_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						userID, hasID := p.Args["user_id"].(int)
						if hasID {
							return selectStarsByUser(userID, database), nil
						}
						return nil, nil
					},
				},
				"all_stars": &graphql.Field{
					Type: graphql.NewList(starType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return selectAllStars(database), nil
					},
				},
			},
		})


	/*
	   	Mutation object type with fields with CRUD for users and stars
	*/
	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			/* 	Create new User item
				?query=mutation{add_user(token:"somethingUnique"){id,token}}
			*/
			"add_user": &graphql.Field{
				Type:        userType,
				Description: "Add new user",
				Args: graphql.FieldConfigArgument{
					"token": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					rand.Seed(time.Now().UnixNano())
					user := User{
						Token:  params.Args["token"].(string),
					}
					insertUser(&user, database)
					return user, nil
				},
			},

			/* 	Create new star item
				?query=mutation{add_star(name:"",user_id:0,img:"",link:""){id,name,user_id,img,link}}
			*/
			"add_star": &graphql.Field{
				Type:        starType,
				Description: "Add new star",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"user_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"img": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"link": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					rand.Seed(time.Now().UnixNano())
					star := Star{
						Name:  params.Args["name"].(string),
						UserID:  params.Args["user_id"].(int),
						Img:  params.Args["img"].(string),
						Link:  params.Args["link"].(string),
					}
					insertStar(&star, database)
					return star, nil
				},
			},

			/* 	Update user item by id
				?query=mutation{update_user(id:"",token:""){id,token}}
			*/
			"update_user": &graphql.Field{
				Type:        userType,
				Description: "Update user by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"token": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, hasID := params.Args["id"].(int)
					token, hasToken := params.Args["token"].(string)
					user := User{}
					if hasID {
						user.ID = id
					}
					if hasToken {
						user.Token = token
					}
					ret := updateUser(user, database)
					return ret, nil
				},
			},

			/* 	Update star item by id
				?query=mutation{update_star(id:"",name:"",user_id:0,img:"",link:""){id,name,user_id,img,link}}
			*/
			"update_star": &graphql.Field{
				Type:        starType,
				Description: "Update star by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"user_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"img": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"link": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, hasID := params.Args["id"].(int)
					name, hasName := params.Args["name"].(string)
					userID, hasUserID := params.Args["user_id"].(int)
					img, hasImg := params.Args["img"].(string)
					link, hasLink := params.Args["link"].(string)
					star := Star{}
					if hasID {
						star.ID = id
					}
					if hasName {
						star.Name = name
					}
					if hasUserID {
						star.UserID = userID
					}
					if hasImg {
						star.Img = img
					}
					if hasLink {
						star.Link = link
					}
					ret := updateStar(star, database)
					return ret, nil
				},
			},

			/* 	Delete user item by id
				?query=mutation{delete_user(id:""){id,token}}
			*/
			"delete_user": &graphql.Field{
				Type:        userType,
				Description: "Delete user by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					user := User{ID:id}
					ret := deleteUser(user, database)
					return ret, nil
				},
			},

			/* 	Delete star item by id
				?query=mutation{delete_star(id:""){id,name,user_id,img,link}}
			*/
			"delete_star": &graphql.Field{
				Type:        starType,
				Description: "Delete star by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					star := Star{ID: id}
					ret := deleteStar(star, database)
					return ret, nil
				},
			},
		},
	})

	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	return schema
}
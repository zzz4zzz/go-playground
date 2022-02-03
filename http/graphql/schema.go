package graphql

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/zzz4zzz/go-playground/db"
	"github.com/zzz4zzz/go-playground/db/models"
)

var userType *graphql.Object
var queryType *graphql.Object
var Schema graphql.Schema

var createdAt = &graphql.Field{
	Type: graphql.DateTime,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		m, ok := p.Source.(models.ModelTimer)
		if ok {
			return m.GetCreatedAt(), nil
		}
		return nil, errors.New("failed to resolve created_at")
	},
}
var updatedAt = &graphql.Field{
	Type: graphql.DateTime,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		m, ok := p.Source.(models.ModelTimer)
		if ok {
			return m.GetUpdatedAt(), nil
		}
		return nil, errors.New("failed to resolve updated_at")
	},
}

func init() {
	bankAccountType := graphql.NewObject(graphql.ObjectConfig{
		Name: "BankAccount",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "Bank account name",
			},
			"number": &graphql.Field{
				Type:        graphql.String,
				Description: "Bank account number",
			},
			"created_at": createdAt,
			"updated_at": updatedAt,
		},
	})
	userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": createdAt,
			"updated_at": updatedAt,
			"bank_accounts": &graphql.Field{
				Type: graphql.NewList(bankAccountType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := p.Source.(models.User)
					return user.BankAccounts, nil
				},
			},
		},
	})

	queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var users []models.User
					db.DB.Preload("BankAccounts").Find(&users)
					return users, nil
				},
			},
		},
	})

	Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})

}

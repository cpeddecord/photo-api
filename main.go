package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cpeddecord/imgs-to-json"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var ImageList []imgstojson.ImgData

func init() {
	raw, err := ioutil.ReadFile("./images.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(raw, &ImageList)
}

var imageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Image",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"caption": &graphql.Field{
			Type: graphql.String,
		},
		"copyright": &graphql.Field{
			Type: graphql.String,
		},
		"createdDate": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"directory": &graphql.Field{
			Type: graphql.String,
		},
		"fNum": &graphql.Field{
			Type: graphql.String,
		},
		"filename": &graphql.Field{
			Type: graphql.String,
		},
		"focalLength": &graphql.Field{
			Type: graphql.String,
		},
		"imageHeight": &graphql.Field{
			Type: graphql.Int,
		},
		"imageWidth": &graphql.Field{
			Type: graphql.Int,
		},
		"iso": &graphql.Field{
			Type: graphql.Int,
		},
		"lens": &graphql.Field{
			Type: graphql.String,
		},
		"shutterSpeed": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"keywords": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"image": &graphql.Field{
			Type:        imageType,
			Description: "Get single Image",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					for _, image := range ImageList {
						if image.ID == idQuery {
							return image, nil
						}
					}
				}

				return imgstojson.ImgData{}, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func executeQuery(q string, s graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: q,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}

	return result
}

func main() {
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":3000", nil)
}

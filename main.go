package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cpeddecord/imgs-to-json"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var ImageList []imgstojson.ImgData

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

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
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"copyright": &graphql.Field{
			Type: graphql.String,
		},
		"createdDate": &graphql.Field{
			Type: graphql.String,
		},
		"keywords": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "list of keywords/tags",
		},
		"directory": &graphql.Field{
			Type: graphql.String,
		},
		"filename": &graphql.Field{
			Type: graphql.String,
		},
		"imageHeight": &graphql.Field{
			Type:        graphql.Int,
			Description: "Pixel Height",
		},
		"imageWidth": &graphql.Field{
			Type:        graphql.Int,
			Description: "Pixel Width",
		},
		"lens": &graphql.Field{
			Type:        graphql.String,
			Description: "Lens used",
		},
		"focalLength": &graphql.Field{
			Type:        graphql.String,
			Description: "Lens Focal Length",
		},
		"fNum": &graphql.Field{
			Type:        graphql.String,
			Description: "Lens F/number",
		},
		"shutterSpeed": &graphql.Field{
			Type:        graphql.String,
			Description: "Shutter Speed of the Camera",
		},
		"iso": &graphql.Field{
			Type:        graphql.Int,
			Description: "Camera's ISO/ASA Sensitivity",
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"image": &graphql.Field{
			Type:        imageType,
			Description: "Get single image",
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
		"imageList": &graphql.Field{
			Type:        graphql.NewList(imageType),
			Description: "Get multiple images",
			Args: graphql.FieldConfigArgument{
				"tag": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"tagContains": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				tagQuery, isTagOK := params.Args["tag"].(string)
				tagContains, isTagContainsOK := params.Args["tagContains"].(string)

				var imgs []imgstojson.ImgData

				if isTagContainsOK {
					for _, image := range ImageList {
						for _, s := range image.Keywords {
							if strings.Contains(s, tagContains) {
								imgs = append(imgs, image)
							}
						}
					}

					return imgs, nil
				}

				if isTagOK {
					for _, image := range ImageList {
						if contains(image.Keywords, tagQuery) {
							imgs = append(imgs, image)
						}
					}

					return imgs, nil

				}

				return ImageList, nil
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

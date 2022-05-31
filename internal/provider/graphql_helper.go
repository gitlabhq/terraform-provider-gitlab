//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xanzy/go-gitlab"
)

func SendGraphQLRequest(ctx context.Context, client *gitlab.Client, graphQLCall string, objectToParseForResponse interface{}) (interface{}, error) {

	request, err := client.NewRequest("POST", "", nil, []gitlab.RequestOptionFunc{func(request *retryablehttp.Request) error {
		//The "New Request" method automatically appends the "/api/v4" path onto the API, which needs to be replaced by "/api/graphql", so we need to use
		//a RequestOptionFunction to overwrite the URL.
		log.Print([]byte(graphQLCall))
		request.URL.Path = "/api/graphql"
		err := request.SetBody([]byte(graphQLCall))
		if err != nil {
			return err
		}

		return nil
	}})
	if err != nil {
		return nil, err
	}
	var testObject interface{}
	response, err := client.Do(request, testObject)
	if err != nil {
		return nil, err
	}

	test, _ := ioutil.ReadAll(response.Body)

	log.Printf("Resp Body:  %s", test)
	log.Printf("Past Do: %v", testObject)

	return objectToParseForResponse, nil
}

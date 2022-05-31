//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xanzy/go-gitlab"
)

func SendGraphQLRequest(ctx context.Context, client *gitlab.Client, graphQLCall string, objectToParseForResponse interface{}) (interface{}, error) {
	//The "New Request" method automatically appends the "/api/v4" path onto the API, which needs to be replaced by "/api/graphql", so we need to use
	//a RequestOptionFunction to overwrite the URL.
	request, err := client.NewRequest("POST", "", ctx, []gitlab.RequestOptionFunc{func(request *retryablehttp.Request) error {
		request, err := retryablehttp.NewRequestWithContext(ctx, "POST", client.BaseURL().RawPath, graphQLCall)
		if err != nil {
			return err
		}
		return nil
	}})
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request, nil)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	log.Fatalf("Body: %s", body)

	err = json.Unmarshal(body, objectToParseForResponse)
	if err != nil {
		return nil, err
	}

	return objectToParseForResponse, nil
}

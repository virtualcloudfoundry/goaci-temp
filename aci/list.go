package aci

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/virtualcloudfoundry/goaci/api"
)

// ListContainerGroups lists an Azure Container Instance Groups, if a resource
// group is given it will list by resource group.
// It optionally accepts a resource group name and will filter based off of it
// if it is not empty.
// From: https://docs.microsoft.com/en-us/rest/api/container-instances/containergroups/list
// From: https://docs.microsoft.com/en-us/rest/api/container-instances/containergroups/listbyresourcegroup
func (c *Client) ListContainerGroups(resourceGroup string) (*ContainerGroupListResult, error) {
	urlParams := url.Values{
		"api-version": []string{apiVersion},
	}

	// Create the url.
	uri := api.ResolveRelative(BaseURI, containerGroupListURLPath)
	// List by resource group if they passed one.
	if resourceGroup != "" {
		uri = api.ResolveRelative(BaseURI, containerGroupListByResourceGroupURLPath)

	}
	uri += "?" + url.Values(urlParams).Encode()

	// Create the request.
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("Creating get container group list uri request failed: %v", err)
	}

	// Add the parameters to the url.
	if err := api.ExpandURL(req.URL, map[string]string{
		"subscriptionId": c.auth.SubscriptionID,
		"resourceGroup":  resourceGroup,
	}); err != nil {
		return nil, fmt.Errorf("Expanding URL with parameters failed: %v", err)
	}
	return c.listContainersByReq(req)
}

func (c *Client) ListContainersGroupsByUri(uri string) (*ContainerGroupListResult, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("Creating get container group list uri request failed: %v", err)
	}

	return c.listContainersByReq(req)
}

func (c *Client) listContainersByReq(req *http.Request) (*ContainerGroupListResult, error) {

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Sending get container group list request failed: %v", err)
	}
	defer resp.Body.Close()

	// 200 (OK) is a success response.
	if err := api.CheckResponse(resp); err != nil {
		return nil, err
	}

	// Decode the body from the response.
	if resp.Body == nil {
		return nil, errors.New("Create container group list returned an empty body in the response")
	}
	var list ContainerGroupListResult
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("Decoding get container group response body failed: %v", err)
	}

	return &list, nil
}

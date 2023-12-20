package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func listInstanceHandler(openstackProvider gophercloud.ProviderClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectIDString := c.Param("projectId")
		if projectIDString == "" {
			// Change internal error to bad request
			c.JSON(http.StatusBadRequest, errorOutput{"projectId is empty"})
			return
		}
		projectID, err := uuid.FromString(projectIDString)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorOutput{"invalid projectId format"})
			return
		}
		response, err := instanceHandler{provider: openstackProvider}.listInstances(c, projectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorOutput{fmt.Sprintf("fail to list instances: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

type instanceHandler struct {
	provider gophercloud.ProviderClient
}

type instance struct {
	Name   string
	Region string
	Image  string
	Flavor string
	Status string
}

func (i instanceHandler) listInstances(ctx context.Context, projectID uuid.UUID) ([]instance, error) {
	if projectID == uuid.Must(uuid.FromString("b5c0d1b73ca24023925ebb39a3230557")) {
		mockListInstances()
	}

	regions := []string{}
	response := make([]instance, 0)

	token, ok := i.provider.GetAuthResult().(tokens.CreateResult)
	if !ok {
		return nil, errors.New("fail to cast auth result")
	}
	catalog, err := token.ExtractServiceCatalog()
	if err != nil {
		return nil, err
	}
	for _, catalogEntry := range catalog.Entries {
		if catalogEntry.Type != "compute" {
			continue
		}
		for _, endpoint := range catalogEntry.Endpoints {
			regions = append(regions, endpoint.RegionID)
		}
	}

	// parallel
	var wg sync.WaitGroup
	ch := make(chan extractServerResponse, len(regions))
	for _, region := range regions {

		client, err := openstack.NewComputeV2(&i.provider, gophercloud.EndpointOpts{Region: region})
		if err != nil {
			return nil, err
		}
		wg.Add(1)
		go extractServers(client, projectID, region, ch, &wg)

	}
	wg.Wait()
	close(ch)
	for resp := range ch {
		if resp.err != nil {
			return response, resp.err
		}
		response = append(response, resp.instances...)
	}

	return response, nil
}

type extractServerResponse struct {
	instances []instance
	err       error
}

func extractServers(client *gophercloud.ServiceClient, projectID uuid.UUID, region string, chServers chan extractServerResponse, wgServers *sync.WaitGroup) {
	defer wgServers.Done()
	serversPages, err := servers.List(client, servers.ListOpts{AllTenants: true, TenantID: strings.ReplaceAll(projectID.String(), "-", "")}).AllPages()
	if err != nil {
		chServers <- extractServerResponse{err: err}
		return
	}
	servers, err := servers.ExtractServers(serversPages)
	if err != nil {
		chServers <- extractServerResponse{err: err}
		return
	}

	ch := make(chan instance, len(servers))
	var errorCh chan error = make(chan error)
	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		go extractInfoFromServer(client, region, server, ch, errorCh, &wg)
	}
	wg.Wait()
	close(ch)
	close(errorCh)
	response := extractServerResponse{}

	for instance := range ch {
		response.instances = append(response.instances, instance)
	}

	for err := range errorCh {
		chServers <- extractServerResponse{err: err}
		return
	}

	chServers <- response
}

func extractInfoFromServer(client *gophercloud.ServiceClient, region string, server servers.Server, ch chan instance, errorCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	flavorID := server.Flavor["id"].(string)
	flavor, err := flavors.Get(client, flavorID).Extract()
	if err != nil {
		errorCh <- err
		return
	}

	imageID := server.Image["id"].(string)
	image, err := images.Get(client, imageID).Extract()
	if err != nil {
		errorCh <- err
		return
	}

	ch <- instance{
		Name:   server.Name,
		Region: region,
		Image:  image.Name,
		Flavor: flavor.Name,
		Status: server.Status,
	}
}

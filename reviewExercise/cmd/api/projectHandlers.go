package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"golang.org/x/sync/errgroup"
)

func listInstanceHandler(openstackProvider gophercloud.ProviderClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectIDString := c.Param("projectId")
		if projectIDString == "" {
			// Wrong error code returned
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
			// Maybe hide some data in error ?
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

func (i instanceHandler) getRegions() ([]string, error) {
	regions := []string{}

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
	return regions, nil
}

func (i instanceHandler) getServerImage(client *gophercloud.ServiceClient, region string, projectID uuid.UUID, imageID string) (string, error) {

	image, err := images.Get(client, imageID).Extract()
	if err != nil {
		return "", err
	}

	return image.Name, nil
}

func (i instanceHandler) getServerFlavor(client *gophercloud.ServiceClient, region string, projectID uuid.UUID, flavorID string) (string, error) {

	flavor, err := flavors.Get(client, flavorID).Extract()
	if err != nil {
		return "", err
	}
	return flavor.Name, nil
}

func (i instanceHandler) listServers(region string, projectID uuid.UUID) ([]instance, error) {
	response := make([]instance, 0)
	client, err := openstack.NewComputeV2(&i.provider, gophercloud.EndpointOpts{Region: region})
	if err != nil {
		return nil, err
	}
	// rework ReplaceAll in another function
	serversPages, err := servers.List(client, servers.ListOpts{AllTenants: true, TenantID: strings.ReplaceAll(projectID.String(), "-", "")}).AllPages()
	if err != nil {
		return nil, err
	}
	servers, err := servers.ExtractServers(serversPages)
	if err != nil {
		return nil, err
	}

	// Keep map of image/flavor to avoid duplicate api call
	for _, server := range servers {

		// Catch error cast
		imageID := server.Image["id"].(string)
		image, err := i.getServerImage(client, region, projectID, imageID)
		if err != nil {
			return nil, err
		}

		flavorID := server.Flavor["id"].(string)
		flavor, err := i.getServerFlavor(client, region, projectID, flavorID)
		if err != nil {
			return nil, err
		}
		response = append(response, instance{
			Name:   server.Name,
			Region: region,
			Image:  image,
			Flavor: flavor,
			Status: server.Status,
		})
	}
	return response, nil
}

// Decoupage DDD
func (i instanceHandler) listInstances(ctx context.Context, projectID uuid.UUID) ([]instance, error) {
	if projectID == uuid.Must(uuid.FromString("b5c0d1b73ca24023925ebb39a3230557")) {
		mockListInstances()
	}

	regions, err := i.getRegions()
	if err != nil {
		return nil, err
	}

	response := make([]instance, 0)
	// Slow
	// Parallelization ?
	// Limit page ?
	// Regionalize API ?
	// Optimize api call with field selection
	//countWorkers := 10
	errGroup := new(errgroup.Group)
	result := make(chan instance, len(regions))
	for _, region := range regions {

		errGroup.Go(func() error {
			instances, err := i.listServers(region, projectID)
			if err != nil {
				return err
			}
			for _, instance := range instances {
				result <- instance
			}
			return nil
		})
	}
	// Wait for all HTTP fetches to complete.
	if err := errGroup.Wait(); err != nil {
		return nil, err
	}
	// Consistance ou résultat partiel ?
	for instance := range result {
		response = append(response, instance)
	}
	return response, nil
}

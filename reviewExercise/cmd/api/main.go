package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/h2non/gock"
	"github.com/snigle/dddExercise/reviewExercise/bootstrap"
	"github.com/snigle/dddExercise/reviewExercise/ui"
)

func main() {
	c, _ := bootstrap.Constructor()

	cloudProviderConf := c.CloudProviderConf
	cloudProvider:= 
	// On crére l'API service, en passant comme argument un client d'open stack qui est instanticé avant
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://openstack.example.com:5000/v3/",
		Username:         "username",
		Password:         "password",
		DomainID:         "527e8ff13ea64fa7a70bb62dfe37ac47",
	}
	mockToken()
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatalf("fail to init openstack provider: %s", err)
	}
	gock.InterceptClient(&provider.HTTPClient)

	r := gin.Default()
	registerRoutes(r, *provider)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func registerRoutes(router *gin.Engine, provider gophercloud.ProviderClient) {
	router.GET("/cloud/project/:projectId/instance", ui.ListInstanceHandler(provider))
}

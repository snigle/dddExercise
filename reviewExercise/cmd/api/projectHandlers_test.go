package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

// ProjectHandlerSuite is a test suite for ProjectHandler.
type ProjectHandlerSuite struct {
	suite.Suite
	router *gin.Engine
}

// TestProjectHandler runs ProjectHandlerSuite.
func TestProjectHandler(t *testing.T) {
	suite.Run(t, &ProjectHandlerSuite{})
}

// SetupTest setups each test.
func (s *ProjectHandlerSuite) SetupTest() {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://openstack.example.com:5000/v3/",
		Username:         "mock",
		Password:         "mock",
		DomainID:         "527e8ff13ea64fa7a70bb62dfe37ac47",
	}
	gock.New("https://openstack.example.com:5000").Post("/v3/auth/tokens").Reply(201).JSON(`
	{
		"token": {
			"catalog": [
				{
					"endpoints": [
						{
							"id": "compute-endpoint-mock-id",
							"interface": "public",
							"legacy_endpoint_id": "",
							"region_id": "MOCK",
							"url": "https://compute.mock.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d"
						}
					],
					"id": "compute-catalog-mock-id",
					"type": "compute"
				},
				{
					"endpoints": [
						{
							"id": "image-endpoint-mock-id",
							"interface": "public",
							"legacy_endpoint_id": "",
							"region_id": "MOCK",
							"url": "https://image.compute.mock.cloud.ovh.net/"
						}
					],
					"id": "image-catalog-mock-id",
					"type": "image"
				}
			],
			"expires_at": "2017-09-23T11:33:05.371752Z",
			"issued_at": "2017-09-22T11:33:05.371773Z"
		}
	}
	`)
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		s.Require().NoError(err)
	}
	gock.InterceptClient(&provider.HTTPClient)

	s.router = gin.Default()
	registerRoutes(s.router, *provider)
}

// TearDownTest tears down each test.
func (s *ProjectHandlerSuite) TearDownTest() {
}

// TestStuff tests stuff.
func (s *ProjectHandlerSuite) TestListBadParam() {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cloud/project/toto/instance", nil)
	s.router.ServeHTTP(w, req)

	s.Require().Equal(400, w.Code)
	expected, err := json.Marshal(errorOutput{"invalid projectId format"})
	s.Require().NoError(err)
	s.Require().Equal(string(expected), w.Body.String())

}

// TODO add test for all errors cases

func (s *ProjectHandlerSuite) TestListSuccess() {

	gock.New("https://compute.mock.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/servers/detail").
		MatchParams(map[string]string{
			"all_tenants": "true",
			"tenant_id":   "4a08acc431d840b9ac9a1e60c5f0477c",
		}).Reply(200).JSON(`
	{"servers":[
		{
			"status": "ACTIVE",
			"key_name": "key-name",
			"image": {
				"id": "2fc5e2cb-28ad-4c8d-a723-098f42b8b288"
			},
			"flavor": {
				"id": "031b7771-c514-4e6d-af77-c1ea52f6fc38"
			},
			"id": "e6e8a496-5f32-4259-9877-c022429bff23",
			"name": "server-name"
		}
	]}
	`)

	gock.New("https://compute.mock.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/flavors/031b7771-c514-4e6d-af77-c1ea52f6fc38").Reply(200).JSON(
		`{
			"flavor": {
				"name": "small"
			}
		}`,
	)

	gock.New("https://compute.mock.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/images/2fc5e2cb-28ad-4c8d-a723-098f42b8b288").Reply(200).JSON(
		`{
			"image": {
				"name": "ubuntu"
			}
		}`,
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cloud/project/4a08acc431d840b9ac9a1e60c5f0477c/instance", nil)
	s.router.ServeHTTP(w, req)

	body := w.Body.String()
	s.Require().Equal(200, w.Code, body)
	expected, err := json.Marshal([]instance{{
		Name:   "server-name",
		Region: "MOCK",
		Image:  "ubuntu",
		Flavor: "small",
		Status: "ACTIVE",
	}})
	s.Require().NoError(err)
	s.Require().Equal(string(expected), body)

}

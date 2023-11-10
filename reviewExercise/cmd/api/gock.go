package main

import (
	"time"

	"github.com/h2non/gock"
)

func mockToken() {
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
							"region_id": "GRA1",
							"url": "https://compute.gra1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d"
						},
						{
							"id": "compute-endpoint-mock2-id",
							"interface": "public",
							"legacy_endpoint_id": "",
							"region_id": "DE1",
							"url": "https://compute.de1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d"
						},
						{
							"id": "compute-endpoint-mock2-id",
							"interface": "public",
							"legacy_endpoint_id": "",
							"region_id": "SYD1",
							"url": "https://compute.syd1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d"
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
							"region_id": "GRA1",
							"url": "https://image.compute.gra1.cloud.ovh.net/"
						},
						{
							"id": "image-endpoint-mock2-id",
							"interface": "public",
							"legacy_endpoint_id": "",
							"region_id": "DE1",
							"url": "https://image.compute.de1.cloud.ovh.net/"
						},
						{
							"id": "image-endpoint-mock2-id",
							"interface": "public",
							"legacy_endpoint_id": "",
							"region_id": "SYD1",
							"url": "https://image.compute.syd1.cloud.ovh.net/"
						}
					],
					"id": "image-catalog-mock-id",
					"type": "image"
				}
			],
			"expires_at": "2017-09-23T11:33:05.371752Z",
			"issued_at": "2017-09-22T11:33:05.371773Z",
			"methods": [
				"manager"
			],
			"project": {
				"domain": {
					"name": "Default"
				},
				"id": "639ff29120f5458d92ddc4063ed8374d",
				"name": "2635455855851552"
			},
			"roles": [
				{
					"id": "administrator-role-id",
					"name": "administrator"
				}
			],
			"user": {
				"domain": {
					"name": "Default"
				},
				"name": "639ff29120f5458d92ddc4063ed8374d",
				"id": "user-id"
			}
		}
	}
	`)
}

func mockListInstances() {

	// GET SERVERS

	gock.New("https://compute.gra1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/servers/detail").
		MatchParams(map[string]string{
			"all_tenants": "true",
			"tenant_id":   "b5c0d1b73ca24023925ebb39a3230557",
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

	gock.New("https://compute.de1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/servers/detail").
		MatchParams(map[string]string{
			"all_tenants": "true",
			"tenant_id":   "b5c0d1b73ca24023925ebb39a3230557",
		}).Reply(200).JSON(`
	{"servers":[
		{
			"status": "ACTIVE",
			"key_name": "key-name",
			"image": {
				"id": "bc072e96-bf59-4418-a472-144dc74146e6"
			},
			"flavor": {
				"id": "509eae67-0b0b-40a1-a48d-f3173d828ac9"
			},
			"id": "e6e8a496-5f32-4259-9877-c022429bff23",
			"name": "server-name"
		}
	]}
	`)

	gock.New("https://compute.syd1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/servers/detail").
		MatchParams(map[string]string{
			"all_tenants": "true",
			"tenant_id":   "b5c0d1b73ca24023925ebb39a3230557",
		}).Reply(200).JSON(`
	{"servers":[]}
	`).Delay(time.Second)

	// GET Flavors
	gock.New("https://compute.gra1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/flavors/031b7771-c514-4e6d-af77-c1ea52f6fc38").Reply(200).JSON(
		`{
			"flavor": {
				"name": "small"
			}
		}`,
	)
	gock.New("https://compute.de1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/flavors/509eae67-0b0b-40a1-a48d-f3173d828ac9").Reply(200).JSON(
		`{
			"flavor": {
				"name": "medium"
			}
		}`,
	)

	// GET Images
	gock.New("https://compute.gra1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/images/2fc5e2cb-28ad-4c8d-a723-098f42b8b288").Reply(200).JSON(
		`{
			"image": {
				"name": "ubuntu"
			}
		}`,
	)
	gock.New("https://compute.de1.cloud.ovh.net/v2/639ff29120f5458d92ddc4063ed8374d").Get("/images/bc072e96-bf59-4418-a472-144dc74146e6").Reply(200).JSON(
		`{
			"image": {
				"name": "debian"
			}
		}`,
	)

}

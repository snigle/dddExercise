package bootstrap

import (
	"errors"

	"github.com/gophercloud/gophercloud"
)

type Config struct {
	openStackConf OpenStackConf
}

type OpenStackConf struct {
	IdentityEndpoint string
	Username         string
	Password         string
	DomainID         string
}

func Constructor() (Config, error) {

	openStackConf := OpenStackConf{
		IdentityEndpoint: "https://openstack.example.com:5000/v3/",
		Username:         "username",
		Password:         "password",
		DomainID:         "527e8ff13ea64fa7a70bb62dfe37ac47",
	}
	if openStackConf.Username == "" {
		return Config{}, errors.New("missing username")
	}
	// FAIL FAST
	return Config{openStackConf: openStackConf}, nil
}

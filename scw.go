package main

import (
	"os"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/config"
)

var endpoint *api.ScalewayAPI
var gateway = os.Getenv("SCW_GATEWAY")

func init() {
	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	endpoint, err = api.NewScalewayAPI(config.Organization, config.Token, "scw-update", "par1") // endpoints seems to be unified now....
	endpoint.Logger = api.NewDisableLogger()
	if err != nil {
		panic(err)
	}
}

func getServers() []api.ScalewayServer {
	s, _ := endpoint.GetServers(true, 0)
	return *s
}

func getImage(name, arch string) (*api.ScalewayImage, error) {
	imageID, err := endpoint.GetImageID(name, arch)
	if err != nil {
		return nil, err
	}
	return endpoint.GetImage(imageID.Identifier)
}

func terminateServer(id string) {
	endpoint.DeleteServerForce(id)
	api.WaitForServerStopped(endpoint, id)
}

func createServer(name string, comercialType string, image *string, ipv4 string, hasIPv6 bool, tags []string) error {
	f := false
	id, err := endpoint.PostServer(api.ScalewayServerDefinition{
		Name:              name,
		Tags:              tags,
		CommercialType:    comercialType,
		Image:             image,
		PublicIP:          ipv4,
		EnableIPV6:        hasIPv6,
		DynamicIPRequired: &f,
	})
	if err != nil {
		return err
	}
	api.StartServer(endpoint, id, true)
	api.WaitForServerReady(endpoint, id, gateway)
	return nil
}

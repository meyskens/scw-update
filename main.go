package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

func main() {
	fmt.Println("Scaleway mass updater")
	fmt.Println("Copyright 2018 Maartje Eyskens")
	fmt.Println("==============================")

	if len(os.Args) < 3 {
		fmt.Println("\nUsage: scw-update [name-regex-match] [new-image-name]")
		return
	}

	matchedServers := []api.ScalewayServer{}

	for _, server := range getServers() {
		if matched, _ := regexp.MatchString(os.Args[1], server.Name); matched {
			matchedServers = append(matchedServers, server)
		}
	}

	for _, server := range matchedServers {
		replaceServer(server)
	}
}

func replaceServer(server api.ScalewayServer) {
	image, err := getImage(os.Args[2], server.Arch)
	if err != nil {
		fmt.Println("Image not found for", server.Arch)
		return
	}

	name := server.Name
	commercialType := server.CommercialType
	tags := server.Tags
	hasIPv6 := server.EnableIPV6
	v4ID := ""
	if v4 := server.PublicAddress; v4.Dynamic != nil {
		v4ID = v4.Identifier
	}

	fmt.Println("Deleting", name)
	terminateServer(server.Identifier)
	fmt.Println("Recreating", name)
	createServer(name, commercialType, &image.Identifier, v4ID, hasIPv6, tags)

}

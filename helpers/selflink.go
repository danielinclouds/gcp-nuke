package helpers

import (
	"net/url"
	"regexp"
	"strings"
)

type SubnetworkSelfLink struct {
	Project string
	Region  string
	Name    string
}

// TODO:
// 1. Handle problems with regexp, for example when an empty string is returned

// Function ParseSubnetworkSelfLink parses Subnetwork selfLink into SubnetworkSelfLink struct
// e.g. https://www.googleapis.com/compute/v1/projects/test-123/regions/europe-west1/subnetworks/sub1
// SubnetworkSelfLink {
//     Project: test-123
//     Region:  europe-west1
//     Name:    sub1
// }
func ParseSubnetworkSelfLink(u string) (SubnetworkSelfLink, error) {

	SubnetworkSelfLinkUrl, err := url.Parse(u)
	if err != nil {
		return SubnetworkSelfLink{}, err
	}

	projectsChunk := regexp.MustCompile(`/projects/.+?/`).Find([]byte(SubnetworkSelfLinkUrl.Path))
	project := strings.TrimPrefix(string(projectsChunk), "/projects/")
	project = strings.TrimSuffix(project, "/")

	regionsChunk := regexp.MustCompile(`/regions/.+?/`).Find([]byte(SubnetworkSelfLinkUrl.Path))
	region := strings.TrimPrefix(string(regionsChunk), "/regions/")
	region = strings.TrimSuffix(region, "/")

	resourcesChunk := regexp.MustCompile(`/subnetworks/.+`).Find([]byte(SubnetworkSelfLinkUrl.Path))
	resourceName := strings.TrimPrefix(string(resourcesChunk), "/subnetworks/")
	resourceName = strings.TrimSuffix(resourceName, "/")

	return SubnetworkSelfLink{
			Project: project,
			Region:  region,
			Name:    resourceName,
		},
		nil

}

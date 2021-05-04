package helpers

import (
	"net/url"
	"regexp"
	"strings"
)

type SubnetworkSelfLink struct {
	Projects     string
	Regions      string
	ResourceName string
}

// TODO:
// 1. Handle problems with regexp, for example when an empty string is returned

// Function ParseSubnetworkSelfLink parses Subnetwork selfLink into SubnetworkSelfLink struct
// e.g. https://www.googleapis.com/compute/v1/projects/test-123/regions/europe-west1/subnetworks/sub1
// SubnetworkSelfLink {
//     Projects:     test-123
//     Regions:      europe-west1
//     ResourceName: sub1
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
			Projects:     project,
			Regions:      region,
			ResourceName: resourceName,
		},
		nil

}

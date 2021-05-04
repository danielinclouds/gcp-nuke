package gcp

import (
	"context"
	"os"

	"github.com/danielinclouds/gcp-nuke/helpers"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListVPC(projectID string, credJSON []byte) {

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("gcp-nuke.json"))
	if err != nil {
		panic(err.Error())
	}

	networkListCall := computeService.Networks.List(projectID)
	networkList, err := networkListCall.Do()
	if err != nil {
		panic(err.Error())
	}

	for _, n := range networkList.Items {
		log.Infof("Network: %s", n.Name)
	}

}

func DeleteAllVPC(projectID string, credJSON []byte) {

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("gcp-nuke.json"))
	if err != nil {
		panic(err.Error())
	}

	networks, err := computeService.Networks.List(projectID).Context(ctx).Do()
	if err != nil {
		panic(err.Error())
	}

	for _, n := range networks.Items {

		// TODO:
		// 1. Don't skip default network
		if n.Name == "default" {
			log.Debug("Skipping default network")
			continue
		}

		log.Debugf("Deleting network: %s", n.Name)

		deleteAllSubnetworks(credJSON, n.Subnetworks)
		deleteVPC(projectID, credJSON, n.Name)
	}
}

func deleteVPC(projectID string, credJSON []byte, network string) {

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("gcp-nuke.json"))
	if err != nil {
		panic(err.Error())
	}

	networkDeleteCall := computeService.Networks.Delete(projectID, "daniel")
	operation, err := networkDeleteCall.Do()
	if err != nil {
		panic(err.Error())
	}

	gresp, err := computeService.GlobalOperations.Wait(projectID, operation.Name).Context(ctx).Do()
	if err != nil {
		panic(err.Error())
	}

	if gresp.Error != nil {
		for _, m := range gresp.Error.Errors {
			log.Error(m.Message)
		}
		panic("FAILED")
	}

}

func deleteAllSubnetworks(credJSON []byte, subnetworks []string) {

	for _, subnetwork := range subnetworks {
		sub, err := helpers.ParseSubnetworkSelfLink(subnetwork)
		if err != nil {
			panic(err.Error())
		}

		log.Debugf("Deleting subnetwork: %s", sub.ResourceName)
		deleteSubnetwork(credJSON, sub)
	}

}

func deleteSubnetwork(credJSON []byte, subnetwork helpers.SubnetworkSelfLink) {

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("gcp-nuke.json"))
	if err != nil {
		panic(err.Error())
	}

	operation, err := computeService.Subnetworks.
		Delete(subnetwork.Projects, subnetwork.Regions, subnetwork.ResourceName).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	resp, err := computeService.RegionOperations.
		Wait(subnetwork.Projects, subnetwork.Regions, operation.Name).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	if resp.Error != nil {
		for _, m := range resp.Error.Errors {
			log.Error(m.Message)
		}
		panic("FAILED")
	}

}

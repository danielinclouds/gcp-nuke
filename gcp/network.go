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

func ListVPC(projectId string, credJSON []byte) {

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("gcp-nuke.json"))
	if err != nil {
		panic(err.Error())
	}

	networkListCall := computeService.Networks.List(projectId)
	networkList, err := networkListCall.Do()
	if err != nil {
		panic(err.Error())
	}

	for _, n := range networkList.Items {
		log.Infof("Network: %s", n.Name)
	}

}

func DeleteAllVPC(projectId string, credJSON []byte) {

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("gcp-nuke.json"))
	if err != nil {
		panic(err.Error())
	}

	networks, err := computeService.Networks.List(projectId).Context(ctx).Do()
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
		deleteVPC(projectId, credJSON, n.Name)
	}
}

func deleteVPC(projectId string, credJSON []byte, network string) {

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("gcp-nuke.json"))
	if err != nil {
		panic(err.Error())
	}

	networkDeleteCall := computeService.Networks.Delete(projectId, "daniel")
	operation, err := networkDeleteCall.Do()
	if err != nil {
		panic(err.Error())
	}

	gresp, err := computeService.GlobalOperations.Wait(projectId, operation.Name).Context(ctx).Do()
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

		log.Debugf("Deleting subnetwork: %s", sub.Name)
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
		Delete(subnetwork.Project, subnetwork.Region, subnetwork.Name).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	resp, err := computeService.RegionOperations.
		Wait(subnetwork.Project, subnetwork.Region, operation.Name).
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

package gcp

import (
	"context"

	"github.com/danielinclouds/gcp-nuke/config"
	"github.com/danielinclouds/gcp-nuke/helpers"
)

func ListVPC(cfg *config.Config) {

	networkListCall := cfg.ComputeService.Networks.List(cfg.Project)
	networkList, err := networkListCall.Do()
	if err != nil {
		panic(err.Error())
	}

	for _, n := range networkList.Items {
		cfg.Log.Infof("Network: %s", n.Name)
	}

}

func DeleteAllVPC(cfg *config.Config) {

	networks, err := cfg.ComputeService.Networks.List(cfg.Project).Context(context.Background()).Do()
	if err != nil {
		panic(err.Error())
	}

	for _, n := range networks.Items {

		// TODO:
		// 1. Don't skip default network
		if n.Name == "default" {
			cfg.Log.Debug("Skipping default network")
			continue
		}

		cfg.Log.Debugf("Deleting network: %s", n.Name)

		deleteAllSubnetworks(cfg, n.Subnetworks)
		deleteVPC(cfg, n.Name)
	}
}

func deleteVPC(cfg *config.Config, network string) {

	networkDeleteCall := cfg.ComputeService.Networks.Delete(cfg.Project, "daniel")
	operation, err := networkDeleteCall.Do()
	if err != nil {
		panic(err.Error())
	}

	gresp, err := cfg.ComputeService.GlobalOperations.
		Wait(cfg.Project, operation.Name).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	if gresp.Error != nil {
		for _, m := range gresp.Error.Errors {
			cfg.Log.Error(m.Message)
		}
		panic("FAILED")
	}

}

func deleteAllSubnetworks(cfg *config.Config, subnetworks []string) {

	for _, subnetwork := range subnetworks {
		sub, err := helpers.ParseSubnetworkSelfLink(subnetwork)
		if err != nil {
			panic(err.Error())
		}

		cfg.Log.Debugf("Deleting subnetwork: %s", sub.Name)
		deleteSubnetwork(cfg, sub)
	}

}

func deleteSubnetwork(cfg *config.Config, subnetwork helpers.SubnetworkSelfLink) {

	operation, err := cfg.ComputeService.Subnetworks.
		Delete(subnetwork.Project, subnetwork.Region, subnetwork.Name).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	resp, err := cfg.ComputeService.RegionOperations.
		Wait(subnetwork.Project, subnetwork.Region, operation.Name).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	if resp.Error != nil {
		for _, m := range resp.Error.Errors {
			cfg.Log.Error(m.Message)
		}
		panic("FAILED")
	}

}

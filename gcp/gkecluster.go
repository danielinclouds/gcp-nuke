package gcp

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/danielinclouds/gcp-nuke/config"
)

func ListGKEClusters(cfg *config.Config) {
	if isServiceDisabled(cfg, "container.googleapis.com") {
		cfg.Log.Debug("Kubernetes Engine API is disabled")
		return
	}

	resp, err := cfg.ContainerService.Projects.Locations.Clusters.
		List(fmt.Sprintf("projects/%s/locations/-", cfg.Project)).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	for _, cluster := range resp.Clusters {
		cfg.Log.Infof("GKE Cluster: %s", cluster.SelfLink)
	}

}

func DeleteAllGKEClusters(cfg *config.Config) {
	if isServiceDisabled(cfg, "container.googleapis.com") {
		cfg.Log.Debug("Kubernetes Engine API is disabled")
		return
	}

	resp, err := cfg.ContainerService.Projects.Locations.Clusters.
		List(fmt.Sprintf("projects/%s/locations/-", cfg.Project)).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	var wg sync.WaitGroup
	for _, cluster := range resp.Clusters {
		wg.Add(1)
		go deleteGKECluster(cluster.SelfLink, cfg, &wg)
	}

	wg.Wait()
}

func deleteGKECluster(selfLink string, cfg *config.Config, wg *sync.WaitGroup) {

	defer wg.Done()

	selfLinkUrl, err := url.Parse(selfLink)
	if err != nil {
		panic(err.Error())
	}
	clusterName := strings.TrimPrefix(selfLinkUrl.Path, "/v1/")

	resp, err := cfg.ContainerService.Projects.Locations.Clusters.
		Delete(clusterName).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	operationUrl, err := url.Parse(resp.SelfLink)
	if err != nil {
		panic(err.Error())
	}
	operationName := strings.TrimPrefix(operationUrl.Path, "/v1/")

	done := make(chan bool, 1)
	go func() {

		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {

			operation, err := cfg.ContainerService.Projects.Locations.Operations.
				Get(operationName).
				Context(context.Background()).
				Do()
			if err != nil {
				panic(err.Error())
			}
			cfg.Log.Debugf("Deleting cluster %s status: %s", clusterName, operation.Status)

			if operation.Status == "DONE" {
				done <- true
				ticker.Stop()
				break
			}
		}
	}()

	select {
	case <-done:
		cfg.Log.Debugf("Finished deleting cluster: %s", clusterName)

	case <-time.After(3 * time.Minute):
		cfg.Log.Debugf("Timeout while deleting cluster: %s", clusterName)
	}
}

package gcp

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/danielinclouds/gcp-nuke/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListGKEClusters(cfg *config.Config) {
	if isServiceDisabled(cfg, "container.googleapis.com") {
		log.Debug("Kubernetes Engine API is disabled")
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
		log.Infof("GKE Cluster: %s", cluster.SelfLink)
	}

}

func DeleteAllGKEClusters(cfg *config.Config) {
	if isServiceDisabled(cfg, "container.googleapis.com") {
		log.Debug("Kubernetes Engine API is disabled")
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
			log.Debugf("Deleting cluster %s status: %s", clusterName, operation.Status)

			if operation.Status == "DONE" {
				done <- true
				ticker.Stop()
				break
			}
		}
	}()

	select {
	case <-done:
		log.Debugf("Finished deleting cluster: %s", clusterName)

	case <-time.After(3 * time.Minute):
		log.Debugf("Timeout while deleting cluster: %s", clusterName)
	}
}

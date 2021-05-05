package gcp

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListGKEClusters(projectId string, credJSON []byte) {

	ctx := context.Background()
	containerService, err := container.NewService(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	resp, err := containerService.Projects.Locations.Clusters.
		List(fmt.Sprintf("projects/%s/locations/-", projectId)).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	for _, cluster := range resp.Clusters {
		log.Infof("GKE Cluster: %s", cluster.SelfLink)
	}

}

func DeleteAllGKEClusters(projectId string, credJSON []byte) {

	ctx := context.Background()
	containerService, err := container.NewService(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	resp, err := containerService.Projects.Locations.Clusters.
		List(fmt.Sprintf("projects/%s/locations/-", projectId)).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	var wg sync.WaitGroup
	for _, cluster := range resp.Clusters {
		wg.Add(1)
		go deleteGKECluster(cluster.SelfLink, credJSON, &wg)
	}

	wg.Wait()
}

func deleteGKECluster(selfLink string, credJSON []byte, wg *sync.WaitGroup) {

	defer wg.Done()

	ctx := context.Background()
	containerService, err := container.NewService(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	selfLinkUrl, err := url.Parse(selfLink)
	if err != nil {
		panic(err.Error())
	}
	clusterName := strings.TrimPrefix(selfLinkUrl.Path, "/v1/")

	resp, err := containerService.Projects.Locations.Clusters.Delete(clusterName).Context(ctx).Do()
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

			operation, err := containerService.Projects.Locations.Operations.Get(operationName).Context(ctx).Do()
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

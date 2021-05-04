package gcp

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/serviceusage/v1"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListNonDefaultServices(projectID string, credJSON []byte) {

	ctx := context.Background()
	serviceusageService, err := serviceusage.NewService(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	// TODO: Handle multiple pages
	resp, err := serviceusageService.Services.
		List(fmt.Sprintf("projects/%s", projectID)).
		Filter("state:ENABLED").
		PageSize(200).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	var enabledServices []string
	for _, service := range resp.Services {
		enabledServices = append(enabledServices, service.Name)
	}

	enabledServices = removeDefaultServices(projectID, credJSON, enabledServices)
	for _, service := range enabledServices {
		log.Infof("API Service: %s", service)
	}

}

func DisableAllNonDefaultServices(projectID string, credJSON []byte) {

	ctx := context.Background()
	serviceusageService, err := serviceusage.NewService(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	// TODO: Handle multiple pages
	resp, err := serviceusageService.Services.
		List(fmt.Sprintf("projects/%s", projectID)).
		Filter("state:ENABLED").
		PageSize(200).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	var enabledServices []string
	for _, service := range resp.Services {
		enabledServices = append(enabledServices, service.Name)
	}

	enabledServices = removeDefaultServices(projectID, credJSON, enabledServices)
	for _, service := range enabledServices {
		disableService(service, credJSON)
	}

}

func disableService(serviceName string, credJSON []byte) {

	ctx := context.Background()
	serviceusageService, err := serviceusage.NewService(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	log.Debugf("Disable service: %s", serviceName)
	operation, err := serviceusageService.Services.
		Disable(serviceName, &serviceusage.DisableServiceRequest{DisableDependentServices: true}).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	// Handle when parent service already disabled child service
	if operation.Name == "operations/noop.DONE_OPERATION" {
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {

		operation, err = serviceusageService.Operations.Get(operation.Name).Context(ctx).Do()
		if err != nil {
			panic(err.Error())
		}

		if operation.Done == true {
			ticker.Stop()
			break
		}
	}

}

func removeDefaultServices(projectID string, credJSON []byte, enabledServices []string) []string {

	projectNumber := getProjectNumber(projectID, credJSON)

	defaultServices := []string{
		"bigquery.googleapis.com",
		"bigquerystorage.googleapis.com",
		"cloudapis.googleapis.com",
		"clouddebugger.googleapis.com",
		"cloudtrace.googleapis.com",
		"datastore.googleapis.com",
		"logging.googleapis.com",
		"monitoring.googleapis.com",
		"servicemanagement.googleapis.com",
		"serviceusage.googleapis.com",
		"sql-component.googleapis.com",
		"storage-api.googleapis.com",
		"storage-component.googleapis.com",
		"storage.googleapis.com",
		"iam.googleapis.com",
		"iamcredentials.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"cloudasset.googleapis.com",
	}

	services := make(map[string]int)

	for _, s := range enabledServices {
		services[s]++
	}

	for _, s := range defaultServices {

		delete(services, fmt.Sprintf("projects/%d/services/%s", projectNumber, s))
	}

	var nonDefaultServices []string
	for k := range services {
		nonDefaultServices = append(nonDefaultServices, k)
	}

	return nonDefaultServices
}

func getProjectNumber(projectID string, credJSON []byte) int64 {

	ctx := context.Background()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		log.Fatal(err)
	}

	project, err := cloudresourcemanagerService.Projects.Get(projectID).Context(ctx).Do()
	if err != nil {
		panic(err.Error())
	}

	return project.ProjectNumber
}

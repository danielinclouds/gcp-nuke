package gcp

import (
	"context"
	"fmt"
	"time"

	"github.com/danielinclouds/gcp-nuke/config"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/serviceusage/v1"
)

func ListNonDefaultServices(cfg *config.Config) {

	// TODO: Handle multiple pages
	resp, err := cfg.ServiceusageService.Services.
		List(fmt.Sprintf("projects/%s", cfg.Project)).
		Filter("state:ENABLED").
		PageSize(200).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	var enabledServices []string
	for _, service := range resp.Services {
		enabledServices = append(enabledServices, service.Name)
	}

	enabledServices = removeDefaultServices(cfg, enabledServices)
	for _, service := range enabledServices {
		cfg.Log.Infof("API Service: %s", service)
	}

}

func DisableAllNonDefaultServices(cfg *config.Config) {

	// TODO: Handle multiple pages
	resp, err := cfg.ServiceusageService.Services.
		List(fmt.Sprintf("projects/%s", cfg.Project)).
		Filter("state:ENABLED").
		PageSize(200).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	var enabledServices []string
	for _, service := range resp.Services {
		enabledServices = append(enabledServices, service.Name)
	}

	enabledServices = removeDefaultServices(cfg, enabledServices)
	for _, service := range enabledServices {
		disableService(cfg, service)
	}

}

func disableService(cfg *config.Config, serviceName string) {

	cfg.Log.Debugf("Disable service: %s", serviceName)
	operation, err := cfg.ServiceusageService.Services.
		Disable(serviceName, &serviceusage.DisableServiceRequest{DisableDependentServices: true}).
		Context(context.Background()).
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

		operation, err = cfg.ServiceusageService.Operations.Get(operation.Name).Context(context.Background()).Do()
		if err != nil {
			panic(err.Error())
		}

		if operation.Done == true {
			ticker.Stop()
			break
		}
	}

}

func removeDefaultServices(cfg *config.Config, enabledServices []string) []string {

	projectNumber := getProjectNumber(cfg)

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

func getProjectNumber(cfg *config.Config) int64 {

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(context.Background(), option.WithCredentialsJSON(cfg.Credentials.JSON))
	if err != nil {
		cfg.Log.Fatal(err)
	}

	project, err := cloudresourcemanagerService.Projects.
		Get(cfg.Project).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	return project.ProjectNumber
}

func isServiceDisabled(cfg *config.Config, api string) bool {

	service, err := cfg.ServiceusageService.Services.
		Get(fmt.Sprintf("projects/%s/services/%s", cfg.Project, api)).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	if service.State == "ENABLED" {
		return false
	}

	return true
}

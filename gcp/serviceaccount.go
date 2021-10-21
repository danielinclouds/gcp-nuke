package gcp

import (
	"context"
	"fmt"
	"os"

	"github.com/danielinclouds/gcp-nuke/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListServiceAccounts(cfg *config.Config) {

	resp, err := cfg.IamService.Projects.ServiceAccounts.
		List(fmt.Sprintf("projects/%s", cfg.Project)).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	for _, sa := range resp.Accounts {

		if sa.Email == cfg.Credentials.Email {
			log.Infof("Skipping current %s service account", cfg.Credentials.Email)
			continue
		}

		log.Infof("Service account: %s", sa.Name)
	}

}

func DeleteAllServiceAccounts(cfg *config.Config) {

	resp, err := cfg.IamService.Projects.ServiceAccounts.
		List(fmt.Sprintf("projects/%s", cfg.Project)).
		Context(context.Background()).
		Do()
	if err != nil {
		panic(err.Error())
	}

	for _, sa := range resp.Accounts {

		if sa.Email == cfg.Credentials.Email {
			log.Debugf("Skipping current %s service account", cfg.Credentials.Email)
			continue
		}

		log.Debugf("Delete service account: %s", sa.Name)
		deleteServiceAccount(cfg, sa.Name)

	}

}

func deleteServiceAccount(cfg *config.Config, name string) {

	_, err := cfg.IamService.Projects.ServiceAccounts.Delete(name).Context(context.Background()).Do()
	if err != nil {
		panic(err.Error())
	}

}

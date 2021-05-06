package gcp

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListServiceAccounts(projectId string, credentials Credentials) {

	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsJSON(credentials.JSON))
	if err != nil {
		panic(err.Error())
	}

	resp, err := iamService.Projects.ServiceAccounts.
		List(fmt.Sprintf("projects/%s", projectId)).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	for _, sa := range resp.Accounts {

		if sa.Email == credentials.Email {
			log.Infof("Skipping current %s service account", credentials.Email)
			continue
		}

		log.Infof("Service account: %s", sa.Name)
	}

}

func DeleteAllServiceAccounts(projectId string, credentials Credentials) {

	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsJSON(credentials.JSON))
	if err != nil {
		panic(err.Error())
	}

	resp, err := iamService.Projects.ServiceAccounts.
		List(fmt.Sprintf("projects/%s", projectId)).
		Context(ctx).
		Do()
	if err != nil {
		panic(err.Error())
	}

	for _, sa := range resp.Accounts {

		if sa.Email == credentials.Email {
			log.Debugf("Skipping current %s service account", credentials.Email)
			continue
		}

		log.Debugf("Delete service account: %s", sa.Name)
		deleteServiceAccount(sa.Name, credentials)

	}

}

func deleteServiceAccount(name string, credentials Credentials) {

	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsJSON(credentials.JSON))
	if err != nil {
		panic(err.Error())
	}

	_, err = iamService.Projects.ServiceAccounts.Delete(name).Context(ctx).Do()
	if err != nil {
		panic(err.Error())
	}

}

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

func ListServiceAccounts(projectId string, credJSON []byte) {

	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsJSON(credJSON))
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
		log.Infof("Service account: %s", sa.Name)
	}

}

func DeleteAllServiceAccounts(projectId string, credJSON []byte) {

	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsJSON(credJSON))
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

		// TODO
		// 1. Don't skip gcp-nuke service account, skip current SA
		if sa.Email == fmt.Sprintf("gcp-nuke@%s.iam.gserviceaccount.com", projectId) {
			log.Debug("Skipping gcp-nuke service account")
			continue
		}

		log.Debugf("Deleting service account: %s", sa.Name)
		deleteServiceAccount(sa.Name, credJSON)

	}

}

func deleteServiceAccount(name string, credJSON []byte) {

	ctx := context.Background()
	iamService, err := iam.NewService(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	_, err = iamService.Projects.ServiceAccounts.Delete(name).Context(ctx).Do()
	if err != nil {
		panic(err.Error())
	}

}

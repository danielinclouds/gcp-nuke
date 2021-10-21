package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/danielinclouds/gcp-nuke/config"
	"github.com/danielinclouds/gcp-nuke/credentials"
	"github.com/danielinclouds/gcp-nuke/gcp"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/serviceusage/v1"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/urfave/cli/v2"
)

func Command() {

	app := &cli.App{
		Name:      "gcp-nuke",
		Usage:     "The GCP project cleanup tool",
		Version:   "v0.1.0",
		UsageText: "gcp-nuke [option]...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project",
				Aliases:  []string{"p"},
				Usage:    "GCP project id",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "credentials",
				Aliases:  []string{"c"},
				Usage:    "Path to file with GCP credentials",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "dry-run",
				Usage:    "Dry run without taking any action.",
				Value:    false,
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {

			creds, err := credentials.FindCredentials(c.String("credentials"))
			if err != nil {
				return err
			}

			// Logging
			var log = &logrus.Logger{
				Out:       os.Stdout,
				Formatter: new(logrus.TextFormatter),
				Hooks:     make(logrus.LevelHooks),
				Level:     logrus.DebugLevel,
			}

			// Bucket
			storageClient, err := storage.NewClient(context.Background(), option.WithCredentialsJSON(creds.JSON))
			if err != nil {
				panic(err.Error())
			}

			// GKE clusters
			containerService, err := container.NewService(context.Background(), option.WithCredentialsJSON(creds.JSON))
			if err != nil {
				panic(err.Error())
			}

			// Service usage
			serviceusageService, err := serviceusage.NewService(context.Background(), option.WithCredentialsJSON(creds.JSON))
			if err != nil {
				panic(err.Error())
			}

			// Networks
			computeService, err := compute.NewService(context.Background(), option.WithCredentialsJSON(creds.JSON))
			if err != nil {
				panic(err.Error())
			}

			// PubSub
			pubSubClient, err := pubsub.NewClient(context.Background(), c.String("project"), option.WithCredentialsJSON(creds.JSON))
			if err != nil {
				panic(err.Error())
			}

			// Service Account
			iamService, err := iam.NewService(context.Background(), option.WithCredentialsJSON(creds.JSON))
			if err != nil {
				panic(err.Error())
			}

			// Assets
			assetsClient, err := asset.NewClient(context.Background(), option.WithCredentialsJSON(creds.JSON))
			if err != nil {
				panic(err.Error())
			}

			// Close clients
			defer storageClient.Close()
			defer pubSubClient.Close()
			defer assetsClient.Close()

			// Config
			cfg := config.Config{
				Project:             c.String("project"),
				Credentials:         creds,
				Log:                 log,
				StorageClient:       storageClient,
				ContainerService:    containerService,
				ServiceusageService: serviceusageService,
				ComputeService:      computeService,
				PubSubClient:        pubSubClient,
				IamService:          iamService,
				AssetsClient:        assetsClient,
			}

			gcp.ListPubSub(&cfg)
			gcp.ListGKEClusters(&cfg)
			gcp.ListBuckets(&cfg)
			gcp.ListVPC(&cfg)
			gcp.ListServiceAccounts(&cfg)
			gcp.ListNonDefaultServices(&cfg)
			// gcp.ListAssets(&cfg)

			if c.Bool("dry-run") {
				return nil
			}

			gcp.DeleteAllGKEClusters(&cfg)
			gcp.DeleteAllPubSub(&cfg)
			gcp.DeleteAllBuckets(&cfg)
			gcp.DeleteAllVPC(&cfg)
			// gcp.DeleteAllServiceAccounts(&cfg)
			// gcp.DisableAllNonDefaultServices(&cfg)

			return nil
		},
	}

	cli.AppHelpTemplate = fmt.Sprintf(`%s
ENVIRONMENT VARIABLES:
   GOOGLE_CREDENTIALS
   GOOGLE_CLOUD_KEYFILE_JSON

EXAMPLES:
   # Delete all resources from project using credentials from env
   export GOOGLE_CREDENTIALS=$(cat gcp-nuke.json)
   gcp-nuke --project PROJECT_ID
   
   # Delete all resources from project using credentials file
   gcp-nuke --project PROJECT_ID --credentials creds.json
   
   # List resources in project
   gcp-nuke --project PROJECT_ID --dry-run
	
	`, cli.AppHelpTemplate)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

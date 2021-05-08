package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/danielinclouds/gcp-nuke/gcp"

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

			credentials, err := gcp.FindCredentials(c.String("credentials"))
			if err != nil {
				panic(err.Error())
			}

			// gcp.ListPubSub(c.String("project"), credentials.JSON)
			gcp.ListBuckets(c.String("project"), credentials.JSON)
			// gcp.ListGKEClusters(c.String("project"), credentials.JSON)
			// gcp.ListVPC(c.String("project"), credentials.JSON)
			// gcp.ListServiceAccounts(c.String("project"), credentials)
			// gcp.ListNonDefaultServices(c.String("project"), credentials.JSON)
			// gcp.ListAssets(c.String("project"), credentials.JSON)

			if c.Bool("dry-run") == true {
				return nil
			}

			// gcp.DeleteAllGKEClusters(c.String("project"), credentials.JSON)
			// gcp.DeleteAllPubSub(c.String("project"), credentials.JSON)
			gcp.DeleteAllBuckets(c.String("project"), credentials.JSON)
			// gcp.DeleteAllVPC(c.String("project"), credentials.JSON)
			// gcp.DeleteAllServiceAccounts(c.String("project"), credentials)
			// gcp.DisableAllNonDefaultServices(c.String("project"), credentials.JSON)

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

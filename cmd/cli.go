package cmd

import (
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
		UsageText: "e.g. gcp-nuke --project test-nuke-262510",
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
				Usage:    "File with GCP credentials",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {

			credentials, err := gcp.FindCredentials(c.String("credentials"))
			if err != nil {
				panic(err.Error())
			}

			gcp.ListPubSub(c.String("project"), credentials.JSON)
			gcp.ListBuckets(c.String("project"), credentials.JSON)
			gcp.ListGKEClusters(c.String("project"), credentials.JSON)
			gcp.ListVPC(c.String("project"), credentials.JSON)
			gcp.ListServiceAccounts(c.String("project"), credentials.JSON)
			gcp.ListNonDefaultServices(c.String("project"), credentials.JSON)
			gcp.ListAssets(c.String("project"), credentials.JSON)

			// ===================================================================

			gcp.DeleteAllGKEClusters(c.String("project"), credentials.JSON)
			gcp.DeleteAllPubSub(c.String("project"), credentials.JSON)
			gcp.DeleteAllBuckets(c.String("project"), credentials.JSON)
			gcp.DeleteAllVPC(c.String("project"), credentials.JSON)
			// gcp.DeleteAllServiceAccounts(c.String("project"), credentials.JSON)
			// gcp.DisableAllNonDefaultServices(c.String("project"), credentials.JSON)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

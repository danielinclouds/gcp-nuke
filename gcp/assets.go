package gcp

import (
	"context"
	"fmt"

	"github.com/danielinclouds/gcp-nuke/config"
	"google.golang.org/api/iterator"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

func ListAssets(cfg *config.Config) {

	if isServiceDisabled(cfg, "cloudasset.googleapis.com") {
		cfg.Log.Debug("Assets API is disabled")
		return
	}

	req := &assetpb.SearchAllResourcesRequest{
		Scope: fmt.Sprintf("projects/%s", cfg.Project),
	}

	it := cfg.AssetsClient.SearchAllResources(context.Background(), req)
	fmt.Println("Remaining assets:")
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("%s %s\n", resp.AssetType, resp.Name)
	}
}

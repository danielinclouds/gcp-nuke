package gcp

import (
	"context"
	"fmt"

	asset "cloud.google.com/go/asset/apiv1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

func ListAssets(projectID string, credJSON []byte) {
	ctx := context.Background()
	client, err := asset.NewClient(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	req := &assetpb.SearchAllResourcesRequest{
		Scope: fmt.Sprintf("projects/%s", projectID),
	}

	it := client.SearchAllResources(ctx, req)
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

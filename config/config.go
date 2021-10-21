package config

import (
	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/danielinclouds/gcp-nuke/credentials"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/serviceusage/v1"
)

type Config struct {
	Project             string
	Credentials         credentials.Credentials
	StorageClient       *storage.Client
	ContainerService    *container.Service
	ServiceusageService *serviceusage.Service
	ComputeService      *compute.Service
	PubSubClient        *pubsub.Client
	IamService          *iam.Service
	AssetsClient        *asset.Client
}

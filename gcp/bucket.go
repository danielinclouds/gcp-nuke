package gcp

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/danielinclouds/gcp-nuke/config"
)

func ListBuckets(cfg *config.Config) {

	it := cfg.StorageClient.Buckets(context.Background(), cfg.Project)

	for {
		bucket, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		cfg.Log.Infof("Bucket: %s", bucket.Name)
	}

}

func DeleteAllBuckets(cfg *config.Config) {

	it := cfg.StorageClient.Buckets(context.Background(), cfg.Project)

	for {
		bucket, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		deleteBucket(cfg, bucket.Name)
	}
}

func deleteBucket(cfg *config.Config, bucketName string) {

	cfg.Log.Debugf("Delete bucket: %s", bucketName)
	disableBucketVersioning(cfg, bucketName)
	emptyBucket(cfg, bucketName)

	err := cfg.StorageClient.Bucket(bucketName).Delete(context.Background())
	if err != nil {
		panic(err.Error())
	}

}

func disableBucketVersioning(cfg *config.Config, bucketName string) {

	cfg.Log.Debugf("Disable bucket versioning: %s", bucketName)
	_, err := cfg.StorageClient.Bucket(bucketName).Update(context.Background(), storage.BucketAttrsToUpdate{
		VersioningEnabled: false,
	})
	if err != nil {
		panic(err.Error())
	}

}

func emptyBucket(cfg *config.Config, bucketName string) {

	bucket := cfg.StorageClient.Bucket(bucketName)

	query := &storage.Query{
		Prefix:   "",
		Versions: true,
	}

	it := bucket.Objects(context.Background(), query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			cfg.Log.Fatal(err)
		}

		object := bucket.Object(attrs.Name).Generation(attrs.Generation)
		cfg.Log.Debugf("Delete object: %s generation: %d", attrs.Name, attrs.Generation)
		err = object.Delete(context.Background())
		if err != nil {
			panic(err.Error())
		}
	}

}

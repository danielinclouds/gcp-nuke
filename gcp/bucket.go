package gcp

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListBuckets(projectID string, credJSON []byte) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	it := client.Buckets(ctx, projectID)

	for {
		bucket, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		log.Infof("Bucket: %s", bucket.Name)
	}

}

func DeleteAllBuckets(projectID string, credJSON []byte) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	it := client.Buckets(ctx, projectID)

	for {
		bucket, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		deleteBucket(bucket.Name, credJSON)
	}
}

func deleteBucket(bucketName string, credJSON []byte) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	log.Debugf("Delete bucket: %s", bucketName)
	disableBucketVersioning(bucketName, credJSON)
	emptyBucket(bucketName, credJSON)

	err = client.Bucket(bucketName).Delete(ctx)
	if err != nil {
		panic(err.Error())
	}

}

func disableBucketVersioning(bucketName string, credJSON []byte) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	log.Debugf("Disable bucket versioning: %s", bucketName)
	_, err = client.Bucket(bucketName).Update(ctx, storage.BucketAttrsToUpdate{
		VersioningEnabled: false,
	})
	if err != nil {
		panic(err.Error())
	}

}

func emptyBucket(bucketName string, credJSON []byte) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	bucket := client.Bucket(bucketName)

	query := &storage.Query{
		Prefix:   "",
		Versions: true,
	}

	log.Debugf("Delete bucket objects: %s", bucketName)
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		object := bucket.Object(attrs.Name).Generation(attrs.Generation)
		err = object.Delete(ctx)
		if err != nil {
			panic(err.Error())
		}
	}

}

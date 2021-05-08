package gcp

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListPubSub(projectId string, credJSON []byte) {
	if isServiceDisabled(projectId, credJSON, "pubsub.googleapis.com") {
		log.Debug("PubSub API is disabled")
		return
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	defer client.Close()

	subscriptionIterator := client.Subscriptions(ctx)

	for {
		subscription, err := subscriptionIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		log.Infof("Subscription: %s", subscription.ID())

	}

	topicIterator := client.Topics(ctx)

	for {
		topic, err := topicIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		log.Infof("Topic: %s", topic.ID())
	}
}

func DeleteAllPubSub(projectId string, credJSON []byte) {
	if isServiceDisabled(projectId, credJSON, "pubsub.googleapis.com") {
		log.Debug("PubSub API is disabled")
		return
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId, option.WithCredentialsJSON(credJSON))
	if err != nil {
		panic(err.Error())
	}

	defer client.Close()

	subscriptionIterator := client.Subscriptions(ctx)

	for {
		subscription, err := subscriptionIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		log.Debugf("Deleting subscription: %s", subscription.ID())
		err = subscription.Delete(ctx)
		if err != nil {
			panic(err.Error())
		}

	}

	topicIterator := client.Topics(ctx)

	for {
		topic, err := topicIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		log.Debugf("Deleting topic: %s", topic.ID())
		err = topic.Delete(ctx)
		if err != nil {
			panic(err.Error())
		}
	}
}

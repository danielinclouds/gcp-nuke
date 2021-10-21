package gcp

import (
	"context"
	"os"

	"github.com/danielinclouds/gcp-nuke/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func ListPubSub(cfg *config.Config) {
	if isServiceDisabled(cfg, "pubsub.googleapis.com") {
		log.Debug("PubSub API is disabled")
		return
	}

	subscriptionIterator := cfg.PubSubClient.Subscriptions(context.Background())

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

	topicIterator := cfg.PubSubClient.Topics(context.Background())

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

func DeleteAllPubSub(cfg *config.Config) {
	if isServiceDisabled(cfg, "pubsub.googleapis.com") {
		log.Debug("PubSub API is disabled")
		return
	}

	subscriptionIterator := cfg.PubSubClient.Subscriptions(context.Background())

	for {
		subscription, err := subscriptionIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		log.Debugf("Deleting subscription: %s", subscription.ID())
		err = subscription.Delete(context.Background())
		if err != nil {
			panic(err.Error())
		}

	}

	topicIterator := cfg.PubSubClient.Topics(context.Background())

	for {
		topic, err := topicIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		log.Debugf("Deleting topic: %s", topic.ID())
		err = topic.Delete(context.Background())
		if err != nil {
			panic(err.Error())
		}
	}
}

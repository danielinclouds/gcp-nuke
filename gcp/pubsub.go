package gcp

import (
	"context"

	"github.com/danielinclouds/gcp-nuke/config"
	"google.golang.org/api/iterator"
)

func ListPubSub(cfg *config.Config) {
	if isServiceDisabled(cfg, "pubsub.googleapis.com") {
		cfg.Log.Debug("PubSub API is disabled")
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

		cfg.Log.Infof("Subscription: %s", subscription.ID())

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

		cfg.Log.Infof("Topic: %s", topic.ID())
	}
}

func DeleteAllPubSub(cfg *config.Config) {
	if isServiceDisabled(cfg, "pubsub.googleapis.com") {
		cfg.Log.Debug("PubSub API is disabled")
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

		cfg.Log.Debugf("Deleting subscription: %s", subscription.ID())
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

		cfg.Log.Debugf("Deleting topic: %s", topic.ID())
		err = topic.Delete(context.Background())
		if err != nil {
			panic(err.Error())
		}
	}
}

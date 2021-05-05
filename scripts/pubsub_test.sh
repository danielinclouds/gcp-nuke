#!/bin/zsh

gcloud pubsub topics create topic-1
gcloud pubsub subscriptions create sub-1 --topic=topic-1

gcloud pubsub topics create topic-2
gcloud pubsub subscriptions create sub-2 --topic=topic-2
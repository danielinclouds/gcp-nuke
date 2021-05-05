#!/bin/zsh

gcloud iam service-accounts create gke-default-node --display-name=gke-default-node
gcloud projects add-iam-policy-binding test-223613 \
    --member='serviceAccount:gke-default-node@test-223613.iam.gserviceaccount.com' \
    --role='roles/editor'


gcloud container clusters create test-1 \
    --zone europe-west1-b \
    --machine-type=e2-standard-2 \
    --service-account="gke-default-node@test-223613.iam.gserviceaccount.com" \
    --num-nodes 1

gcloud container clusters create test-2 \
    --region europe-west1 \
    --machine-type=e2-standard-2 \
    --service-account="gke-default-node@test-223613.iam.gserviceaccount.com" \
    --num-nodes 1


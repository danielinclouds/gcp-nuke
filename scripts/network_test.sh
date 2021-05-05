#!/bin/zsh

gcloud compute networks create daniel --subnet-mode=custom
gcloud compute networks subnets create s1 --network=daniel --range="10.0.0.0/16" --region=us-central1
gcloud compute networks subnets create s2 --network=daniel --range="10.1.0.0/16" --region=europe-west1


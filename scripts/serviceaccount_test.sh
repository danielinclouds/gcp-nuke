#!/bin/zsh

gcloud iam service-accounts list

gcloud iam service-accounts create test-1 --display-name=test-1
gcloud iam service-accounts create test-2 --display-name=test-2
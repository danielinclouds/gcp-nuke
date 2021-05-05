#!/bin/zsh

# Create bucket
gsutil mb gs://danielinclouds-test-123
gsutil mb gs://danielinclouds-test-456

# create versioning
gsutil versioning set on gs://danielinclouds-test-123

# add file v1
gsutil cp ../files/file_1 gs://danielinclouds-test-123/file

# add file v2
gsutil cp ../files/file_2 gs://danielinclouds-test-123/file

# add file to folder
gsutil cp ../files/file_1 gs://danielinclouds-test-123/folder/file


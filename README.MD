# gcp-nuke



## Usage
Authenticate with env variables:
```
GOOGLE_CREDENTIALS=$(cat credentials.json)
export GOOGLE_CREDENTIALS

go run main.go --project test-123
```

Authenticate with credentials file:
```
go run main.go --project test-123 --credentials credentials.json
```


## Supported resources
Currently gcp-nuke deletes following resources:
- GCS buckets
- GKE clusters
- Networks
- Subnetworks
- PubSub
- Service Accounts
- API Services


## Release
```sh
git tag -a v0.1.0 -m "Some release"
git push --tags
goreleaser release --rm-dist
```


## GitHub Actions
Add workflow file to your GCP project repository to schedule cleanup.
```yaml
name: Nuke GCP project

on:
  schedule:
    - cron: '30 17 * * *'  # Cron runs in UTC0 timezone

env:
  PROJECT_ID:         ${{ secrets.PROJECT_ID }}
  GOOGLE_CREDENTIALS: ${{ secrets.GOOGLE_CREDENTIALS }}


jobs:
  nuke:
    runs-on: ubuntu-latest
    steps:
    
    - name: Install gcp-nuke
      run: |
        wget https://github.com/danielinclouds/gcp-nuke/releases/download/v0.4.0/gcp-nuke_0.4.0_Linux_x86_64.tar.gz
        tar zxvf gcp-nuke_0.4.0_Linux_x86_64.tar.gz
        chmod +x gcp-nuke
        sudo mv gcp-nuke /usr/local/bin
    
    - name: Nuke GCP Project
      run: |
        gcp-nuke --project "$PROJECT_ID"
```
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
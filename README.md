# gcshare
Simple CLI tool to quickly share a local file with an email address via Google Cloud Storage SignedURL

## Requirements
- GCP account with Cloud Storage bucket already created
- GCP [service account](https://cloud.google.com/iam/docs/service-account-overview) with minimum of storage.objects.get permission for bucket
- [gcloud](https://cloud.google.com/sdk/docs/install) installed locally and [service account credentials available](https://cloud.google.com/sdk/gcloud/reference/auth/activate-service-account)
- Email with credentials available to use for sending (example: gmail with app password)

## Usage
1. `go build` to compile
2. `cp example_config.yml config.yml` then edit the yml with the appropriate fields
3. Syntax: `gcshare <filepath> <recipient_email>` will upload a local file to cloud storage, generate signed URL, and send URL to recipient for access

Note: Email password currently stored in plaintext config.yml. Use at your own discretion. Recommend if using this tool having a separate email designated just for sending signedURLs. 

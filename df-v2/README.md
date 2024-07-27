# Dialogflow using cloud functions v2

Deplpoying Dialogflow v2 using cloud functions v2.
```shell
# --allow-unauthenticated
gcloud functions deploy go-http-function --gen2 --runtime go122 --trigger-http \
--entry-point HelloShirts --source . --project=$GCP_PROJECT --region=us-central1

# Get details
gcloud functions describe go-http-function --gen2 --region=us-central1 --project=$GCP_PROJECT

# Get function URL:
FUNCTION_URL=$(gcloud functions describe go-http-function --gen2 --region=us-central1 --project=$GCP_PROJECT  --format="value(serviceConfig.uri)")

# Send test data
curl -m 70 -X POST https://go-http-function-2dbml6flea-uc.a.run.app -H "Authorization:bearer $(gcloud auth print-identity-token)" -H "Content-Type:application/json" -d '{"name": "Hello World"}'
gcloud functions call go-http-function --gen2 --region us-central1 --project=$GCP_PROJECT  --data='{"name": "Jack"}'
```

Other examples:

```shell
# enable a func
gcloud  functions deploy cx-webhook --gen2 --runtime go122 --trigger-http --entry-point HandleWebhookRequest --source . --project=$GCP_PROJECT --region=us-central1

# Echo fun
gcloud beta functions deploy log-test-v2 --gen2 --runtime go122 --trigger-http --entry-point  ManojS --source . --project=$GCP_PROJECT  --region=us-central1
gcloud beta functions deploy log-test-v3 --runtime go122 --trigger-http --entry-point  ManojS --source . --project=$GCP_PROJECT  --region=us-central1

# Get details
gcloud beta functions describe cx-webhook --gen2 --region=us-central1

# Get function URL:
CX_FUNCTION_URL=$(gcloud alpha functions describe cx-webhook --gen2 --region=us-central1 --format="value(serviceConfig.uri)")

# Send test data
curl -m 70 -X POST https://cx-webhook-2yi7hjkwba-uc.a.run.app -H "Authorization:bearer $(gcloud auth print-identity-token)" -H "Content-Type:application/json" -d '{"name": "Hello World"}'
gcloud alpha functions call cx-webhook --gen2 --region us-central1
```

### Testing on GCP functions

This is testing of the GCP cloud functions v2 using my personal account.

```shell
# Deploy unauthenticated
gcloud beta functions deploy go-all --gen2 --runtime go122 --trigger-http --entry-point HelloShirts --source . --project=$GCP_PROJECT --region=us-central1 --allow-unauthenticated
gcloud beta functions describe go-all --gen2 --region=us-central1
curl https://go-all-2dbml6flea-uc.a.run.app
curl -m 70 -X POST https://go-all-2dbml6flea-uc.a.run.app -H "Authorization:bearer $(gcloud auth print-identity-token)" -H "Content-Type:application/json" -d '{"name": "Hello World"}'

# Deploy authenticated
gcloud beta functions deploy go-auth --gen2 --runtime go122 --trigger-http --entry-point HelloShirts --source . --project=$GCP_PROJECT --region=us-central1
curl -m 70 -X POST https://go-all-2dbml6flea-uc.a.run.app -H "Authorization:bearer $(gcloud auth print-identity-token)" -H "Content-Type:application/json" -d '{"name": "Hello World"}'
```
Dialogflow CX
```shell
gcloud beta functions deploy cx-webhook --gen2 --runtime go122 --trigger-http --entry-point HandleWebhookRequest --source . --project=$GCP_PROJECT --region=us-central1
curl -m 70 -X POST https://cx-webhook-2dbml6flea-uc.a.run.app -H "Authorization:bearer $(gcloud auth print-identity-token)" -H "Content-Type:application/json" -d '{ "FulfillmentInfo": { "tag": "welcome" }}'
```
## Notes

- [Getting Started Tutorial: Cloud Functions (2nd gen)](https://cloud.google.com/functions/docs/2nd-gen/getting-started)
-
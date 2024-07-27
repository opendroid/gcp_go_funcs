# Dead Letter Topic
These functions test the Pub/Sub retries and dead-letter topics.

## Test Cloud Functions

First deploy the pub-sub-triggered and HTTP function using --gen2, it uses cloud run to deploy functions.
The configure the topics and subscriptions. Enable  "Push authentication" with audience URL as your cloud-run URL. Also 
make sure to provide the "Cloud Run Invoker" and "Cloud Function Invoker" permission to cloud-run service account.

In this example, the cloud functions `AckPubMessage` and `BadAckFunc` tests the subscription retries and dead-letter topics.
Main Pub/Sub topic __radio-pluto__, has  cloud http trigger func `BadAckFunc` is attached to it. This always returns code 500.
Since __radio-pluto__ topic is set to retry 5 times, the message will be sent 5 times to `BadAckFunc`. 
Once all retries fails, the message will be pushed to __pluto-dead-letter__ topic.

Dead letter topic: __pluto-dead-letter__ has listener `AckPubMessage`. Once the messages is received this function 
consumes it. Once consumed message is discarded (if configured).

```shell
# Activate  right project configuration
gcloud config configurations activate gcp-experiments
# Note that --runtime go122  is supported 
# https://cloud.google.com/functions/docs/runtime-support#go
# Ensure that the Cloud functions has permission "Cloud Run Invoker" and "Cloud Function Invoker"
gcloud functions deploy AckPubMessage  --gen2 --runtime go122 --trigger-topic pluto-dead-letter --project=gcp-experiments-334602
# Deploy a http func so we can return a error code (to test retry and dead-letter)
gcloud functions deploy BadAckFunc  --gen2 --runtime go122 --trigger-http --project=gcp-experiments-334602
```

### Testing Message Retries
Pushing messages using command:
```shell
# Push a message to a 'radio-pluto' topic
gcloud pubsub topics publish radio-pluto --message='{"name": "GCP", "rating": "12-star"}' --project=gcp-experiments-334602
```
Then observe the logs for retries. 

## Multiple Subscriptions on a topic.
This is not tested yet.
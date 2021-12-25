# Dead Letter Topic
These functions test the Pub/Sub retries and dead-letter topics.

## Test Cloud Functions
The cloud functions `AckPubMessage` and `BadAckFunc` test the subscription retries and dead-letter topics.

Main Pub/Sub topic __radio-pluto__, has  cloud http trigger func `BadAckFunc` is attached to it that always returns code 500.
Since __radio-pluto__ topic is set to retry 5 times, the message will be sent 5 times to `BadAckFunc`. 
Once all retries fails, the message will be pushed to __pluto-dead-letter__ topic.

Dead letter topic: __pluto-dead-letter__ has listener `AckPubMessage`. Once the messages is received this function 
consumes it. Once consumed message is discarded (if configured).

```shell
# Activate  right project configuration
gcloud config configurations activate gcp-experiments
# Note that --runtime go116  is supported as beta
# https://cloud.google.com/functions/docs/concepts/go-runtime
gcloud functions deploy AckPubMessage --runtime go116 --trigger-topic pluto-dead-letter --project=gcp-experiments-334602
# Deploy a http func so we can return a error code (to test retry and dead-letter)
gcloud functions deploy BadAckFunc --runtime go116 --trigger-http --allow-unauthenticated --project=gcp-experiments-334602 --region=us-west1
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
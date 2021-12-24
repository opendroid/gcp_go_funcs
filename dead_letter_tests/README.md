# Dead Letter Tests

## Dead Letter Topic
These pair of cloud functions test the subscription retries and dead-letter topic.
Main topic: __radio-pluto__, a cloud http trigger func 'BadAckFunc' is attached to it that always returns code 500.
The __radio-pluto__ topic is set to retry 5 times and after all retries fail, will push the message to 
__pluto-dead-letter__ topic.

Dead letter topic: __pluto-dead-letter__ has a listener 'AckPubMessage'. Once the messages is received this function 
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

## Multiple Subscriptions on a topic.


## Testing
Pushing messages:
```shell
# Push a message to a 'radio-pluto' topic
gcloud pubsub topics publish radio-pluto --message='{"name": "GCP", "rating": "12-star"}' --project=gcp-experiments-334602
```

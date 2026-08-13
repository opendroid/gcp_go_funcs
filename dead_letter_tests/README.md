# Dead Letter Topic

These functions test Google Cloud Pub/Sub retry behavior and dead-letter topics using **2nd-generation Cloud Functions**, which run on Cloud Run.

The test flow is:

```text
radio-pluto
    │
    ▼
subscription
    │
    ▼
BadAckFunc
    │
    │ HTTP 500
    │ retry up to configured max delivery attempts
    ▼
pluto-dead-letter
    │
    ▼
dead-letter-reader
    │
    ▼
AckPubMessage
    │
    ▼
HTTP 200 / ACK
```

## Test Functions

Two functions are used:

* `BadAckFunc`

  * HTTP-triggered function.
  * Always returns HTTP `500`.
  * Used to force Pub/Sub message redelivery.

* `AckPubMessage`

  * HTTP-triggered function used by the dead-letter subscription.
  * Successfully consumes the message and returns HTTP `200`.

The main Pub/Sub topic is:

```text
radio-pluto
```

Its push subscription invokes `BadAckFunc`.

The subscription is configured with a dead-letter policy and a maximum number of delivery attempts. After the configured delivery attempts fail, Pub/Sub forwards the message to:

```text
pluto-dead-letter
```

The dead-letter topic has the push subscription:

```text
dead-letter-reader
```

which invokes `AckPubMessage`.

> Pub/Sub calls this setting `max-delivery-attempts`. For example, `5` means approximately five total delivery attempts, not the initial attempt plus five additional retries.

---

## Deploy the Cloud Functions

Activate the correct gcloud configuration:

```shell
gcloud config configurations activate gcp-experiments
```

Go 1.26 is supported:

```text
https://cloud.google.com/functions/docs/runtime-support#go
```

Deploy AckPubMessage as a Pub/Sub-triggered CloudEvent function. The function is invoked when a message is published to the pluto-dead-letter topic:

```shell
gcloud functions deploy AckPubMessage \
  --gen2 \
  --runtime=go126 \
  --region=us-central1 \
  --project=$YOUR_PROJECT_ID
```

Deploy `BadAckFunc` as an HTTP function so it can explicitly return HTTP `500`:

```shell
gcloud functions deploy BadAckFunc \
  --gen2 \
  --runtime=go126 \
  --region=us-central1 \
  --trigger-http \
  --project=$YOUR_PROJECT_ID
```

Because these are Gen2 Cloud Functions, the functions are backed by Cloud Run services.

For example, `AckPubMessage` may have URLs such as:

```text
Cloud Functions URL:
https://us-central1-$YOUR_PROJECT_ID.cloudfunctions.net/AckPubMessage

Cloud Run URLs:
https://ackpubmessage-2dbml6flea-uc.a.run.app
https://ackpubmessage-$YOUR_PROJECT_NUMBER.us-central1.run.app
```

For the Pub/Sub push subscription in this test, use the **Cloud Functions URL**.

---

## Configure Push Authentication

The Pub/Sub push subscription authenticates using an OIDC ID token.

The calling service account used in this example is:

```text
$YOUR_PROJECT_NUMBER-compute@developer.gserviceaccount.com
```

The Pub/Sub service agent is:

```text
service-$YOUR_PROJECT_NUMBER@gcp-sa-pubsub.iam.gserviceaccount.com
```

### 1. Allow Pub/Sub to generate an OIDC token

Grant `Service Account Token Creator` to the Pub/Sub service agent on the calling service account:

```shell
gcloud iam service-accounts add-iam-policy-binding \
  $YOUR_PROJECT_NUMBER-compute@developer.gserviceaccount.com \
  --member="serviceAccount:service-$YOUR_PROJECT_NUMBER@gcp-sa-pubsub.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountTokenCreator" \
  --project=$YOUR_PROJECT_ID
```

### 2. Allow the calling service account to invoke the function

Grant the calling service account `Cloud Run Invoker` on the underlying Gen2 Cloud Run service:

```shell
gcloud run services add-iam-policy-binding ackpubmessage \
  --member="serviceAccount:$YOUR_PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role="roles/run.invoker" \
  --region=us-central1 \
  --project=$YOUR_PROJECT_ID
```

If required by the Cloud Functions configuration, also verify the corresponding Cloud Functions invocation permissions.

---

## Configure the Dead-Letter Push Subscription

The dead-letter topic is:

```text
projects/$YOUR_PROJECT_ID/topics/pluto-dead-letter
```

The subscription consuming it is:

```text
projects/$YOUR_PROJECT_ID/subscriptions/dead-letter-reader
```

### Topic path in the push URL

The Go Functions Framework attempts to determine the Pub/Sub topic from the HTTP request path.

If Pub/Sub invokes only:

```text
https://us-central1-$YOUR_PROJECT_ID.cloudfunctions.net/AckPubMessage
```

the function receives the request at `/` and may log:

```text
WARNING: failed to extract Pub/Sub topic name from the URL request path: "/",
configure your subscription's push endpoint to use the following path pattern:
'projects/PROJECT_NAME/topics/TOPIC_NAME'
```

Therefore the push endpoint should include the Pub/Sub topic path:

```text
https://us-central1-$YOUR_PROJECT_ID.cloudfunctions.net/AckPubMessage/projects/$YOUR_PROJECT_ID/topics/pluto-dead-letter
```

### Important: configure the OIDC audience separately

The OIDC audience **must remain the base function URL**:

```text
https://us-central1-$YOUR_PROJECT_ID.cloudfunctions.net/AckPubMessage
```

Do **not** use the topic-qualified push URL as the OIDC audience.

Otherwise Pub/Sub can generate an ID token whose `aud` claim contains:

```text
/AckPubMessage/projects/$YOUR_PROJECT_ID/topics/pluto-dead-letter
```

and Cloud Run may reject it with:

```text
401
The request was not authorized to invoke this service.
The access token could not be verified.
```

Configure the subscription with:

```shell
gcloud pubsub subscriptions modify-push-config dead-letter-reader \
  --project=$YOUR_PROJECT_ID \
  --push-endpoint="https://us-central1-$YOUR_PROJECT_ID.cloudfunctions.net/AckPubMessage/projects/$YOUR_PROJECT_ID/topics/pluto-dead-letter" \
  --push-auth-service-account="$YOUR_PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --push-auth-token-audience="https://us-central1-$YOUR_PROJECT_ID.cloudfunctions.net/AckPubMessage"
```

The distinction is intentional:

```text
HTTP Push Endpoint
https://...cloudfunctions.net/AckPubMessage/projects/.../topics/pluto-dead-letter
                                         │
                                         └── lets the Functions Framework
                                             determine the source topic


OIDC Audience
https://...cloudfunctions.net/AckPubMessage
                                         │
                                         └── identifies the protected
                                             function/service
```

Verify the configuration:

```shell
gcloud pubsub subscriptions describe dead-letter-reader \
  --project=$YOUR_PROJECT_ID \
  --format='yaml(topic,pushConfig,deadLetterPolicy,retryPolicy)'
```

Expected push configuration:

```yaml
pushConfig:
  oidcToken:
    audience: https://us-central1-$YOUR_PROJECT_ID.cloudfunctions.net/AckPubMessage
    serviceAccountEmail: $YOUR_PROJECT_NUMBER-compute@developer.gserviceaccount.com
  pushEndpoint: https://us-central1-$YOUR_PROJECT_ID.cloudfunctions.net/AckPubMessage/projects/$YOUR_PROJECT_ID/topics/pluto-dead-letter
topic: projects/$YOUR_PROJECT_ID/topics/pluto-dead-letter
```

---

## Configure Retries and Dead-Letter Policy

The dead-letter configuration belongs to the **subscription consuming `radio-pluto`**, not to the topic itself.

Conceptually:

```text
radio-pluto
    │
    ▼
source subscription
    │
    ├── delivery attempt 1 → BadAckFunc → 500
    ├── delivery attempt 2 → BadAckFunc → 500
    ├── delivery attempt 3 → BadAckFunc → 500
    ├── delivery attempt 4 → BadAckFunc → 500
    └── delivery attempt 5 → BadAckFunc → 500
                                      │
                                      ▼
                             pluto-dead-letter
```

Configure the source subscription with a dead-letter topic and maximum delivery attempts, for example:

```shell
gcloud pubsub subscriptions update SOURCE_SUBSCRIPTION \
  --project=$YOUR_PROJECT_ID \
  --dead-letter-topic=pluto-dead-letter \
  --max-delivery-attempts=5
```

A retry policy can also be configured if delayed/exponential redelivery is desired.

---

## Testing the Dead-Letter Reader

Before testing the entire retry chain, verify that the dead-letter topic can successfully invoke `AckPubMessage`.

Publish directly to the dead-letter topic:

```shell
gcloud pubsub topics publish pluto-dead-letter \
  --project=$YOUR_PROJECT_ID \
  --message='{"test":"dead-letter-reader"}'
```

Inspect the `AckPubMessage` Cloud Run logs:

```shell
gcloud logging read \
  'resource.type="cloud_run_revision"
   AND resource.labels.service_name="ackpubmessage"' \
  --project=$YOUR_PROJECT_ID \
  --limit=20 \
  --freshness=5m \
  --format='table(timestamp,httpRequest.status,textPayload)'
```

Successful delivery should show:

```text
STATUS
200
```

HTTP `200` means Pub/Sub considers the pushed message acknowledged.

You can also verify that the earlier topic-path warning has disappeared:

```shell
gcloud logging read \
  'resource.type="cloud_run_revision"
   AND resource.labels.service_name="ackpubmessage"
   AND textPayload:"failed to extract Pub/Sub topic name"' \
  --project=$YOUR_PROJECT_ID \
  --freshness=5m \
  --limit=20
```

For newly delivered messages, this should return no matching warning.

---

## Testing Message Retries

Publish a message to the main topic:

```shell
gcloud pubsub topics publish radio-pluto \
  --project=$YOUR_PROJECT_ID \
  --message='{"name":"GCP","rating":"12-star"}'
```

The expected sequence is:

```text
1. Message published to radio-pluto
2. Pub/Sub pushes message to BadAckFunc
3. BadAckFunc returns HTTP 500
4. Pub/Sub redelivers according to the subscription retry policy
5. After approximately max-delivery-attempts failures,
   Pub/Sub forwards the message to pluto-dead-letter
6. dead-letter-reader pushes the message to AckPubMessage
7. AckPubMessage returns HTTP 200
8. Pub/Sub acknowledges the dead-letter message
```

Observe the logs for `BadAckFunc` to confirm repeated HTTP `500` responses and the logs for `AckPubMessage` to confirm eventual HTTP `200`.

---

## Multiple Subscriptions on a Topic

Not tested yet.

Each Pub/Sub subscription has independent:

* retry behavior
* delivery-attempt tracking
* push configuration
* dead-letter policy

Therefore multiple subscriptions on the same topic can use different retry counts and different dead-letter topics.

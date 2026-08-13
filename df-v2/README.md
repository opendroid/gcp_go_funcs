# Dialogflow using Cloud Run functions

Deploying Go functions for Dialogflow using **Cloud Run functions**.

Cloud Functions 2nd gen is now part of the Cloud Run functions model. Functions are deployed and managed as Cloud Run services using `gcloud run deploy`.

```shell
GCP_PROJECT=gcp-experiments-334602
GCP_REGION=us-central1
```

## Basic HTTP Function

Deploy `HelloShirts` as a Cloud Run function:

```shell
gcloud run deploy hello-shirts \
  --source=. \
  --function=HelloShirts \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --no-allow-unauthenticated
```

For Cloud Run functions:

* `--function` specifies the Go function entry point.
* `--base-image=go126` selects the Go 1.26 runtime.
* `--source=.` builds and deploys the function directly from the current source directory.
* `--no-allow-unauthenticated` keeps the service private.

Google documents `--allow-unauthenticated` and `--no-allow-unauthenticated` as the explicit controls for public versus authenticated Cloud Run services.

### Get service details

```shell
gcloud run services describe hello-shirts \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION"
```

### Get service URL

```shell
FUNCTION_URL=$(gcloud run services describe hello-shirts \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --format='value(status.url)')

echo "$FUNCTION_URL"
```

Cloud Run exposes the deployed service URL through `status.url`.

### Send test data

For an authenticated service:

```shell
curl -m 70 \
  -X POST "$FUNCTION_URL" \
  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -H "Content-Type: application/json" \
  -d '{"name":"Hello World"}'
```

The caller must have permission to invoke the Cloud Run service, such as `roles/run.invoker`.

---

## Dialogflow CX Webhook

Deploy `HandleWebhookRequest` as the `cx-webhook` Cloud Run service:

```shell
gcloud run deploy cx-webhook \
  --source=. \
  --function=HandleWebhookRequest \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --no-allow-unauthenticated
```

Get details:

```shell
gcloud run services describe cx-webhook \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION"
```

Get the service URL:

```shell
CX_FUNCTION_URL=$(gcloud run services describe cx-webhook \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --format='value(status.url)')

echo "$CX_FUNCTION_URL"
```

Send test data:

```shell
curl -m 70 \
  -X POST "$CX_FUNCTION_URL" \
  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -H "Content-Type: application/json" \
  -d '{
    "fulfillmentInfo": {
      "tag": "welcome"
    }
  }'
```

---

## Other Examples

### Echo function

Deploy `Echo` as `echo`:

```shell
gcloud run deploy echo \
  --source=. \
  --function=Echo \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --no-allow-unauthenticated
```

Deploy the same entry point as another service:

```shell
gcloud run deploy echo-v2 \
  --source=. \
  --function=Echo \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --no-allow-unauthenticated
```

---

## Testing Cloud Run Functions

These examples use the personal GCP project:

```shell
GCP_PROJECT=gcp-experiments-334602
GCP_REGION=us-central1
```

### Deploy an unauthenticated function

Deploy `HelloShirts` as the public `go-all` service:

```shell
gcloud run deploy go-all \
  --source=. \
  --function=HelloShirts \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --allow-unauthenticated
```

Get details:

```shell
gcloud run services describe go-all \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION"
```

Get its URL:

```shell
GO_ALL_URL=$(gcloud run services describe go-all \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --format='value(status.url)')

echo "$GO_ALL_URL"
```

Because the service allows unauthenticated access, it can be invoked without an ID token:

```shell
curl "$GO_ALL_URL"
```

Or send POST data:

```shell
curl -m 70 \
  -X POST "$GO_ALL_URL" \
  -H "Content-Type: application/json" \
  -d '{"name":"Hello World"}'
```

`--allow-unauthenticated` grants public invocation access to the Cloud Run service.

---

### Deploy an authenticated function

Deploy `HelloShirts` as the private `go-auth` service:

```shell
gcloud run deploy go-auth \
  --source=. \
  --function=HelloShirts \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --no-allow-unauthenticated
```

Get its URL:

```shell
GO_AUTH_URL=$(gcloud run services describe go-auth \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --format='value(status.url)')

echo "$GO_AUTH_URL"
```

Invoke it with an identity token:

```shell
curl -m 70 \
  -X POST "$GO_AUTH_URL" \
  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -H "Content-Type: application/json" \
  -d '{"name":"Hello World"}'
```

Cloud Run authenticated HTTP requests use an ID token in the `Authorization: Bearer` header.

---

## Dialogflow CX

Deploy the Dialogflow CX webhook:

```shell
gcloud run deploy cx-webhook \
  --source=. \
  --function=HandleWebhookRequest \
  --base-image=go126 \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --no-allow-unauthenticated
```

Get the URL:

```shell
CX_FUNCTION_URL=$(gcloud run services describe cx-webhook \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --format='value(status.url)')

echo "$CX_FUNCTION_URL"
```

Test the webhook:

```shell
curl -m 70 \
  -X POST "$CX_FUNCTION_URL" \
  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -H "Content-Type: application/json" \
  -d '{
    "fulfillmentInfo": {
      "tag": "welcome"
    }
  }'
```

---

## Cloud Functions v2 → Cloud Run Command Mapping

| Cloud Functions v2          | Cloud Run functions                         |
| --------------------------- | ------------------------------------------- |
| `gcloud functions deploy`   | `gcloud run deploy`                         |
| `--gen2`                    | Not required                                |
| `--runtime=go126`           | `--base-image=go126`                        |
| `--entry-point=HelloShirts` | `--function=HelloShirts`                    |
| `--trigger-http`            | Not required for an HTTP Cloud Run function |
| `--source=.`                | `--source=.`                                |
| `gcloud functions describe` | `gcloud run services describe`              |
| `serviceConfig.uri`         | `status.url`                                |
| `gcloud functions call`     | Use HTTP/`curl` against the Cloud Run URL   |

Google describes a Cloud Run function as a Cloud Run service deployed from source and supports `gcloud run deploy` for the current function resource model.

## Notes

* Prefer `gcloud run deploy` for new Cloud Run functions.
* Use `--function=<GO_ENTRY_POINT>` for Functions Framework-style Go functions.
* Use `--base-image=go126` for Go 1.26.
* Use `--allow-unauthenticated` only when the endpoint should be public.
* Use `--no-allow-unauthenticated` for private endpoints.
* For private endpoints, callers require Cloud Run invocation permission and must provide a valid ID token.
* Retrieve the deployed endpoint with:

```shell
gcloud run services describe SERVICE \
  --project="$GCP_PROJECT" \
  --region="$GCP_REGION" \
  --format='value(status.url)'
```

# API Gateway

### Creating an API config

```
gcloud beta api-gateway api-configs create gae-first \
  --api=gae-hello-api --openapi-spec=openapi2-appengine.yaml \
  --project=sinmetal-api-gateway-us --backend-auth-service-account=sinmetal-api-gateway-us@appspot.gserviceaccount.com

gcloud beta api-gateway api-configs delete gae-first \
  --api=gae-hello-api --project=sinmetal-api-gateway-us
```

updateで --openapi-spec=openapi2-appengine.yaml が更新できないから、こいつはなるべく更新しない構成の方がいいのか？
以下のように書いてあるから、BETAの時だけなのかもしれないが

```
(BETA) Update an API Gateway API config.
NOTE: Only the name and labels may be updated on an API config.
```

### Deploying an API Gateway

```
gcloud beta api-gateway gateways create gae-gateway-first \
  --api=gae-hello-api --api-config=gae-first \
  --location=us-central1 --project=sinmetal-api-gateway-us

gcloud beta api-gateway gateways describe gae-gateway-first \
  --location=us-central1 --project=sinmetal-api-gateway-us
```
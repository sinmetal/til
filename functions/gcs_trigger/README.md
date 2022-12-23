```
export PROJECT_ID=sinmetal-gcs-trigger-20221222
export SERVICE_ACCOUNT=gcs-trigger@sinmetal-gcs-trigger-20221222.iam.gserviceaccount.com
```

```
GCS_SERVICE_ACCOUNT="$(gsutil kms serviceaccount)"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:${GCS_SERVICE_ACCOUNT}" \
    --role='roles/pubsub.publisher'
```

```
SERVICE_ACCOUNT=gcs-trigger@sinmetal-gcs-trigger-20221222.iam.gserviceaccount.com
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:${SERVICE_ACCOUNT}" --role='roles/eventarc.eventReceiver'
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:${SERVICE_ACCOUNT}" --role='roles/cloudfunctions.invoker'
```

```
gcloud functions deploy gcs-trigger-tokyo-function \
  --gen2 \
  --region=asia-northeast1 \
  --runtime=go119 \
  --source=./ \
  --entry-point=GCSTriggerFunction \
  --trigger-event-filters="type=google.cloud.storage.object.v1.finalized" \
  --trigger-event-filters="bucket=sinmetal-gcs-trigger-tokyo" \
  --trigger-service-account="gcs-trigger@sinmetal-gcs-trigger-20221222.iam.gserviceaccount.com"
```
# impersonate service account

https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials example

## error sample

impersonateしたいServiceAccountのTokenCreatorの権限を持っていない
```
Post "https://storage.googleapis.com/upload/storage/v1/b/hoge/o?alt=json&name=2021-05-18+18%3A57%3A56.12854+%2B0900+JST+m%3D%2B0.003766193&prettyPrint=false&projection=full&uploadType=multipart": impersonate: status code 403: {
  "error": {
    "code": 403,
    "message": "The caller does not have permission",
    "status": "PERMISSION_DENIED"
  }
}
```

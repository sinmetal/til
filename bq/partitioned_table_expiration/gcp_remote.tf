terraform {
  backend "gcs" {
    bucket = "tfstate20230306"
    prefix = "terraform/state"
  }
}
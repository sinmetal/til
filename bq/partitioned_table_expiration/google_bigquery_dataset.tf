resource "google_bigquery_dataset" "first_dataset" {
  project       = "terraform20230306"
  dataset_id    = "first_dataset"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"
}
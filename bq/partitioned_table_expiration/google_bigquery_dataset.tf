resource "google_bigquery_dataset" "first_dataset" {
  dataset_id                  = "first_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}
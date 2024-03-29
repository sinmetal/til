resource "google_bigquery_table" "bar" {
  project    = "terraform20230306"
  dataset_id = google_bigquery_dataset.first_dataset.dataset_id
  table_id   = "bar"

  time_partitioning {
    type          = "DAY"
    expiration_ms = 10 * 24 * 3600 * 1000
  }

  schema = <<EOF
[
  {
    "name": "permalink",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "The Permalink"
  },
  {
    "name": "state",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "State where the head office is located"
  }
]
EOF

}

resource "google_bigquery_table" "foo" {
  project    = "terraform20230306"
  dataset_id = google_bigquery_dataset.first_dataset.dataset_id
  table_id   = "foo"

  time_partitioning {
    type          = "DAY"
    expiration_ms = 1 * 24 * 3600 * 1000
  }

  schema = <<EOF
[
  {
    "name": "permalink",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "The Permalink"
  },
  {
    "name": "state",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "State where the head office is located"
  }
]
EOF

}

resource "google_bigquery_table" "hoge" {
  project    = "terraform20230306"
  dataset_id = google_bigquery_dataset.first_dataset.dataset_id
  table_id   = "hoge"

  time_partitioning {
    type = "DAY"
  }

  schema = <<EOF
[
  {
    "name": "permalink",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "The Permalink"
  },
  {
    "name": "state",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "State where the head office is located"
  }
]
EOF

}
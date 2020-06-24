provider "airflow" {
  base_endpoint = "http://localhost:28080/"
}

resource "airflow_variable" "foo" {
  key   = "foo"
  value = "bar"
}

resource "airflow_variable" "hello" {
  key   = "hello"
  value = "world"
}

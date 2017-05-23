nginx {
  address           = "http://localhost:8080"
  status_endpoint   = "status"
  upstream_endpoint = "upstream_conf"
}

mapping {
  entry "sample-static-project" {
    type = "static"

    upstreams = [
      "10.144.253.11:7480",
      "10.144.253.12:7480",
    ]
  }
}

mapping {
  entry "sample-consul-project" {
    type         = "consul"
    service      = "service-name"
    delay_remove = "10s"
    delay_insert = "2s"
  }

  entry "sample-consul-project-defaults" {
    type    = "consul"
    service = "service-name"
  }
}

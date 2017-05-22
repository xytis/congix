nginx {
  address           = "http://localhost:8080"
  status_endpoint   = "status"
  upstream_endpoint = "upstream_conf"
}

mapping {
  entry "sample-static-project" {
    static = {
      server = "10.144.253.11:7480"
    }

    static = {
      server = "10.144.253.12:7480"
    }
  }
}

mapping {
  entry "sample-consul-project" {
    service = "service-name"
  }
}

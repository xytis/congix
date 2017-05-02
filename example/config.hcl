nginx {
  address            = "http://localhost:8080"
  status_endpoint    = "status"
  upstreams_endpoint = "upstream_conf"
}

consul {
  address = "http://localhost:8500"
}

mapping {
  sample-static-project {
    static = {
      server = "10.144.253.11:7480"
    }

    static = {
      server = "10.144.253.12:7480"
    }
  }

  sample-consul-project {
    service = "service-name"
  }
}

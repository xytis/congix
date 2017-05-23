# congix
Daemon which connects nginx-plus API with consul service catalog

## Building

    dep ensure -update
    go build
    # Optional
    gox "linux darwin" "386 amd64" -ouput="dist/congix_{{.OS}}-{{.Arch}}"

## Testing

    go test ./...

## Principles

Congix works with two logical states of services in consul:

- undefined
- healthy

Any other state is aggregated to one of the two states above.
Congix requires a configuration file which maps nginx-plus upstream zone name to upstream server list.
List can be static or dynamic:

    entry "nginx-plus-upstream-zone-name" {
      type = "consul"
      service = "consul-service-name"
      delay_remove = 10s
      delay_insert = 2s
    }

    entry "nginx-plus-upstream-zone-name" {
      type = "static"
      upstreams = [
         "127.0.0.1:8080",
      ]
    }

Congix stores last known state in disk, on change it performs diff based update to nginx API.

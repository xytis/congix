# congix
Daemon which connects nginx-plus API with consul service catalog

## Building

    dep ensure -update
    go build
    # Optional
    gox "linux darwin" "386 amd64" -ouput="dist/congix_{{.OS}}-{{.Arch}}"

## Testing

    go test ./...

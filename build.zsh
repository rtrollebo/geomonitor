# Build for embedded architecture:

VERSION=$(git describe --always --dirty)
GOOS=linux GOARCH=arm64 go build -ldflags="-X main.version=${VERSION}" -o bin/geomonitor main.go


run:
	go run cmd/main.go

test:
	go test github.com/htaunay/noise

depends:
	go get -u github.com/spf13/cobra 
	go get -u github.com/spf13/viper 

PROJECT_PATH=github.com/htaunay/noise

cli:
	go run noise-cli/main.go

install:
	go install $(PROJECT_PATH)
	go install $(PROJECT_PATH)/noise-cli

test:
	go test $(PROJECT_PATH)

depends:
	go get -u github.com/spf13/cobra 
	go get -u github.com/spf13/viper 

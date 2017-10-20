PROJECT_PATH=github.com/htaunay/noise

cli:
	go run noise-cli/main.go

install:
	go install $(PROJECT_PATH)
	go install $(PROJECT_PATH)/noise-cli
	go install $(PROJECT_PATH)/noise-gui

test:
	go test $(PROJECT_PATH) -v

depends:
	go get -u -v github.com/spf13/cobra 
	go get -u -v github.com/spf13/viper 
	go get -u -v github.com/go-gl/glfw
	go get -u -v github.com/faiface/glhf
	go get -u -v github.com/faiface/pixel

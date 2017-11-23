PROJECT_PATH=github.com/htaunay/noise

cli:
	go run noise-cli/main.go

gui:
	go run noise-gui/main.go

install: depends
	go install $(PROJECT_PATH)
	go install $(PROJECT_PATH)/noise-cli
	go install $(PROJECT_PATH)/noise-gui

test:
	go test $(PROJECT_PATH) -v

depends:
	go get -v github.com/spf13/cobra
	go get -v github.com/spf13/viper
	go get -v github.com/go-gl/glfw/v3.2/glfw
	go get -v github.com/faiface/glhf
	go get -v github.com/faiface/pixel

travis:
	go get -v github.com/spf13/cobra
	go get -v github.com/spf13/viper
	go get -v github.com/faiface/glhf
	go get -v github.com/faiface/pixel
	go install $(PROJECT_PATH)
	go install $(PROJECT_PATH)/noise-cli
	go test $(PROJECT_PATH)
	

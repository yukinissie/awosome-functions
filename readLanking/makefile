.PHONY: local_run dev_build dev_deploy

dev_build:
	mkdir -p ./build ./build/develop
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o ./build/develop/main main.go
	zip ./build/develop/deployment.zip ./build/develop/main

local_run: dev_build
	cd .. && sls invoke local -f readLanking -d '{}' --verbose

dev_deploy: dev_build
	cd .. && sls deploy

# productionへはGitHub Actionsで

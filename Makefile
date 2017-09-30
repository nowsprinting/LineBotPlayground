# for local test execute
export LINEBOT_CHANNEL_SECRET=012345678901234567890123456789ab
export LINEBOT_CHANNEL_ACCESS_TOKEN=012345678901234567890123456789ab012345678901234567890123456789ab

ifdef RUN
	RUNFUNC := -run $(RUN)
endif

version:
	echo package backend > backend/version.go
	echo const version = \"$(shell git describe --tags)\" >> backend/version.go

test: version
	gcloud config set project testapp
	go test ./backend -v -covermode=count -coverprofile=coverage.out $(RUNFUNC)

deploy: version
	gcloud config set project line-bot-playground
	gcloud app deploy backend/app.yaml --version 1

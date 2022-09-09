test:
	go test ./... -coverprofile cover.prof && go tool cover -func cover.prof

static:
	staticcheck ./...

cover:
	go test ./... -coverprofile cover.prof && go tool cover -html=cover.prof

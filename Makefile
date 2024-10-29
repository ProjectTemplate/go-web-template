
.PHONY: build-admin
build-admin:
	cd ./app/admin/cmd/server &&  go build .

.PHONY: run-admin
run-admin:
	cd ./app/admin/cmd/server && go run main.go
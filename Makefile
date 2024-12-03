
.PHONY: build-admin
build-admin:
	cd ./app/admin/cmd/server &&  go build .

.PHONY: run-admin
run-admin:
	cd ./app/admin/cmd/server && go run main.go


.PHONY: test
test:
	cd ./base/common/constant && go test -gcflags "all=-N -l" ./...
	cd ./base/common/utils && go test -gcflags "all=-N -l" ./...
	cd ./base/lib/config && go test -gcflags "all=-N -l" ./...
	cd ./base/lib/logger && go test -gcflags "all=-N -l" ./...
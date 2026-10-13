GOPATH = cmd/bot
RUN = go run $(GOPATH)/main.go $(filter-out $@,$(MAKECMDGOALS))
.PHONY: test run ping

run:
	@$(RUN)

test:
	@go test ./...

ping:
	@go test -v -count=1 ./internal

vet:
	@go vet ./...

# %:
# 	@:

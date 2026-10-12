GOPATH = cmd/bot
RUN = go run $(GOPATH)/main.go $(filter-out $@,$(MAKECMDGOALS))
.PHONY: test run

test:
	@go test -v ./...

run:
	@$(RUN)

# %:
# 	@:

GOPATH = src/cmd/bot
RUN = go run $(GOPATH)/main.go $(filter-out $@,$(MAKECMDGOALS))

run:
	@$(RUN)

# %:
# 	@:

GOPATH = src/cmd/bot/main.go
RUN = go run $(GOPATH) $(filter-out $@,$(MAKECMDGOALS))

run:
	@$(RUN)

# %:
# 	@:


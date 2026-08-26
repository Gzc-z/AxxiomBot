RUN = go run main.go $(filter-out $@,$(MAKECMDGOALS))

run:
	@$(RUN)

# %:
# 	@:


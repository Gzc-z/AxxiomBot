SCRIPT = ./script.sh
RUN = go run main.go $(filter-out $@,$(MAKECMDGOALS))

run:
	@$(RUN) || true

run_wit_err:
	@$(run)

r:
	@$(SCRIPT)


%:
	@:


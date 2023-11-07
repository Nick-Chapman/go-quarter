
default = bbc

run: run-$(default)

run-%:
	go run . ../quarter-forth/$*.list

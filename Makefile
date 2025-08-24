.PHONY: run/ssg
run/ssg:
	@echo "Run cmd/ssg to generate static content"
	go run ./cmd/ssg

.PHONY: run/server
run/server:
	@echo "Run static server to deliver the content"
	go run ./cmd/web
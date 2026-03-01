# ─── Test ─────────────────────────────────────────────────────────────────────

test-unit:
	cd app && go test ./internal/validator/... ./internal/service/... ./internal/handler/... -v

test-integration:
	cd app && go test ./internal/repository/... -v -timeout 120s

test-coverage:
	cd app && go test -coverprofile=coverage.out ./... -timeout 120s
	cd app && go tool cover -func=coverage.out
	cd app && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: app/coverage.html"

test-coverage-ci:
	cd app && go test -coverprofile=coverage.out -covermode=atomic ./... -timeout 120s
	cd app && go tool cover -func=coverage.out

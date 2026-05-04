test:
	go clean -testcache
	go test -v -coverprofile=coverage.out "$(shell go list ./... | grep -v "examples")"
	go tool cover -html=coverage.out -o coverage.html

sd_parity:
	release=$$(cat stable_diffusion.release); \
	cd "$$SD" && git diff $$release HEAD -- include/stable-diffusion.h

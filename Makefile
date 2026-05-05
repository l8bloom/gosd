SHELL := /bin/bash


test:
	go clean -testcache
	go test -v -coverprofile=coverage.out "$(shell go list ./... | grep -v "examples")"
	go tool cover -html=coverage.out -o coverage.html

sd_parity:
	release=$$(cat stable_diffusion.release); \
	cd "$$SD" && git pull && git diff $$release HEAD -- include/stable-diffusion.h

run_gen_image_example:
	time go run examples/image_gen/image_gen.go

run_gen_video_example:
	time go run examples/video_gen/video_gen.go

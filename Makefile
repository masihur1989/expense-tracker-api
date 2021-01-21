dev: build-docs
	air

build-docs:
	swag init -g server.go --output docs
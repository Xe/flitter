all: build

build:
	@cd havok ; go get ./... ; go build .

container: build
	@cd havok ; docker build -t xena/havok .

.PHONY: all build container

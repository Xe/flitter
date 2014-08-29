NAME = deis/init
VERSION = 0.1.0

.PHONY: all build test tag_latest release ssh

all: build

build:
	docker build -t $(NAME):$(VERSION) --rm image

test:
	env NAME=$(NAME) VERSION=$(VERSION) ./test/runner.sh

shell:
	docker run --rm -it ${NAME}:${VERSION} /sbin/my_init -- /bin/bash

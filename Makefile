GIT_COMMIT := $(shell git rev-list -1 HEAD)
DT := $(shell date +%Y.%m.%d.%H%M%S)
ME := $(shell whoami)
HOST := $(shell hostname)
PRODUCT := infosec
REPOROOT := digitalcircle
PLATFORM := linux/amd64
STACK := foundation
PWD := $(shell pwd)

run:
	DSN="host=localhost user=xxx password=xxx dbname=xxx port=5432 sslmode=disable TimeZone=America/Sao_Paulo" \
	REDIS=redis://localhost:6379 \
	CGO_ENABLED=0 go run -ldflags "-X github.com/digitalcircle-com-br/buildinfo.Ver=$(GIT_COMMIT) -X github.com/digitalcircle-com-br/buildinfo.BuildDate=$(DT) -X github.com/digitalcircle-com-br/buildinfo.BuildUser=$(ME) -X github.com/digitalcircle-com-br/buildinfo.BuildHost=$(HOST) -X github.com/digitalcircle-com-br/buildinfo.Product=$(PRODUCT)" $(MAIN)

arm64:
	$(eval PRODUCT_TAG := arm64)

product-gateway:
	$(eval PRODUCT := gateway)

product-auth:
	$(eval PRODUCT := auth)

product-authmgr:
	$(eval PRODUCT := authmgr)

product-static:
	$(eval PRODUCT := static)

docker-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deploy/$(PRODUCT)/main \
	-ldflags "-X github.com/digitalcircle-com-br/buildinfo.Ver=$(GIT_COMMIT) -X github.com/digitalcircle-com-br/buildinfo.BuildDate=$(DT) -X github.com/digitalcircle-com-br/buildinfo.BuildUser=$(ME) -X github.com/digitalcircle-com-br/buildinfo.BuildHost=$(HOST) -X github.com/digitalcircle-com-br/buildinfo.Product=$(PRODUCT)" \
	./cmd/${PRODUCT}/main.go && \
	cd deploy/$(PRODUCT) && \
	docker build --platform $(PLATFORM) -t $(REPOROOT)/$(STACK)-$(PRODUCT):latest .
	
docker-push: 
	docker push  $(REPOROOT)/$(STACK)-$(PRODUCT):latest

docker-pull:
	docker pull  $(REPOROOT)/$(STACK)-$(PRODUCT):latest

docker-build-push: docker-build-product docker-push


## Images - will build each image as per above

img-gateway: product-gateway docker-build docker-push

img-auth: product-auth docker-build docker-push

img-static: product-static docker-build docker-push
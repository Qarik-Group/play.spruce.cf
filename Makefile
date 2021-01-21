IMAGE_NAME ?= starkandwayne/play-spruce-cf
IMAGE_TAG  ?= latest

build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

push: build
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

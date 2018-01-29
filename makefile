IMAGE = stanfoot/s3-keep-files-amount
IMAGE_TAR = docker-image.tar

.PHONY: clean
clean:
	@rm -rf vendor
install:
	@dep ensure -v
update:
	@dep ensure -update -v
docker-build:
	@docker build . -t ${IMAGE}:latest
docker-save:
	@docker image save ${IMAGE}:latest > ./${IMAGE_TAR}
docker-load:
	@docker load -i ./${IMAGE_TAR}
docker-test:
	@structure-test -test.v -image ${IMAGE}:latest container-structure-test.yaml
docker-tag:
	@docker tag ${IMAGE}:latest ${IMAGE}:${tag}
docker-push:
	@docker push ${IMAGE}
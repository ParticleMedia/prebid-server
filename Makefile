# Makefile
GIT_BRANCH=$(strip $(shell git symbolic-ref --short HEAD))
AVRO_OUTPUT_DIR=./logging/model
AVRO_SCHEMA_DIR=./logging/schema

pr:
	@echo "Creating pull request..."
	@git push --set-upstream origin $(GIT_BRANCH)
	@hub pull-request

all: deps test build

.PHONY: deps test build image

# deps will clean out the vendor directory and use go mod for a fresh install
deps:
	GOPROXY="https://proxy.golang.org" go mod vendor -v && go mod tidy -v
	
# test will ensure that all of our dependencies are available and run validate.sh
test: deps
# If there is no indentation, Make will treat it as a directive for itself; otherwise, it's regarded as a shell script.
# https://stackoverflow.com/a/4483467
ifeq "$(adapter)" ""
	./validate.sh
else
	go test github.com/prebid/prebid-server/adapters/$(adapter) -bench=.
endif

# build will ensure all of our tests pass and then build the go binary
build:
	@go build -mod=vendor

# image will build a docker image
image:
	@docker build -t prebid-server .


run:
	@./prebid-server -logtostderr=1

avro:
	-rm -rf $(AVRO_OUTPUT_DIR)
	@mkdir $(AVRO_OUTPUT_DIR)
	@gogen-avro --containers=false \
				--sources-comment=false \
				--short-unions=false \
				--package=model \
				$(AVRO_OUTPUT_DIR) $(AVRO_SCHEMA_DIR)/*.avsc

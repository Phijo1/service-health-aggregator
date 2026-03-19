# Service Health Aggregator Assessment

## How to run

The service can be run using `go run main.go` in the base directory. It can also be built as a Docker image, and run using Docker:

`docker build -t health-aggregator .`

`docker run -p 8080:8080 health-aggregator`

Once running, the service is available at:

http://127.0.0.1:8080/health/aggregate

## Decisions made

For the sake of saving time, I included everything in `main.go`, the config loading, the http requests, and the aggregation. If I'd had more time, I would have split these out into their own files and directories for a modular structure. Something like:

`cmd/server/main.go` The actual service.

`internal/config/config.go` Reads config.

`internal/http/handler.go` Handles the service http output.

`internal/health/check.go` Checks input url health.

`internal/health/aggregator.go` Aggregates service health and returns healthy/degraded/down.

The separation would improve testability, and maintainability, however dealing with the boilerplate of imports and project structure would have been a time sink. Also, the complexity of the service is low, so the single `main.go` file can still easily be understood.

The Dockerfile was simple enough, I broke it up into two parts: the builder which builds the golang binary and copies all of the repo's code in to do so, and the actual output image which should just have the binary and config on it. I chose to use alpine because its incredibly lightweight, and does all I need it to.

The Github Actions make use of existing actions to test the code and build the docker image, though nothing is done with the image. To extend it, I'd add a step uploading it to Dockerhub or whatever equivalent I was using.

The Terraform was the bit I am least familiar with, I utilised AI (chatgpt) to help build it out and teach me what each section was doing in detail. The region is split out into the variables file so that it can be overwritten on `terraform apply`, and the outputs file specifies the service url as an output so we don't have to go digging to find it.

Lastly, the tests make use of golang's excellent test framework. They can be run using `go test ./...` (this would be helpful if there were other tests in subdirectories had I split everything out). I have one test which verifies the expected behaviour of the service, and then one which tests the timeout handling, and another which tests the aggregator function, in a more typical functional unit test manner.
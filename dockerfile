# This is a multi-stage image build example
# Which will help :
# - shorten the image size
# - improve security by not including the Go toolchain and other dependencies in the actual container
FROM golang:1.19 AS build-stage

WORKDIR /app

COPY  go.mod go.sum ./
RUN go mod download

COPY . ./
COPY ./.env ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /rules-engine

FROM build-stage AS run-test-stage
RUN go test -v -coverprofile cover.out ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /

COPY --from=build-stage /rules-engine /rules-engine
COPY rules/ rules/
COPY .env .

EXPOSE 4301

CMD ["/rules-engine"]
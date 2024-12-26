FROM golang:1.23 AS build-stage
LABEL authors="obinn"
WORKDIR /build
COPY ../go.mod go.sum ./
RUN go mod download
COPY .. .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main /build/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /main /main
USER nonroot:nonroot
ENTRYPOINT ["/main"]
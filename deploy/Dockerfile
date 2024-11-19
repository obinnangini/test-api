

FROM golang:1.22
LABEL authors="obinn"
WORKDIR /build
COPY ../go.mod go.sum ./
RUN go mod download
COPY .. .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main /build/main.go
CMD ["/main"]
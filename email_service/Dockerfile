FROM golang:1.22.5 as build_stage

WORKDIR /app

# Copy lock files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Add project files
COPY *.go ./

# Compile Application
RUN CGO_ENABLED=0 GOOS=linux go build -o /mail_service

# Test In Container
FROM build_stage as run_tests
RUN test -v ./..

# Run in separate container
FROM gcr.io/distroless/base-debian11 as release_stage

WORKDIR /

COPY --from=build_stage /mail_service /mail_service
USER nonroot:nonroot

ENTRYPOINT ["/mail_service"]
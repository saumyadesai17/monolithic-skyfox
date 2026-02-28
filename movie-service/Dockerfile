
FROM artifactory.idfcfirstbank.com/neev-docker/golang:1.22 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY main.go ./
COPY movies.json ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o movie-service .

# Final stage - use minimal image
FROM artifactory.idfcfirstbank.com/neev-docker/golang:alpine

WORKDIR /app

# Copy binary and data from builder
COPY --from=builder /app/movie-service .
COPY --from=builder /app/movies.json .

EXPOSE 4567

CMD ["./movie-service"]


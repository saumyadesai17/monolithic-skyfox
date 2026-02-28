# ---------- BUILD STAGE ----------
FROM golang:1.22 AS builder

ARG VERSION=1.0
ARG OS=linux
ARG ARCH=amd64

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /code

# download deps first (cache friendly)
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build static binary
RUN CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} \
    go build -o /out/server main.go

# ---------- RUNTIME STAGE ----------
FROM debian:buster-slim

WORKDIR /app

COPY --from=builder /out/server ./out/server
COPY config/config-local.yml ./config.yml

EXPOSE 8080

ENTRYPOINT ["./out/server","-configFile","config.yml"]
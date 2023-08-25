.PHONY: pack build

PROJECT?=yh-process
AUTHOR?=khaosles
RELEASE?=1.0.0
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell TZ="Asia/Shanghai" date '+%Y-%m-%d_%H:%M:%S')
OUTDIR?=../../manifest/docker
TAG?=pack

pack:
	go-bindata -pkg packed -o ./internal/packed/packed.go ./manifest/config/config.yaml

build:
	go mod tidy && \
	cd ./cmd/gfs && GOOS=linux GOARCH=amd64  go build -tags "pack" -o ${OUTDIR}/yh-process-gfs   \
    && cd ../download && GOOS=linux GOARCH=amd64  go build -tags "pack" -o ${OUTDIR}/yh-process-download \
    && echo "build succeed!"


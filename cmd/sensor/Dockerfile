# base image
FROM golang:1.19.3-alpine AS builder
# create appuser
RUN adduser -D -g '' none
# create workspace
WORKDIR /opt/app/
COPY ../../go.mod ../../go.sum ./
# fetch dependancies
RUN go mod download && \
    go mod verify
# copy the source code as the last step
COPY . .
# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix 'static' -tags timetzdata -o /go/bin/sensor ./cmd/sensor


# build a small image
FROM alpine:3.17.0
LABEL language="golang"
# import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
# copy the static executable
COPY --from=builder --chown=none:1000 /go/bin/sensor /sensor
# use a non-root user
USER none
# run app
# ENTRYPOINT ["./sensor" ,"--mongoURI=mongodb://mongo-database:27017"]
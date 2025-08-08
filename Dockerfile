FROM --platform=$BUILDPLATFORM golang:1.23.4-alpine3.21 AS build
ARG HTTP_PORT
ARG TARGETPLATFORM
RUN echo "build for platform: $TARGETPLATFORM"
# Allow go to retrieve the dependencies for the build step
WORKDIR /go-modules
RUN apk update && apk upgrade && apk add --no-cache git && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . ./

# Compile the binary, we don't want to run the cgo resolver
ARG TARGETOS
ARG TARGETARCH
RUN echo "building for GOOS: $TARGETOS, GOARCH: $TARGETARCH"
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w"  -tags timetzdata -mod=vendor -a -installsuffix  -o market-service

# final stage
FROM scratch AS final
WORKDIR /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go-modules/market-service .
ENV TZ=Asia/Bangkok

VOLUME /conf
VOLUME /scripts

EXPOSE $HTTP_PORT
ENTRYPOINT ["./market-service"]

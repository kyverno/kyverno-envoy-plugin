FROM golang:1.22 as build
RUN go install golang.org/x/lint/golint@latest
WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o sidecar-injector

FROM scratch
WORKDIR /
COPY --from=build /build/sidecar-injector /

ENTRYPOINT ["/sidecar-injector"]

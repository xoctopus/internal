FROM golang:1.21-alpine AS builder

WORKDIR /go/src
COPY ./loopvar.go   ./loopvar.go
COPY ./go.mod.go121 ./go.mod

# build
#ARG COMMIT_SHA
RUN cat go.mod && go build .

# runtime
FROM gcr.io/distroless/static AS runtime

COPY --from=builder /go/src/loopvar /go/bin/loopvar
EXPOSE 80

WORKDIR /go/bin
ENTRYPOINT ["/go/bin/loopvar"]

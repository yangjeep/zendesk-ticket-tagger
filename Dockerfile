########################
# STEP 1 build go binary
########################
FROM golang:1.22.2-bullseye as builder
RUN apt-get update && apt-get install -yq build-essential
WORKDIR /app
COPY . .

# Ensure access to private repos, get dependencies, build binary
ARG TOKEN
RUN git config --global url."https://${TOKEN}:x-oauth-basic@github.com".insteadOf "https://github.com" \
    && go env -w GOPRIVATE=github.com/searchspring/* \
    && go get -d -v ./... \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=true -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/main .

############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /go/bin/main /go/bin/main
WORKDIR /go/bin
ENTRYPOINT ["./main"]

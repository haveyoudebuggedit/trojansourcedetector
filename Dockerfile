FROM golang AS build

RUN mkdir -p /srv/src
WORKDIR /srv/src
COPY . /srv/src
RUN go test -v ./...
RUN go build -o trojansourcedetector cmd/trojansourcedetector/main.go

FROM alpine

COPY --from=build /srv/src/trojansourcedetector /
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /trojansourcedetector /entrypoint.sh

RUN mkdir -p /github/workspace
VOLUME /github/workspace
WORKDIR /github/workspace

ENTRYPOINT ["/entrypoint.sh"]

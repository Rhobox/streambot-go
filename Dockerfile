FROM golang:alpine as build

WORKDIR /streambot

COPY . .

# Deployment needs alpine for a CA, so we allow CGO
RUN go build

# Using alpine instead of scratch so that we have an up-to-date certificate authority
FROM alpine:latest

COPY --from=build /streambot/streambot .

ENTRYPOINT ["./streambot"]

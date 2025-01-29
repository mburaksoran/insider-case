# Builder
FROM golang:1.23.5-alpine as builder

WORKDIR /app

RUN pwd && ls -l

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN  go build -o /insider-case github.com/mburaksoran/insider-case

# Final docker image
FROM alpine:latest AS final

RUN apk update \
    && apk upgrade\
    && apk add --no-cache tzdata curl

WORKDIR /app
COPY --from=builder /insider-case .
COPY --from=builder /app /app/

EXPOSE 8080
CMD [ "./insider-case" ]
# Builder stage
FROM golang:latest as builder

WORKDIR /build

COPY . /build/

ARG ENV_MONGO_DB_PASSWORD
RUN sed -i "s/{ENV_MONGO_DB_PASSWORD}/$ENV_MONGO_DB_PASSWORD/" /build/.env
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /build/main

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/main /app/
COPY --from=builder /build/.env /app/
RUN echo "this is .env in /app: "
RUN cat /app/.env

RUN chmod +x /app/main

EXPOSE 8080

CMD ["./main"]

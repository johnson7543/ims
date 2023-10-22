# Builder stage
FROM golang:latest as builder

WORKDIR /build

COPY . /build/

RUN sed -i "s/{ENV_MONGO_DB_PASSWORD}/$ENV_MONGO_DB_PASSWORD/" /build/.env
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /build/ims

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/ims /app/
COPY --from=builder /build/.env /app/
RUN echo "this is .env in /app: "
RUN cat /app/.env

RUN chmod +x /app/ims

EXPOSE 8080

CMD ["./ims"]

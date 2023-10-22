# Builder stage
FROM golang:latest as builder

WORKDIR /build

# Copy the entire application and .env file
COPY . /build/
RUN echo "this is .env in /build: "
RUN cat /build/.env
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

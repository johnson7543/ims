# Builder stage
FROM golang:latest as builder

WORKDIR /build

# Copy the entire application and .env file
COPY . /build/

RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /build/ims


# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the ims binary from the builder stage
COPY --from=builder /build/ims /app/

# Copy the .env file from the builder stage
COPY --from=builder /build/.env /app/
RUN chmod +x /app/ims
EXPOSE 8080

CMD ["./ims"]

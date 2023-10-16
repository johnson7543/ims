# Builder stage
FROM golang:latest as builder

WORKDIR /build

RUN echo "$ENV_VARIABLE_NAME" | base64 --decode >> .env

# Copy the entire application and .env file
COPY . /build/

RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /build/ims


# Final stage
FROM alpine:latest

WORKDIR /app

# Replace placeholder in .env file
RUN echo "$ENV_MONGO_DB_PASSWORD" | base64 --decode | tr -d '\n' > /tmp/mongo_password && \
    sed -i "s/{ENV_MONGO_DB_PASSWORD}/$(cat /tmp/mongo_password)/g" /build/.env

COPY --from=builder /build/.env /app/
COPY --from=builder /build/ims /app/


RUN chmod +x /app/ims
EXPOSE 8080

CMD ["./ims"]

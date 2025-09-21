FROM golang:1.25-alpine AS builder

# Set working directory di dalam container
WORKDIR /build

# Copy file dependency
COPY src/go.mod src/go.sum ./
RUN go mod download && go mod verify

COPY src/ .

ENV CGO_ENABLED=0

# Build aplikasi
RUN go build -v -o app .


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /build/app .

EXPOSE 8080

CMD ["./app"]

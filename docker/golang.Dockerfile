############################
# STEP 1 build executable binary
############################
FROM golang:1.22.5-alpine3.20 AS builder
RUN apk update && apk add --no-cache gcc musl-dev gcompat
WORKDIR /whatsapp
COPY ./src .

ENV GOCACHE=/root/.cache/go-build
# Fetch dependencies.
# RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download
# Build the binary.
RUN go build -o /app/whatsapp

#############################
## STEP 2 build a smaller image
#############################
FROM alpine:3.20
RUN apk update && apk add --no-cache ffmpeg
WORKDIR /app
# Copy compiled from builder.
COPY --from=builder /app/whatsapp /app/whatsapp
# Run the binary.
ENTRYPOINT ["/app/whatsapp"]
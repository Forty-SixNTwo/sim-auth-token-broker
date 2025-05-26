# Builder stage: compiles a single service with shared libs
FROM golang:1.24-alpine AS builder
ARG SERVICE_NAME=broker
WORKDIR /app
COPY services/${SERVICE_NAME}/go.mod services/${SERVICE_NAME}/go.sum services/${SERVICE_NAME}/
COPY libs/ ./libs/
COPY services/${SERVICE_NAME}/ ./services/${SERVICE_NAME}/
WORKDIR /app/services/${SERVICE_NAME}
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/${SERVICE_NAME} ./main.go

# Runtime stage
FROM alpine:latest
ARG SERVICE_NAME=broker
ARG PORT=8080
WORKDIR /app
COPY --from=builder /app/${SERVICE_NAME} ./${SERVICE_NAME}
COPY prefix_map.yaml .
COPY deploy.sh .
ENV PORT=${PORT}
EXPOSE ${PORT}
ENTRYPOINT ["./${SERVICE_NAME}"]

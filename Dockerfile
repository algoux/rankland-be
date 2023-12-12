FROM golang:alpine AS builder

LABEL stage=gobuilder

# Don't use libc
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

# Fix timezone
RUN apk update --no-cache && apk add --no-cache tzdata

# Build
WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/rankland ./main.go

FROM alpine

# Timezone & TLS
RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

# Copy application and init config file
WORKDIR /app/rankland
RUN mkdir -p ./config
RUN touch ./config/config.yaml
COPY --from=builder /app/rankland rankland

#Bind IO
EXPOSE 8000

# Run
ENV GIN_MODE=release
CMD ["./rankland"]

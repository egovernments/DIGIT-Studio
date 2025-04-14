FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o egov-generic-service

#FROM alpine:latest

#RUN apk --no-cache add ca-certificates

#WORKDIR /root/

#COPY --from=builder /app/egov-generic-service .
RUN chmod +x egov-generic-service
#CMD ["./egov-generic-service"]
CMD ["/bin/sh", "-c", "/app/egov-generic-service"]
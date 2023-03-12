FROM golang:1.20 AS gobuild
WORKDIR /build-dir
COPY go.mod .
COPY go.sum .
RUN go mod download all
COPY . .
RUN go build -o /tmp/ton-backend github.com/Contribution-DAO/cdao-ton-token-gate-core


FROM ubuntu AS tonproof
RUN apt-get update && \
    apt-get install -y openssl ca-certificates
COPY --from=gobuild /tmp/ton-backend /app/ton-backend
CMD ["/app/ton-backend"]



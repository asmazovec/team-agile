FROM golang:1.22.5-bookworm AS builder

WORKDIR /build/

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-extldflags -static" \
      github.com/asmazovec/team-agile/cmd/primary-server

#-

FROM debian:bookworm-slim

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates && \
      rm -rf /var/lib/apt/lists/*

WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=builder /build/primary-server ./bin/primary-server

EXPOSE 8080

CMD ["primary-server"]

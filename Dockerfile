# 빌드
FROM golang:1.24.3-bullseye AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app ./cmd/server

# 배포
FROM debian:bullseye-slim AS deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

ENV TZ=Asia/Seoul

CMD ["./app"]
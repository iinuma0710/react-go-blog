# デプロイ用のコンテナに含めるバイナリを作成するコンテナ
FROM golang:1.23-bullseye as deploy-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app


# デプロイ用のコンテナ
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .
CMD ["./app"]


# ローカル開発向けのホットリロード環境
FROM golang:1.23-bullseye as development
WORKDIR /app

RUN apt update && apt install git
RUN go install github.com/air-verse/air@latest

CMD ["air"]
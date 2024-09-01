ARG GO_VER=1.22
# デプロイ用コンテナに含めるバイナリを作成するコンテナ
FROM golang:$GO_VER-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldfags "-w -s" -o app

#-------------------------------------------------

# デプロイ用コンテナ
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]

#-------------------------------------------------

# ローカル開発環境で利用するホットリロード環境
FROM golang:$GO_VER as dev
WORKDIR /app
#RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/air-verse/air@latest
CMD ["air"]
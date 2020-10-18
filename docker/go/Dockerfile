# ベースとなるDockerイメージ指定
FROM golang:latest
RUN go get github.com/derekparker/delve/cmd/dlv \
&& go get github.com/go-sql-driver/mysql

# コンテナ内に作業ディレクトリを作成
RUN mkdir /go/src/work
# コンテナログイン時のディレクトリ指定
WORKDIR /go/src/work

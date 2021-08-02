# address-search-apis

[![ci](https://github.com/nekochans/address-search-apis/actions/workflows/ci.yml/badge.svg)](https://github.com/nekochans/address-search-apis/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/nekochans/address-search-apis/badge.svg?branch=main)](https://coveralls.io/github/nekochans/address-search-apis?branch=main)
[![cd-staging](https://github.com/nekochans/address-search-apis/actions/workflows/cd-staging.yml/badge.svg)](https://github.com/nekochans/address-search-apis/actions/workflows/cd-staging.yml)

住所検索に関する様々なAPIを提供します。

## Getting Started

AWS Lambda + Goで実装しています。

デプロイには [serverless framework](https://www.serverless.com/) を利用しています。

### AWSクレデンシャルの設定

[名前付きプロファイル](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-profiles.html) を利用する前提です。

### 環境変数の設定

[direnv](https://github.com/direnv/direnv) 等を利用して環境変数を設定します。

以下のように環境毎に設定を行います。

```
export DEPLOY_STAGE=stg, dev, qa などデプロイ先のステージを指定
export AWS_REGION=ap-northeast-1 などのAWSリージョンを指定
export AWS_PROFILE=利用しているAWS プロファイル名を指定
export KENALL_SECRET_KEY=https://kenall.jp/ のシークレットキーを指定
```

### Dockerで環境構築

以下で環境を立ち上げます。

```bash
docker-compose up --build -d
```

2回目以降は `docker-compose up -d` でOKです。

環境を再構築した際は `docker-compose up --build -d` を再度実行する必要があります。

`docker-compose exec go sh` でコンテナの中に入ります。

コンテナ内では以下のコマンドが利用可能です。

#### `make build`

アプリケーションのBuildを実行します。

これを実行すると `go.mod` に記載されているpackageをダウンロードしつつ `go.sum` を更新します。

packageのバージョンを統一する為に、この2つのファイルは必ずGitの管理対象として下さい。

#### `make lint`

linterを実行します。

本プロジェクトでは https://github.com/golangci/golangci-lint を利用します。

linterでのerrorが残っている状態では、PRがmerge出来ない設定を行っています。

#### `make format`

ソースコードをformatします。

`make lint` で出力されるerrorの一部は修正出来ますが `make format` だけでは修正出来ないerrorも多いのです。

`make format` が未実行の場合は `make lint` でerrorが出るようになっています。

#### `make test`

テストコードを実行します。

ローカルでの動作確認は基本的にはテストコードで担保します。

### Node.js のインストール

LTS版に指定されているバージョンであれば動作します。（最新版でもおそらく動作すると思います）

複数プロジェクトで異なる Node.js のバージョンを利用する可能性があるので、Node.js 自体をバージョン管理出来るようにしておくのが無難です。

以下は [nodenv](https://github.com/nodenv/nodenv) を使った設定例です。

```bash
nodenv install 14.15.3

nodenv local 14.15.3
```

### デプロイ

1. `npm ci` を実行（初回のみでOK）
1. `make serverless-deploy` を実行

ただしローカルでのデプロイは推奨されません。

もしも、作業中のPRをデプロイしたい場合は、GitHub Actionsを手動でデプロイ出来るようになっているので `.github/workflows/cd-staging.yml` 等を手動で実行してデプロイを行う事を推奨します。

## API仕様

### 共通仕様

リクエスト時のHTTPHeaderに `X-Request-Id` を含めると、そのままレスポンスヘッダーで同じ値を返します。

ログ出力の際にこの値が出力されるようになっているので、マイクロサービス間を横断したログ検索を行う際に利用する事を想定しています。

### 郵便番号から住所を検索する

以下のようにリクエストします。

```bash
curl -v \
-H "X-Request-Id: zzzzzzzzzzzzz-zzzz-(=^・^=)-dddddddd" \
https://{ドメイン名}/v1/1620062 | jq
```

以下のようなレスポンスが返ってきます。

```
< HTTP/2 200
< date: Mon, 02 Aug 2021 03:28:32 GMT
< content-type: application/json
< content-length: 87
< x-lambda-request-id: 91375e71-da96-42b7-9d05-d7bb6e1892da
< x-request-id: zzzzzzzzzzzzz-zzzz-(=^・^=)-dddddddd
< apigw-requestid: Dazu-g54NjMEPlw=
<
{ [87 bytes data]
100    87  100    87    0     0    125      0 --:--:-- --:--:-- --:--:--   125
* Connection #0 to host {ドメイン名} left intact
* Closing connection 0
{
  "postalCode": "1620062",
  "prefecture": "東京都",
  "locality": "新宿区市谷加賀町"
}
```

#### 郵便番号に該当する住所が存在しない場合はステータスコード404を返します

```
< HTTP/2 404
< date: Mon, 02 Aug 2021 03:31:39 GMT
< content-type: application/json
< content-length: 53
< x-lambda-request-id: 839338cb-3053-4379-ab36-bc37df3f5215
< x-request-id: zzzzzzzzzzzzz-zzzz-(=^・^=)-dddddddd
< apigw-requestid: Da0MWiMhNjMEMnQ=
<
{ [53 bytes data]
100    53  100    53    0     0    265      0 --:--:-- --:--:-- --:--:--   265
* Connection #0 to host {ドメイン名} left intact
* Closing connection 0
{
  "message": "住所が見つかりませんでした"
}
```


#### 郵便番号のフォーマットが不正の場合はステータスコード422を返します

```
< HTTP/2 422
< date: Mon, 02 Aug 2021 03:31:08 GMT
< content-type: application/json
< content-length: 74
< x-lambda-request-id: 9d64fe15-2ce8-46e5-a722-f799643a5fd9
< x-request-id: zzzzzzzzzzzzz-zzzz-(=^・^=)-dddddddd
< apigw-requestid: Da0HgiOgNjMEPgw=
<
{ [74 bytes data]
100    74  100    74    0     0    666      0 --:--:-- --:--:-- --:--:--   666
* Connection #0 to host {ドメイン名} left intact
* Closing connection 0
{
  "message": "郵便番号のフォーマットが正しくありません"
}
```

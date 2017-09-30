# LINE Bot Playground

LINE Messaging APIの機能を検証するためのApp



### 仕様

- Webhookで受け取った発言をオウム返しする



### Installation

#### Install golang libraries

    $ go get -u google.golang.org/appengine
    $ go get -u github.com/line/line-bot-sdk-go/linebot


#### app.yaml

`backend/app.yaml`にはLINE BOTのキー情報などを含むため、リポジトリから除外している。下記の書式でファイルを作成すること。

	runtime: go
	api_version: go1

	handlers:
	- url: /.*
	  script: _go_app

	env_variables:
	  LINEBOT_CHANNEL_SECRET: 'LINE Messaging APIチャネルのSECRET'
	  LINEBOT_CHANNEL_ACCESS_TOKEN: 'LINE Messaging APIチャネルのTOKEN'


#### version.go

バージョン番号は、`make test`、`make deploy`の際に生成されるversion.goファイルに定義される。このファイルはリポジトリ管理対象外。

`git clone`したプロジェクトを直接`goapp deploy`コマンドでApp Engineにデプロイすると、定数`version`が未定義なためエラーとなる。`make deploy`を使うこと。



### Test

ローカルでテストを実行する

	$ make test

特定のテストのみ実行する

	$ make test RUN=テスト関数名

なお、

- ローカル実行のとき、project idが`testapp`でないとデータストアに接続できないため、テスト実行時に`testapp `、デプロイ時に商用環境のProject IDを設定している。Project IDをリポジトリに残したくなかったが、やむを得ず。


### Deploy

	$ make deploy

# Logleaf

## 環境構築
### 1. git clone
```bash
git clone https://github.com/umekikazuya/logleaf.git
cd logleaf
```
### 2. Dockerイメージのビルド & コンテナの起動
```bash
docker-compose up --build
```

### 3. DynamoDB初期セットアップ
```bash
aws dynamodb create-table \
  --table-name leaves \
  --attribute-definitions \
    AttributeName=pk,AttributeType=S \
    AttributeName=sk,AttributeType=S \
  --key-schema \
    AttributeName=pk,KeyType=HASH \
    AttributeName=sk,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --endpoint-url http://localhost:8000
```

## Qiita APIトークンの設定
### 1. Qiitaのアクセストークンを取得
Qiitaのアカウントからアクセストークンを生成します。
参考: https://qiita.com/maiamea/items/680cca06f7825595cba0

### 2. 環境変数に設定
`.env.example`をコピーして`.env`ファイルを作成し、アクセストークンを設定します。

```bash
cp .env.example .env
```

`.env`ファイルを開いて、以下のようにアクセストークンを設定します。

```dotenv
QIITA_ACCESS_TOKEN=your_qiita_access_token
```

## API仕様


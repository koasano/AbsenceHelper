# 使用例

1. 終日指定の場合:
./AbsenceHelper.exe -s "Asano 終日不在" -d 2023-07-15

2. 複数日指定の場合:
./AbsenceHelper.exe -s "Asano 終日不在" -d 2023-07-15 -de 2023-07-20

3. 時間指定の場合:
./AbsenceHelper.exe -s "Asano AM不在" -d 2023-07-16 -tb 09:00 -te 13:00


# 注意事項
このツールには credentials.json と config.json が必要です。

## credentials.json
Google Cloud Platformのサービスアカウント認証情報を取得し、credentials.jsonという名前のファイルに保存する必要があります。
詳しくはドキュメントを参照ください。
https://cloud.google.com/iam/docs/service-accounts-create?hl=ja

## config.json
config.jsonファイルは、AbsenceHelperの設定を管理するためのJSON形式の設定ファイルです。このファイルは以下のキーを持ちます:

1. calendar_id: イベントを追加するGoogleカレンダーのIDを指定します。これは通常、メールアドレスの形式をとります。

2. time_zone: イベントの時間帯を指定します。デフォルトは"Asia/Tokyo"ですが、必要に応じて他のタイムゾーンに変更することができます。

3. language: AbsenceHelperが表示するメッセージの言語を指定します。現在、日本語("ja")と英語("en")が利用可能です。

### Example
以下はconfig.jsonの例です:

```
{
    "calendar_id": "example@gmail.com",
    "time_zone": "Asia/Tokyo",
    "language": "ja"
}
```

この例では、カレンダーIDは"example@gmail.com"、時間帯は"Asia/Tokyo"、言語は日本語("ja")と設定されています。

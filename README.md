# ping

## 概要

- ping コマンドを自分で実装しています

## できること

- ICMP Echo Requestを送る
- ICMP Echo Reqlyを受け取って表示する
- 相手のIPv4アドレスを指定して送信する
- `-s`オプションで、送信するbyte長を指定する
- `-c`オプションで、送信する回数を指定する
- `-w`オプションで、タイムアウト時間を指定する

## 使い方

- リポジトリをクローンする
- `sudo go run ./cmd/ping <IPv4 アドレス>`をする
- あるいは、`go build ./cmd/ping`を実行してから、`sudo ./ping <IPv4 アドレス>`を行う
- 使用例：`sudo ./ping -c 10 -s 52 -W 3 8.8.8.8`

## 参考

- [RFC 792}](https://www.rfc-editor.org/info/rfc792)
- [Wikipedia Internet Control Message Protocol](https://ja.wikipedia.org/wiki/Internet_Control_Message_Protocol)

## 注意

- これは学習目的で作成されたものであり、実際のネットワーク環境での使用は推奨されません。
- 作成にあたって、ChatGPTとの対話を利用しました。

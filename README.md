# ping

## 概要
- ping コマンドを自分で実装しています

## できること
- ICMP Echo Requestを送る
- ICMP Echo Reqlyを受け取って表示する
- 相手のIPv4アドレスを指定して送信する
- `-s`オプションで、送信するbyte長を指定する
- `-c`オプションで、送信する回数を指定する

## 使い方
- リポジトリをクローンする
- `sudo go run . <IPv4 アドレス>`をする
- あるいは、`go build`を行い、`sudo ./ping <IPv4 アドレス>`を行う
- 使用例：`sudo ./ping -c 10 -s 52 8.8.8.8`

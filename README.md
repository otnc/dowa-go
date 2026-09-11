# dowa-go

[![Go Reference](https://pkg.go.dev/badge/github.com/otnc/dowa-go.svg)](https://pkg.go.dev/github.com/otnc/dowa-go)
[![CI](https://github.com/otnc/dowa-go/actions/workflows/ci.yml/badge.svg)](https://github.com/otnc/dowa-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> *dowa* means どわーw (冷笑)

冷笑を検知しますw

> [!Note]
>   
> [dowa](https://github.com/otnc/dowa) (npm) のGo移植版です。

## 動作環境

- Go >= 1.21

## インストール

```bash
go get github.com/otnc/dowa-go
```

## 使い方

```go
package main

import (
	"fmt"

	"github.com/otnc/dowa-go"
)

func main() {
	fmt.Println(dowa.Contains("うおw")) // true

	// relaxed モード(検知範囲を拡大)
	fmt.Println(dowa.Contains("どわー", dowa.Options{Relaxed: true})) // true

	fmt.Println(dowa.FindAll("←うおw、爆笑爆笑")) // [うおw 爆笑爆笑]

	// どのパターンにマッチしたかの詳細がほしい場合
	for _, m := range dowa.FindMatches("うおw、爆笑爆笑") {
		fmt.Printf("%s (pattern=%s, strict=%v)\n", m.Text, m.PatternID, m.Strict)
	}
}
```

動くサンプルは [examples/basic](./examples/basic) にもあります。

```bash
go run ./examples/basic
```

## API

- `FindAll(text string, opts ...Options) []string` — マッチした冷笑の配列を返すw(見つからなければ`nil`)
- `Contains(text string, opts ...Options) bool` — 冷笑が含まれるかを真偽値で返すw
- `FindMatches(text string, opts ...Options) []Match` — どのパターンにマッチしたかの詳細(位置・パターンid・strict/relaxed)付きで返すw(見つからなければ`nil`)。パターンごとに個別検索するため、複数パターンの一致範囲が重なる場合は`FindAll`の重複排除された結果とは一致しないことがある
- `Options`
  - `Relaxed bool` (デフォルト: `false`) — `true`で検知範囲を拡大する
- `Match`
  - `Text string` — マッチした文字列
  - `Index int` — マッチ開始位置(バイトオフセット。TypeScript版はUTF-16コードユニット単位のindexを返すが、Goの文字列はUTF-8バイト列なのでこちらはバイトオフセットになる)
  - `PatternID string` — マッチしたパターンのid(内部の定義は[internal/patterns/patterns.go](./internal/patterns/patterns.go)参照)
  - `Strict bool` — そのパターンがstrictかどうか
- `Patterns() []PatternInfo` — 検知に使われている全パターンの概要(id/strict/サンプル)を返すw。パターンを紹介・デバッグしたい場合に

オプション引数は`Options`構造体を渡すことで指定する。省略した場合は`Options{}`(既定の挙動)として扱われる。

```go
dowa.Contains(text)                            // 既定(strict)
dowa.Contains(text, dowa.Options{Relaxed: true}) // relaxed
```

## TypeScript版との違い

- `dowaSchema`(Standard Schema対応)はTS/Zod/Valibotエコシステム固有の機能のため移植していません。バリデーションに組み込みたい場合は`FindMatches`の結果からエラーメッセージを組み立ててください
- `Match.Index`は前述の通りバイトオフセットです

貢献方法については [CONTRIBUTING.md](./CONTRIBUTING.md) を確認してください。

## その他

Twitter (新X) で冷笑ツイートを検知するChrome/Firefox拡張機能を公開中です。

https://github.com/otnc/dowa-twitter-checker

[Chrome拡張機能](https://chromewebstore.google.com/detail/kkojaplhlbbildhofdophfadmbdholdn) / [Firefox拡張機能](https://addons.mozilla.org/firefox/addon/%E5%86%B7%E7%AC%91%E3%83%81%E3%82%A7%E3%83%83%E3%82%AB%E3%83%BC-for-twitter/)

## 著者

otoneko. https://github.com/otnc

## ライセンス

[MIT License](./LICENSE)

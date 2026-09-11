# 貢献ガイド

PRに適切なラベルを付与してください。

レビューが必要な場合やしばらく反応がない場合は、PRのコメントやDiscordでメンションしてください。

## 開発

```bash
go build ./...
go vet ./...
gofmt -l .   # フォーマット崩れがあればファイル名が出力される
go test ./...
```

## 必要な操作

1. このリポジトリをforkしてからローカルにクローンします
2. 必要な修正や追加を行う
3. `gofmt -l .` で何も出力されないことを確認する(出力があれば `gofmt -w .` で整形)
4. `go vet ./...` でエラーが出ないことを確認する
5. `go build ./...` でビルドエラーが出ないことを確認する
6. `go test ./...` で失敗しないかを確認する(必要であればテストケースを追加してから行う)

## 説明

- 冷笑のパターンを追加する場合
  - [internal/patterns/patterns.go](./internal/patterns/patterns.go) の `Definitions` スライスに `Definition` を1つ追加してください(`ID` / `Strict` / `Source` / `Samples`)
    - `Samples` に書いたサンプル文字列は自動でテスト化されます([dowa_test.go](./dowa_test.go)の`TestPatternsMatchSamples`)。個別にテストコードを書く必要はありません
  - 可能であれば **ひらがな** , **カタカナ** , **半角カナ** も同様に含めてください(既存の`mora*`定数が使えないか確認してください)
  - 追加するパターンが明らかに冷笑な場合は `Strict: true` に、冷笑か怪しい場合や文脈によっては冷笑ではないパターンは `Strict: false` (relaxedモードでのみ検知)にしてください

- 機能の追加をする場合
  - PRの説明欄に機能についての説明を記載してください
  - テストケースは作成できる場合は作成してください

> [!Important]
>
> 正規表現はGoの`regexp`パッケージ(RE2エンジン)でコンパイル可能である必要があります。後読み・先読みや後方参照はRE2でサポートされていないため使用できません

## リリース

Goのモジュールはnpmのような「公開」の手順がなく、[SemVer](https://semver.org/lang/ja/)形式のgitタグが存在すれば`go get github.com/otnc/dowa-go@vX.Y.Z`として利用可能になる。

リリースはGitHubの Actions タブから [Release workflow](./.github/workflows/release.yml) を手動実行(workflow_dispatch)して行う。

- `version`に`patch`/`minor`/`major`のいずれかを指定すると最新タグから自動算出、`v1.2.3`のように明示的なバージョンを指定するとそのまま使われる
- 実行するとタグの作成・push、GitHub Releaseの作成、[pkg.go.dev](https://pkg.go.dev/github.com/otnc/dowa-go) へのインデックス通知までを自動で行う

メジャーバージョンが2以上になる場合は、go.modのモジュールパスに`/v2`のようなサフィックスを付ける必要がある([Goのモジュールバージョニングの仕様](https://go.dev/doc/modules/major-version)による)。

## 禁止事項

- .github/ を変更すること
- go.mod のmoduleパスを変更すること

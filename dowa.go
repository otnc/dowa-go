// Package dowa は日本語テキスト中の「冷笑」表現("うおw"、"爆笑爆笑"、😅 など)を正規表現ベースで検知する。TypeScript版 https://github.com/otnc/dowa の移植。
//
// パターン定義と正規表現の組み立ては internal/patterns パッケージが担い、このパッケージは利用者向けのAPIだけを公開する薄いラッパーになっている。
package dowa

import "github.com/otnc/dowa-go/internal/patterns"

// Options は検知の挙動を調整するオプション。ゼロ値(Options{})が既定の挙動になる。
type Options struct {
	// Relaxed を true にすると検知範囲を拡大する(既定: false)。
	Relaxed bool
}

// resolveOptions は可変長引数で渡された Options を1つに解決する。
// 呼び出し側が省略した場合はゼロ値(Relaxed: false)を使う。
func resolveOptions(opts []Options) Options {
	if len(opts) == 0 {
		return Options{}
	}
	return opts[0]
}

// FindAll はテキスト中の冷笑パターンをすべて検出する。
// マッチしなければ nil を返す。
func FindAll(text string, opts ...Options) []string {
	o := resolveOptions(opts)
	return patterns.MatchAll(text, o.Relaxed)
}

// Contains はテキストに冷笑パターンが含まれるかを判定する。
func Contains(text string, opts ...Options) bool {
	return FindAll(text, opts...) != nil
}

// Match は FindMatches が返すマッチ1件分の詳細。
type Match struct {
	// Text はマッチした文字列。
	Text string
	// Index はマッチ開始位置(バイトオフセット)。
	// TypeScript版はUTF-16コードユニット単位のindexを返すが、Goの文字列はUTF-8バイト列なのでこちらはバイトオフセットになる。
	Index int
	// PatternID はマッチしたパターンのid(Patterns参照)。
	PatternID string
	// Strict はそのパターンがstrictかどうか。
	Strict bool
}

// FindMatches はテキスト中の冷笑パターンを、どのパターンにマッチしたかの詳細付きで検出する。パターンごとに独立して検索するため、複数パターンの一致範囲が重なる場合はそれぞれ個別の結果として返る(FindAllの重複排除された結果とは一致しないことがある)。マッチがなければ nil を返す。
func FindMatches(text string, opts ...Options) []Match {
	o := resolveOptions(opts)
	raw := patterns.FindMatches(text, o.Relaxed)
	if raw == nil {
		return nil
	}
	out := make([]Match, len(raw))
	for i, m := range raw {
		out[i] = Match{Text: m.Text, Index: m.Index, PatternID: m.PatternID, Strict: m.Strict}
	}
	return out
}

// PatternInfo は1つの冷笑検知パターンの概要(id・strict・サンプル)。
// パターンを紹介・デバッグしたい場合に参照する。実際の正規表現ソースは非公開(internal/patterns)。
type PatternInfo struct {
	ID      string
	Strict  bool
	Samples []string
}

// Patterns は検知に使われている全パターンの概要を返す。
func Patterns() []PatternInfo {
	defs := patterns.Definitions
	out := make([]PatternInfo, len(defs))
	for i, d := range defs {
		out[i] = PatternInfo{ID: d.ID, Strict: d.Strict, Samples: d.Samples}
	}
	return out
}

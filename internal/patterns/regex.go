package patterns

import (
	"regexp"
	"sort"
	"strings"
)

// compiledPattern は Definition を relaxed / strict それぞれのモード向けにあらかじめコンパイルしたもの。呼び出しのたびに regexp.Compile するのは無駄なので、パッケージ初期化時に1度だけ行う。
type compiledPattern struct {
	def Definition
	// reStrict は strict モード用にコンパイルした正規表現。
	// def.Strict が false のパターンは strict モードでは使われないため nil。
	reStrict *regexp.Regexp
	// reRelaxed は relaxed モード用にコンパイルした正規表現。
	reRelaxed *regexp.Regexp
}

var (
	// regexStrict / regexRelaxed は全strictパターン(・relaxedなら全パターン)を "(?:pattern1)|(?:pattern2)|..." の形で1つにまとめた正規表現。FindAll / Contains で使う。
	regexStrict  *regexp.Regexp
	regexRelaxed *regexp.Regexp

	// compiledPatterns はパターンごとの詳細位置を調べる FindMatches で使う。
	compiledPatterns []compiledPattern
)

func init() {
	regexStrict = regexp.MustCompile(buildSource(false))
	regexRelaxed = regexp.MustCompile(buildSource(true))

	compiledPatterns = make([]compiledPattern, len(Definitions))
	for i, p := range Definitions {
		cp := compiledPattern{def: p}
		if p.Strict {
			cp.reStrict = regexp.MustCompile(p.Source(false))
		}
		cp.reRelaxed = regexp.MustCompile(p.Source(true))
		compiledPatterns[i] = cp
	}
}

// buildSource は relaxedMode に応じて対象パターンを選び、1つの正規表現ソースにまとめる。strict なパターンを常に先に並べる。正規表現の選択肢(|)は最初にマッチしたものを採用するため、順序が逆だと "爆笑" のような relaxed専用の短い一致が "(?:爆笑){2,}" より先に取られてしまう。
func buildSource(relaxedMode bool) string {
	var relevant []Definition
	for _, p := range Definitions {
		if p.Strict {
			relevant = append(relevant, p)
		}
	}
	if relaxedMode {
		for _, p := range Definitions {
			if !p.Strict {
				relevant = append(relevant, p)
			}
		}
	}

	parts := make([]string, len(relevant))
	for i, p := range relevant {
		parts[i] = "(?:" + p.Source(relaxedMode) + ")"
	}
	return strings.Join(parts, "|")
}

// MatchAll はテキストにマッチした冷笑パターンの配列を返す(なければnil)。
func MatchAll(text string, relaxed bool) []string {
	re := regexStrict
	if relaxed {
		re = regexRelaxed
	}
	m := re.FindAllString(text, -1)
	if len(m) == 0 {
		return nil
	}
	return m
}

// Match は FindMatches が返すマッチ1件分の詳細。
type Match struct {
	Text      string
	Index     int
	PatternID string
	Strict    bool
}

// FindMatches はテキスト中の冷笑パターンを、どのパターンにマッチしたかの詳細付きで検出する。パターンごとに独立して検索するため、複数パターンの一致範囲が重なる場合はそれぞれ個別の結果として返る(MatchAllの重複排除された結果とは一致しないことがある)。マッチがなければ nil を返す。
func FindMatches(text string, relaxed bool) []Match {
	var results []Match
	for _, cp := range compiledPatterns {
		if !cp.def.Strict && !relaxed {
			continue
		}
		re := cp.reRelaxed
		if !relaxed {
			re = cp.reStrict
		}
		for _, loc := range re.FindAllStringIndex(text, -1) {
			results = append(results, Match{
				Text:      text[loc[0]:loc[1]],
				Index:     loc[0],
				PatternID: cp.def.ID,
				Strict:    cp.def.Strict,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Index < results[j].Index })
	if len(results) == 0 {
		return nil
	}
	return results
}

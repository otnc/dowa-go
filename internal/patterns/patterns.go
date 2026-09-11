package patterns

import "strings"

// Definition は1つの冷笑検知パターンを表す。
//
// Source は relaxed モードかどうかで内容が変わる場合があるため、固定文字列ではなく関数として保持する(TypeScript版の `string | ((relaxed: boolean) => string)` に相当する)。
type Definition struct {
	ID string
	// Strict が true なら常に検知対象、false なら relaxed モードのときだけ検知対象。
	Strict bool
	// Source は正規表現ソース(フラグなし)を relaxed モードに応じて返す。
	Source func(relaxed bool) string
	// Samples はこのパターンにマッチするはずのサンプル文字列(テストに使用)。
	Samples []string
}

// モーラの異表記(ひらがな/カタカナ/半角カナ)。濁点合成が要るもの(が/ざ/だ行等)は半角カナの濁点が2文字になるため文字クラスではなく(?:...)の選択構造にしている。
const (
	moraA   = "[あぁアァｱｧ]"
	moraI   = "[いぃイィｲｨ]"
	moraU   = "[うぅウゥｳｩ]"
	moraE   = "[えぇエェｴｪ]"
	moraO   = "[おぉオォｵｫ]"
	moraKA  = "[かカｶ]"
	moraKI  = "[きキｷ]"
	moraKU  = "[くクｸ]"
	moraKE  = "[けケｹ]"
	moraKO  = "[こコｺ]"
	moraGU  = "(?:ぐ|グ|ｸﾞ)"
	moraGO  = "(?:ご|ゴ|ｺﾞ)"
	moraSA  = "[さサｻ]"
	moraSI  = "[しシｼ]"
	moraSU  = "[すスｽ]"
	moraSE  = "[せセｾ]"
	moraSO  = "[そソｿ]"
	moraTA  = "[たタﾀ]"
	moraTI  = "[ちチﾁ]"
	moraXTU = "[っッｯ]"
	moraTE  = "[てテﾃ]"
	moraDE  = "(?:で|デ|ﾃﾞ)"
	moraDO  = "(?:ど|ド|ﾄﾞ)"
	moraNA  = "[なナﾅ]"
	moraNO  = "[のノﾉ]"
	moraBA  = "(?:ば|バ|ﾊﾞ)"
	moraBO  = "(?:ぼ|ボ|ﾎﾞ)"
	moraMU  = "[むムﾑ]"
	moraMO  = "[もモﾓ]"
	moraYA  = "[やヤﾔ]"
	moraXYO = "[ょョｮ]"
	moraYO  = "[よヨﾖ]"
	moraRI  = "[りリﾘ]"
	moraRO  = "[ろロﾛ]"
	moraWA  = "[わゎワヮﾜ]"
	moraN   = "[んンﾝ]"
)

// ！/？: 語尾なしだとstrictは2文字以上、relaxedは1文字以上で許容。
const lowPunct = "[！!？?]"

// ‼️❗❓⁉️: 語尾なしでも常に1文字から許容。
const highPunct = `(?:[❗❓]|‼\x{FE0F}?|⁉\x{FE0F}?)`

const anyPunct = "(?:" + lowPunct + "|" + highPunct + ")"

const laugh = `(?:[wｗ]+|(?:(?:爆笑)|笑)+|[（(]笑[）)])`

// stretchTail は語幹の後ろに続く伸ばし棒・促音の繰り返しを許容する部分。
const stretchTail = "[-ｰー～っッｯ]*"

// 手・指のジェスチャー系絵文字に付く肌の色modifier(任意)。
const skinTone = `[\x{1F3FB}-\x{1F3FF}]?`

// stemSuffix は「語幹+伸ばし棒+記号+語尾」の語尾部分を組み立てる。
// 記号だけでも条件を満たせば語尾(w/笑など)なしで許容する。
func stemSuffix(relaxed bool) string {
	lowMin := "2"
	if relaxed {
		lowMin = "1"
	}
	return stretchTail + "(?:" +
		anyPunct + "*" + laugh +
		"|" + lowPunct + "{" + lowMin + ",}" +
		"|" + highPunct + "+" +
		")"
}

// stem は語幹+stemSuffixを組み立てる Definition.Source を返す。
func stem(parts ...string) func(bool) string {
	joined := strings.Join(parts, "")
	return func(relaxed bool) string {
		return joined + stemSuffix(relaxed)
	}
}

// bare は語尾(w/笑など)を要求しない語幹単体(伸ばし棒の繰り返しのみ許容)を返す。
func bare(parts ...string) func(bool) string {
	joined := strings.Join(parts, "") + stretchTail
	return func(bool) string { return joined }
}

// literal は relaxed モードに関わらず内容が変わらない固定の正規表現ソースを返す。
func literal(source string) func(bool) string {
	return func(bool) string { return source }
}

// yan は「(です)やん」で終わる冷笑フレーズ共通のビルダー。
func yan(word string) func(bool) string {
	return func(relaxed bool) string {
		return word + "(?:" + moraDE + moraSU + ")?" + moraYA + moraN + stemSuffix(relaxed)
	}
}

// Definitions は検知に使われている全パターンの定義。
// パターンを紹介・デバッグしたい場合に参照する。
var Definitions = []Definition{
	// --- 絵文字 (strict) ---
	{ID: "emoji-sweat-smile", Strict: true, Source: literal("😅"), Samples: []string{"それはさすがに草😅"}},
	{ID: "emoji-rofl", Strict: true, Source: literal("🤣"), Samples: []string{"それ何回同じネタやってんの🤣🤣"}},
	{ID: "emoji-double-exclamation", Strict: true, Source: literal(`‼\x{FE0F}?`), Samples: []string{"そんな‼️"}},
	{ID: "emoji-bang", Strict: true, Source: literal(`❗\x{FE0F}?`), Samples: []string{"いいね❗"}},
	{ID: "emoji-question", Strict: true, Source: literal(`❓\x{FE0F}?`), Samples: []string{"は❓"}},
	{ID: "emoji-interrobang", Strict: true, Source: literal(`⁉\x{FE0F}?`), Samples: []string{"まじで言ってる⁉️"}},
	{ID: "emoji-eye-roll", Strict: true, Source: literal("🙄"), Samples: []string{"はいはい🙄"}},
	{ID: "emoji-smirk", Strict: true, Source: literal("😏"), Samples: []string{"それな😏"}},
	{ID: "emoji-clown", Strict: true, Source: literal("🤡"), Samples: []string{"ふっ🤡"}},
	{
		ID: "emoji-grin-fist", Strict: true,
		// 😁 + (✊|👊) の組み合わせ(順不同、✊👊は肌の色modifier許容)。
		Source:  literal(`(?:😁(?:✊|👊)` + skinTone + `|(?:✊|👊)` + skinTone + `😁)`),
		Samples: []string{"また学校来いよ😁✊", "待ってるからな👊😁"},
	},
	{ID: "emoji-grin", Strict: true, Source: literal("😁"), Samples: []string{"余裕😁"}},
	{ID: "emoji-fist-raised", Strict: true, Source: literal("✊" + skinTone), Samples: []string{"かかってこいよ✊"}},
	{ID: "emoji-fist-oncoming", Strict: true, Source: literal("👊" + skinTone), Samples: []string{"やる気👊"}},
	{ID: "emoji-point-at-viewer", Strict: true, Source: literal("🫵" + skinTone), Samples: []string{"それお前のことな🫵"}},
	{ID: "emoji-ok-hand", Strict: true, Source: literal("👌" + skinTone), Samples: []string{"了解👌"}},
	{ID: "emoji-joy", Strict: true, Source: literal("😂"), Samples: []string{"おけ😂"}},
	{ID: "emoji-grinning", Strict: true, Source: literal("😃"), Samples: []string{"すごいですね😃"}},
	{ID: "emoji-tongue", Strict: true, Source: literal("😛"), Samples: []string{"残念😛"}},
	{ID: "emoji-wink-tongue", Strict: true, Source: literal("😜"), Samples: []string{"バレたか😜"}},
	{ID: "emoji-open-mouth", Strict: true, Source: literal("😮"), Samples: []string{"まさかそれ本気で言ってる😮"}},
	{ID: "emoji-astonished", Strict: true, Source: literal("😲"), Samples: []string{"それはさすがに草😲"}},
	{ID: "emoji-frowning-open-mouth", Strict: true, Source: literal("😦"), Samples: []string{"それは無理があるでしょ😦"}},
	{ID: "emoji-anguished", Strict: true, Source: literal("😧"), Samples: []string{"見てて痛々しい😧"}},
	{ID: "emoji-fearful", Strict: true, Source: literal("😨"), Samples: []string{"それはさすがにやばい😨"}},
	{ID: "emoji-weary", Strict: true, Source: literal("😩"), Samples: []string{"もう見てられない😩"}},
	{ID: "emoji-zany", Strict: true, Source: literal("🤪"), Samples: []string{"それマジで言ってる🤪"}},
	{ID: "emoji-hot-face", Strict: true, Source: literal("🥵"), Samples: []string{"必死すぎん🥵"}},
	{ID: "emoji-anxious-sweat", Strict: true, Source: literal("😰"), Samples: []string{"それはさすがに焦るわ😰"}},
	{ID: "emoji-scream", Strict: true, Source: literal("😱"), Samples: []string{"うそでしょ😱"}},
	{ID: "emoji-nerd", Strict: true, Source: literal("🤓"), Samples: []string{"詳しいっすね🤓"}},
	{ID: "emoji-hand-over-mouth", Strict: true, Source: literal("🤭"), Samples: []string{"ぷっ🤭"}},
	{ID: "emoji-eyes-hand-over-mouth", Strict: true, Source: literal("🫢"), Samples: []string{"まじで言ってるの🫢"}},
	{ID: "emoji-lying", Strict: true, Source: literal("🤥"), Samples: []string{"よく言うわ🤥"}},
	{ID: "emoji-cowboy", Strict: true, Source: literal("🤠"), Samples: []string{"余裕じゃん🤠"}},
	{ID: "emoji-point-up", Strict: true, Source: literal("👆" + skinTone), Samples: []string{"それそれ👆"}},
	{ID: "emoji-point-down", Strict: true, Source: literal("👇" + skinTone), Samples: []string{"こいつを見て👇"}},

	// --- 絵文字 (relaxedのみ: 単体だと冷笑と断定しづらいもの) ---
	{ID: "emoji-sweat-drop", Strict: false, Source: literal("💦"), Samples: []string{"それはさすがに草だわ💦"}},
	{ID: "emoji-expressionless", Strict: false, Source: literal("😑"), Samples: []string{"……😑"}},
	{ID: "emoji-upside-down", Strict: false, Source: literal("🙃"), Samples: []string{"はいはい、そうですね🙃"}},
	{ID: "emoji-skull", Strict: false, Source: literal("💀"), Samples: []string{"その理論はさすがに無理💀"}},
	{ID: "emoji-melting", Strict: false, Source: literal("🫠"), Samples: []string{"見てるだけでしんど🫠"}},
	{ID: "emoji-ok", Strict: false, Source: literal("🆗"), Samples: []string{"🆗"}},
	{ID: "emoji-pleading", Strict: false, Source: literal("🥺"), Samples: []string{"許して🥺"}},
	{ID: "emoji-grin-squint", Strict: false, Source: literal("😆"), Samples: []string{"それは草😆"}},
	{ID: "emoji-thinking", Strict: false, Source: literal("🤔"), Samples: []string{"それ本気で言ってる?🤔"}},
	{ID: "emoji-salute", Strict: false, Source: literal("🫡"), Samples: []string{"了解しました🫡"}},
	{ID: "emoji-partying", Strict: false, Source: literal("🥳"), Samples: []string{"やったぜ🥳"}},
	{ID: "emoji-squint-tongue", Strict: false, Source: literal("😝"), Samples: []string{"ばれちゃった😝"}},
	{ID: "emoji-monocle", Strict: false, Source: literal("🧐"), Samples: []string{"それで?🧐"}},
	{ID: "emoji-shushing", Strict: false, Source: literal("🤫"), Samples: []string{"へぇ🤫"}},
	{ID: "emoji-shaking", Strict: false, Source: literal("🫨"), Samples: []string{"それやば🫨"}},
	{ID: "emoji-holding-tears", Strict: false, Source: literal("🥹"), Samples: []string{"感動した🥹"}},

	// --- 語幹 + w/笑/爆笑/(笑) (strict) ---
	{ID: "stem-kichi", Strict: true, Source: stem(moraKI, moraTI), Samples: []string{"きちーｗ"}},
	{ID: "stem-ou", Strict: true, Source: stem(moraO, moraU), Samples: []string{"お、おうｗ"}},
	{ID: "stem-uo", Strict: true, Source: stem(moraU, moraO), Samples: []string{"うおw"}},
	{ID: "stem-oke", Strict: true, Source: stem(moraO, moraKE), Samples: []string{"おけ（笑）"}},
	{ID: "stem-dowa", Strict: true, Source: stem(moraDO, moraWA), Samples: []string{"どわーwww"}},
	{ID: "stem-uwa", Strict: true, Source: stem(moraU, moraWA), Samples: []string{"うわーw"}},
	{ID: "stem-yaba", Strict: true, Source: stem(moraYA, moraBA), Samples: []string{"それヤバ笑"}},
	{ID: "stem-samu", Strict: true, Source: stem(moraSA, moraMU), Samples: []string{"そのノリさむw"}},
	{ID: "stem-ita", Strict: true, Source: stem(moraA+"?", moraI, moraTA+"+"), Samples: []string{"アイタタタタw"}},
	{ID: "stem-kimo", Strict: true, Source: stem(moraKI, moraMO), Samples: []string{"きもwドン引きだわ"}},
	{ID: "stem-kita", Strict: true, Source: stem(moraKI, moraTA), Samples: []string{"キター！！！！ｗ"}},

	// --- 語幹 + w/笑/爆笑/(笑) (relaxedのみ: 単体では冷笑以外の文脈でも頻出するため) ---
	{ID: "stem-iya", Strict: false, Source: stem(moraI, moraYA), Samples: []string{"いやwそれは草"}},

	// --- 語幹単体 (relaxedのみ: 語尾のw/笑がなくても検知する) ---
	{ID: "bare-uo", Strict: false, Source: bare(moraU, moraO), Samples: []string{"う、うお、しか言えなくなってて草"}},
	{ID: "bare-dowa", Strict: false, Source: bare(moraDO, moraWA), Samples: []string{"どわ…しか反応できてなくて草"}},
	{ID: "bare-bakushou", Strict: false, Source: literal("爆笑"), Samples: []string{"その返し思わず爆笑してしまった"}},
	{ID: "bare-reishou", Strict: false, Source: literal("冷笑"), Samples: []string{"これが世に言う冷笑ってやつか"}},

	// --- 繰り返しパターン (strict) ---
	{ID: "repeat-bakushou", Strict: true, Source: literal("(?:爆笑){2,}"), Samples: []string{"その言い訳マジで爆笑爆笑"}},
	{ID: "repeat-reishou", Strict: true, Source: literal("(?:冷笑){2,}"), Samples: []string{"これぞ正統派の冷笑冷笑という感じ"}},
	{ID: "repeat-warai", Strict: true, Source: literal("(?:笑){2,}"), Samples: []string{"それは草生えるわ笑笑"}},
	{ID: "paren-warai", Strict: true, Source: literal("[（(]笑[）)]"), Samples: []string{"はいはい、すごいですね（笑）"}},

	// --- フレーズ系 (strict) ---
	{ID: "phrase-omoro", Strict: true, Source: stem(moraO, moraMO, moraRO, moraI+"?"), Samples: []string{"おもろwww"}},
	{ID: "phrase-kakke", Strict: true, Source: stem(moraKA, moraXTU, moraKE), Samples: []string{"かっけーwイキっててウケる"}},
	{ID: "phrase-kakkoyo", Strict: true, Source: stem(moraKA, moraXTU, moraKO, moraYO), Samples: []string{"かっこよwナルシストかよ"}},
	{ID: "phrase-egui", Strict: true, Source: stem(moraE, moraGU), Samples: []string{"その自己評価えぐー！笑"}},
	{ID: "phrase-doshita", Strict: true, Source: stem(moraDO, moraSI, moraTA, moraN+"?"), Samples: []string{"ど、どした？笑 急に早口になって"}},
	// 「(です)やん」+ 語尾の冷笑的な相槌。
	{ID: "phrase-yan", Strict: true, Source: yan(""), Samples: []string{"めっちゃ必死やんw", "冗談やんwノリ悪いなあ", "冗談ですやんw"}},
	{ID: "phrase-sonna-nori", Strict: true, Source: stem(moraSO, moraU, moraI, moraU, moraNO, moraRI, `[…\.・･]*`), Samples: []string{"あぁ、そういうノリ...w理解した"}},
	// 冷笑チェーンネタで使われる定型フレーズ。
	{ID: "phrase-yoseyai", Strict: true, Source: stem(moraYO, moraSE, moraYA, moraI), Samples: []string{"よせやいw"}},
	{ID: "phrase-atabouyo", Strict: true, Source: stem(moraA, moraTA, moraBO, moraU, moraYO), Samples: []string{"あたぼうよw"}},
	{ID: "phrase-teyandei", Strict: true, Source: stem(moraTE, moraYA, moraN, moraDE, moraI), Samples: []string{"てやんでいw"}},

	// --- フレーズ系 (relaxedのみ: 単体では冷笑以外の文脈でも頻出するため) ---
	{ID: "phrase-cho", Strict: false, Source: stem(moraTI, moraXYO), Samples: []string{"ちょwそれは草すぎる"}},
	{ID: "phrase-mattaku", Strict: false, Source: stem(moraXTU, moraTA, moraKU), Samples: []string{"ったくwしょうがないやつだな"}},
	{ID: "phrase-omoroi", Strict: false, Source: stem(moraO, moraMO, moraRO, moraI, moraNA, moraA+"?"), Samples: []string{"おもろいなあwキミw"}},
	{ID: "phrase-sugoi", Strict: false, Source: stem(moraSU, moraGO, moraI, moraNA, moraA+"?"), Samples: []string{"すごいなあwwキミ見損なったわ"}},
}

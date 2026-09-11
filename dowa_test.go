package dowa

import (
	"reflect"
	"testing"
)

// 各パターンの Samples が、そのパターン自身で検知できることを確認する。
// strict パターンは strict モード(既定)で、relaxed専用パターンは relaxed モードで検証する。
func TestPatternsMatchSamples(t *testing.T) {
	for _, p := range Patterns() {
		p := p
		t.Run(p.ID, func(t *testing.T) {
			opts := Options{Relaxed: !p.Strict}
			for _, sample := range p.Samples {
				if !Contains(sample, opts) {
					t.Errorf("pattern %q should match sample %q (relaxed=%v)", p.ID, sample, opts.Relaxed)
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		name string
		text string
		want bool
	}{
		{"plain text", "それは普通の会話です", false},
		{"uow with laugh", "うおw", true},
		{"repeated bakushou", "爆笑爆笑", true},
		{"single bakushou strict", "爆笑してしまった", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Contains(c.text); got != c.want {
				t.Errorf("Contains(%q) = %v, want %v", c.text, got, c.want)
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	got := FindAll("←うおw、爆笑爆笑")
	want := []string{"うおw", "爆笑爆笑"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAll() = %#v, want %#v", got, want)
	}

	if got := FindAll("普通のテキスト"); got != nil {
		t.Errorf("FindAll() = %#v, want nil", got)
	}
}

func TestContainsRelaxed(t *testing.T) {
	text := "その理論はさすがに無理💀"
	if Contains(text) {
		t.Errorf("Contains(%q) should be false in strict mode", text)
	}
	if !Contains(text, Options{Relaxed: true}) {
		t.Errorf("Contains(%q, Relaxed: true) should be true", text)
	}
}

func TestFindMatches(t *testing.T) {
	matches := FindMatches("うおw、爆笑爆笑")
	if len(matches) != 2 {
		t.Fatalf("len(matches) = %d, want 2", len(matches))
	}
	if matches[0].PatternID != "stem-uo" || matches[0].Text != "うおw" {
		t.Errorf("matches[0] = %#v", matches[0])
	}
	if matches[1].PatternID != "repeat-bakushou" || matches[1].Text != "爆笑爆笑" {
		t.Errorf("matches[1] = %#v", matches[1])
	}

	if got := FindMatches("普通のテキスト"); got != nil {
		t.Errorf("FindMatches() = %#v, want nil", got)
	}
}

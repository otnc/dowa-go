package dowa_test

import (
	"fmt"

	"github.com/otnc/dowa-go"
)

func ExampleContains() {
	fmt.Println(dowa.Contains("うおw"))
	fmt.Println(dowa.Contains("どわー", dowa.Options{Relaxed: true}))
	// Output:
	// true
	// true
}

func ExampleFindAll() {
	all := dowa.FindAll("←うおw、爆笑爆笑")
	fmt.Println(all)
	// Output:
	// [うおw 爆笑爆笑]
}

func ExampleFindMatches() {
	matches := dowa.FindMatches("うおw、爆笑爆笑")
	for _, m := range matches {
		fmt.Printf("%s (pattern=%s, strict=%v)\n", m.Text, m.PatternID, m.Strict)
	}
	// Output:
	// うおw (pattern=stem-uo, strict=true)
	// 爆笑爆笑 (pattern=repeat-bakushou, strict=true)
}

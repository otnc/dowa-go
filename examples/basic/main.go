// Command basic は github.com/otnc/dowa-go の基本的な使い方を示すサンプル。
//
// 実行するには次のようにする:
//
//	go run ./examples/basic
package main

import (
	"fmt"

	"github.com/otnc/dowa-go"
)

func main() {
	fmt.Println(dowa.Contains("うおw"))
	fmt.Println(dowa.Contains("どわー", dowa.Options{Relaxed: true}))

	all := dowa.FindAll("←うおw、爆笑爆笑")
	fmt.Println(all)

	for _, m := range dowa.FindMatches("うおw、爆笑爆笑") {
		fmt.Printf("%s (pattern=%s, strict=%v)\n", m.Text, m.PatternID, m.Strict)
	}
}

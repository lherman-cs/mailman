//+build generate

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(http.Dir("assets"), vfsgen.Options{
		Filename:  "assets_prod.go",
		BuildTags: "prod",
	})
	if err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"../../../linkparser"
	"fmt"
	"strings"
)

var exampleHtml = `
<html>
    <body>
        <a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
    </body>
</html>
`

func main() {
	r := strings.NewReader(exampleHtml)
	links, err := linkparser.Parse(r)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", links)
}

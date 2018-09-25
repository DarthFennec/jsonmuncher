package main

import (
	"fmt"
	"html/template"
	"os"
)

var header = `{
  "apiVersion": "v1",
  "generated": "2018-04-11T16:56:56.656249201Z",
  "entries": {
`

var footer = `
    ]
  }
}`

var chart = `    "dummy-chart-{{.Num}}": [
`

var release = `      {
        "name": "dummy-chart-{{.Num}}",
        "home": "https://example.com",
        "sources": [
          "https://example.com",
          "https://example.com"
        ],
        "version": "1.2.{{.Num2}}",
        "description": "Example description",
        "keywords": [
          "A",
          "B"
        ],
        "maintainers": [
          {
            "name": "Bar",
            "email": "bar@example.com"
          }
        ],
        "icon": "https://example.com/foo.png",
        "urls": [
          "https://example.com"
        ],
        "created": "2017-07-06T01:33:50.952906435Z",
        "digest": "249e27501dbfe1bd93d4039b04440f0ff19c707ba720540f391b5aefa3571455"
      }`

type wrapper struct {
	Num  int
	Num2 int
}

func main() {
	fmt.Println("Generating fixture_huge.json for testing")
	f, err := os.Create("./fixture_huge.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	_, err = f.WriteString(header)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	charttmpl, err := template.New("chart").Parse(chart)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reltmpl, err := template.New("release").Parse(release)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var w wrapper
	for i := 0; i < 100; i++ {
		w = wrapper{Num: i}
		err = charttmpl.Execute(f, w)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for j := 0; j < 5000; j++ {
			w.Num2 = j
			err = reltmpl.Execute(f, w)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if j < 4999 {
				_, err = f.WriteString(",\n")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
		if i < 99 {
			_, err = f.WriteString("\n    ],\n")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
	_, err = f.WriteString(footer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done generating testing file")
}

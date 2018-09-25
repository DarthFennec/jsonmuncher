package main

import (
	"os"
	"fmt"
	"bufio"
	"regexp"
	"strings"
)

type inputLine struct {
	project string
	size string
	time string
	space string
	alloc string
}

type outputLine struct {
	project string
	small string
	medium string
	large string
	huge string
}

var projectmap = map[string]string{
	"JsonMuncher"					: "[`github.com/darthfennec/jsonmuncher`][]",
	"JsonParser"					: "[`github.com/buger/jsonparser`][]",
	"EncodingJsonStruct"			: "[`encoding/json`][] (struct)",
	"EncodingJsonInterface"			: "[`encoding/json`][] (interface)",
	"EncodingJsonStreamStruct"		: "[`encoding/json`][] (struct streaming)",
	"EncodingJsonStreamInterface"	: "[`encoding/json`][] (interface streaming)",
	"Jstream"						: "[`github.com/bcicen/jstream`][]",
	"Gojay"							: "[`github.com/francoispqt/gojay`][]",
	"JsonIterator"					: "[`github.com/json-iterator/go`][]",
	"Gabs"							: "[`github.com/jeffail/gabs`][]",
	"GoSimpleJson"					: "[`github.com/bitly/go-simplejson`][]",
	"FFJson"						: "[`github.com/pquerna/ffjson`][]",
	"Jason"							: "[`github.com/antonholmquist/jason`][]",
	"Ujson"							: "[`github.com/mreiferson/go-ujson`][]",
	"Djson"							: "[`github.com/a8m/djson`][]",
	"Ugorji"						: "[`github.com/ugorji/go/codec`][]",
	"EasyJson"						: "[`github.com/mailru/easyjson`][]",
}

var projectorder = [...]string{
	"[`github.com/antonholmquist/jason`][]",
	"[`github.com/bcicen/jstream`][]",
	"[`github.com/bitly/go-simplejson`][]",
	"[`github.com/ugorji/go/codec`][]",
	"[`github.com/jeffail/gabs`][]",
	"[`github.com/mreiferson/go-ujson`][]",
	"[`github.com/json-iterator/go`][]",
	"[`github.com/a8m/djson`][]",
	"[`encoding/json`][] (interface streaming)",
	"[`encoding/json`][] (struct streaming)",
	"[`encoding/json`][] (interface)",
	"[`encoding/json`][] (struct)",
	"[`github.com/francoispqt/gojay`][]",
	"[`github.com/pquerna/ffjson`][]",
	"[`github.com/mailru/easyjson`][]",
	"[`github.com/buger/jsonparser`][]",
	"[`github.com/darthfennec/jsonmuncher`][]",
}

var rex = regexp.MustCompile("^Benchmark(\\w+)(Huge|Large|Medium|Small)-8\\s+\\d+\\s+(\\d+) ns/op\\s+(\\d+) B/op\\s+(\\d+) allocs/op$")

func main() {
	lines := make([]inputLine, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txt := scanner.Text()
		l := rex.FindStringSubmatch(txt)
		if l == nil {
			fmt.Fprintln(os.Stderr, "Warning: cannot match line:", txt)
			continue
		}
		line := inputLine{projectmap[l[1]], l[2], l[3], l[4], l[5]}
		lines = append(lines, line)
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
		os.Exit(1)
	}
	times := make([]outputLine, 0)
	spaces := make([]outputLine, 0)
	allocs := make([]outputLine, 0)
	for _, ln := range projectorder {
		time := outputLine{ln, "", "", "", ""}
		space := outputLine{ln, "", "", "", ""}
		alloc := outputLine{ln, "", "", "", ""}
		for _, l := range lines {
			if l.project != ln {
				continue
			}
			switch l.size {
			case "Small":
				time.small = l.time
				space.small = l.space
				alloc.small = l.alloc
			case "Medium":
				time.medium = l.time
				space.medium = l.space
				alloc.medium = l.alloc
			case "Large":
				time.large = l.time
				space.large = l.space
				alloc.large = l.alloc
			case "Huge":
				time.huge = l.time
				space.huge = l.space
				alloc.huge = l.alloc
			}
		}
		times = append(times, time)
		spaces = append(spaces, space)
		allocs = append(allocs, alloc)
	}
	drawTable("Speed", "ns/op", times)
	drawTable("Memory", "B/op", spaces)
	drawTable("Allocations", "allocs/op", allocs)
}

func drawTable(title string, unit string, lines []outputLine) {
	libmax := 0
	smlmax := 0
	medmax := 0
	lrgmax := 0
	hgemax := 0
	for _, line := range lines {
		extra := 0
		if line.project == "[`github.com/darthfennec/jsonmuncher`][]" {
			extra = 4
		}
		if len(line.project) > libmax {
			libmax = len(line.project)
		}
		if len(line.small) + extra > smlmax {
			smlmax = len(line.small) + extra
		}
		if len(line.medium) + extra > medmax {
			medmax = len(line.medium) + extra
		}
		if len(line.large) + extra > lrgmax {
			lrgmax = len(line.large) + extra
		}
		if len(line.huge) + extra > hgemax {
			hgemax = len(line.huge) + extra
		}
	}
	if libmax < 7 {
		libmax = 7
	}
	smlmax += 1 + len(unit)
	if smlmax < 10 {
		smlmax = 10
	}
	medmax += 1 + len(unit)
	if medmax < 11 {
		medmax = 11
	}
	lrgmax += 1 + len(unit)
	if lrgmax < 10 {
		lrgmax = 10
	}
	hgemax += 1 + len(unit)
	if hgemax < 9 {
		hgemax = 9
	}
	fmt.Printf("### %s\n\n", title)
	fmt.Printf("%-*s | %-*s | %-*s | %-*s | %s\n",
		libmax, "Library", smlmax, "Small JSON", medmax, "Medium JSON",
		lrgmax, "Large JSON", "Huge JSON")
	fmt.Printf(":%s|-%s:|-%s:|-%s:|%s:\n",
		strings.Repeat("-", libmax), strings.Repeat("-", smlmax),
		strings.Repeat("-", medmax), strings.Repeat("-", lrgmax),
		strings.Repeat("-", hgemax))
	for _, line := range lines {
		lib := line.project
		sml := line.small + " " + unit
		med := line.medium + " " + unit
		lrg := line.large + " " + unit
		hge := line.huge + " " + unit
		if lib == "[`github.com/darthfennec/jsonmuncher`][]" {
			sml = "**" + sml + "**"
			med = "**" + med + "**"
			lrg = "**" + lrg + "**"
			hge = "**" + hge + "**"
		}
		fmt.Printf("%-*s | %-*s | %-*s | %-*s | %s\n",
			libmax, lib, smlmax, sml, medmax, med, lrgmax, lrg, hge)
	}
	fmt.Printf("\n")
}

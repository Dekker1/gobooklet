package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

const latexTemplate = `\documentclass[a4paper]{article}
\usepackage[pdftex]{color,graphicx,epsfig}
\usepackage[final]{pdfpages}
\begin{document}
\includepdf[pages=-,nup=1x2,landscape,signature={{.Pages}}]{{printf "{"}}{{.FileName}}{{printf "}"}}
\end{document}
`

type info struct {
	Pages    int
	FileName string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getNrPages(fileName string) int {
	data, err := ioutil.ReadFile(fileName)
	check(err)
	lines := strings.Split(strings.Replace(string(data), "\r", "\n", -1), "\n")

	re := regexp.MustCompile("^[^\\/]*\\/Count\\s(\\d+)")
	for _, line := range lines {
		if re.MatchString(line) {
			match := re.FindStringSubmatch(line)
			pages, err := strconv.Atoi(match[1])
			check(err)
			return pages
		}
	}
	return 0
}

func generateLaTeX(inputFile string, outputFile string) {
	t := template.New("latexTemplate")
	t, err := t.Parse(latexTemplate)
	check(err)

	pages := 0
	for pages < getNrPages(inputFile) {
		pages += 4
	}

	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)

	err = t.Execute(w, info{
		Pages: pages,
		FileName: strings.Replace(
			strings.Replace(inputFile, "%", "\\%", -1),
			"_", "\\_", -1),
	})
	w.Flush()
}

func compileAndClean(texFile string) {
	c := exec.Command("latexmk", "-pdf", texFile)
	check(c.Run())

	c = exec.Command("latexmk", "-c", texFile)
	check(c.Run())
}

func main() {
	var outputFile string
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided, use command in the form 'gobooklet [input] ([output])?'")
		os.Exit(-1)
	} else if len(os.Args) < 3 {
		outputFile = os.Args[1]
		outputFile = strings.TrimRight(outputFile, ".pdf")
		outputFile = outputFile + "_booklet"
	} else {
		outputFile = os.Args[2]
		outputFile = strings.TrimRight(outputFile, ".pdf")
	}

	generateLaTeX(os.Args[1], outputFile+".tex")
	compileAndClean(outputFile + ".tex")
	os.Remove(outputFile + ".tex")
}

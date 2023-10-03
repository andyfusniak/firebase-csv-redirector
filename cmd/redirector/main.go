package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"text/template"
)

func main() {
	var filename string
	flag.StringVar(&filename, "f", "", "-f <filename>")
	flag.Parse()

	if filename == "" {
		fmt.Fprintln(os.Stderr, "Usage: redirector -f <filename>")
		os.Exit(1)
	}

	if err := run(filename); err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(2)
	}
}

func run(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		directive, err := firebaseRedirect(record[0], record[1])
		if err != nil {
			return err
		}
		fmt.Printf("%s", directive)
	}

	return nil
}

const t1 = `{
  "source": "{{.Src}}",
  "destination": "{{.Dst}}",
  "type": 301
},
`

var tmpl = template.Must(template.New("t1").Parse(t1))

func firebaseRedirect(src, dst string) (string, error) {
	tp := struct {
		Src string
		Dst string
	}{Src: src, Dst: dst}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, "t1", tp); err != nil {
		return "", err
	}
	return buf.String(), nil
}

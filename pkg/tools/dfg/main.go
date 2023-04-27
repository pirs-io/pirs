package main

import (
	"os"
	"text/template"
)

type DfgData struct {
	AppName            string
	EnvFile            string
	ModuleDependencies []string
}

func main() {

	data := DfgData{
		AppName:            "process",
		EnvFile:            "example.env",
		ModuleDependencies: []string{"commons", "process-storage"},
	}

	t := template.New("dfg")

	files := []string{"pkg/tools/dfg/Dockerfile.gohtml", "pkg/tools/dfg/go.work.gohtml"}
	for _, file := range files {
		f, err := os.ReadFile(file)
		t, err := t.Parse(string(f))
		if err != nil {
			panic(err)
		}
		err = t.Execute(os.Stdout, data)
		if err != nil {
			panic(err)
		}
	}
}

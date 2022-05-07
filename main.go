package main

import (
	"os"
	"strings"

	"github.com/asottile/dockerfile"
)

func main() {

	numOfArgs := len(os.Args)
	dockerFileName := ""
	if numOfArgs > 1 {
		dockerFileName = os.Args[1]
	} else {
		dockerFileName = "./Dockerfile"
	}

	filePrefix := ""
	if numOfArgs > 2 {
		filePrefix = os.Args[2] + "_"
	}

	parsed, err := dockerfile.ParseFile(dockerFileName)

	if err != nil {
		println("unable to parse dockerfile error: " + err.Error())
		return
	}

	var scriptName = ""
	var script = ""
	var entrypoint = ""
	var command = ""
	for _, v := range parsed {
		switch v.Cmd {
		case "FROM":
			{
				if scriptName != "" {
					if !(entrypoint == "" && command == "") {
						script += entrypoint + " " + command + "\n"
					}
					os.WriteFile(scriptName, []byte(script), 0755)
				} else {
					script = ""
				}
				scriptName = filePrefix + v.Value[0] + ".sh"
			}
		case "RUN":
			{
				script += strings.Join(v.Value, " ") + "\n"
			}
		case "CMD":
			{
				command = strings.Join(v.Value, " ")
			}
		case "ENTRYPOINT":
			{
				entrypoint = strings.Join(v.Value, " ")

			}
		case "ARG":
			{
				script += "export " + strings.Join(v.Value, "") + "\n"
			}
		case "LABEL":
			{
				len := len(v.Value)
				for i := 0; i < len/2; i++ {
					script += "# label:  " + v.Value[i*2] + "=" + v.Value[i*2+1] + "\n"
				}
			}
		case "MAINTAINER":
			{
				script += "# Maintainer " + strings.Join(v.Value, "") + "\n"
			}
		case "EXPOSE":
			{
				// no functionality in shell
			}
		case "ENV":
			{
				len := len(v.Value)
				for i := 0; i < len/2; i++ {
					script += "export " + v.Value[i*2] + "=" + v.Value[i*2+1] + "\n"
				}
			}
		case "COPY", "ADD":
			{
				script += "cp " + strings.Join(v.Value, " ") + "\n"
			}
		case "USER":
			{
				if strings.Contains(v.Value[0], ":") {
					script += "su - $(id -un " + v.Value[0] + ")" + "\n"
				} else {
					script += "su " + v.Value[0] + "\n"
				}
			}
		case "WORKDIR":
			{
				script += "mkdir -p " + v.Value[0] + "\n"
				script += "cd " + v.Value[0] + "\n"
			}
		}
		if !(entrypoint == "" && command == "") {
			script += entrypoint + " " + command + "\n"
		}
		os.WriteFile(scriptName, []byte(script), 0755)
	}

}

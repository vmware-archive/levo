/* Copyright (C) 2014 Pivotal Software, Inc.

All rights reserved. This program and the accompanying materials
are made available under the terms of the under the Apache License,
Version 2.0 (the "License”); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

//This type and two methods augment the "-names"" flag
type nameArray []string

func (self *nameArray) String() string {
	return fmt.Sprint(*self)
}
func (self *nameArray) Set(value string) error {
	if len(*self) > 0 {
		return errors.New("multiple 'names' flag")
	}
	value = strings.Replace(value, " ", "", -1)
	for _, value := range strings.Split(value, ",") {
		*self = append(*self, value)
	}
	return nil
}

//This type and two methods augment the "-model" flag
type modelArray []string

func (self *modelArray) String() string {
	return fmt.Sprint(*self)
}
func (self *modelArray) Set(value string) error {
	*self = append(*self, value)
	return nil
}

//This type and two methods augment the "-include" flag
type templateFeatureArray []string

func (self *templateFeatureArray) String() string {
	return fmt.Sprint(*self)
}
func (self *templateFeatureArray) Set(value string) error {
	if len(*self) > 0 {
		return errors.New("multiple 'include' flag")
	}
	value = strings.Replace(value, " ", "", -1)
	for _, value := range strings.Split(value, ",") {
		*self = append(*self, value)
	}
	return nil
}

//The various flags this tool accepts
var configPath string
var example bool
var zipOutput bool
var modelName string
var modelNames nameArray
var model modelArray
var schemaPath string
var templatePath string
var forceOverwrite bool
var alwaysAsk bool
var packageString string
var projectName string
var templateFeatures templateFeatureArray
var getTemplateFeatures bool
var getVersion bool

func setupFlags() {
	fmt.Printf("")
	modelNames = make(nameArray, 0)
	model = make(modelArray, 0)
	flag.StringVar(&configPath, "config", "", "The full path to your configuration file")
	flag.StringVar(&configPath, "c", "", "")
	flag.StringVar(&projectName, "project", "", "The string to use wherever a template requires the name of the project")
	flag.StringVar(&projectName, "p", "", "")
	flag.StringVar(&packageString, "package", "", "The package string that templates will use in source code files")
	flag.StringVar(&packageString, "k", "", "")
	flag.StringVar(&modelName, "name", "", "The name of a model in the schema")
	flag.StringVar(&modelName, "n", "", "")
	flag.Var(&modelNames, "names", "The names of some models in the schema")
	flag.Var(&modelNames, "N", "")
	flag.Var(&model, "model", "A model definition with the format model_name[ property_name:primitive_type][ property_name:primitive_type][...]. It *must* be wrapped in quotation marks. eg. \"User FirstName:string Age:int Password:string\"")
	flag.Var(&model, "m", "")
	flag.StringVar(&schemaPath, "schema", "", "The full path to the schema")
	flag.StringVar(&schemaPath, "s", "", "")
	flag.StringVar(&templatePath, "template", "", "The full path to the template")
	flag.StringVar(&templatePath, "t", "", "")
	flag.BoolVar(&getTemplateFeatures, "list", false, "When this parameter is used in conjunction with the -template parameter, levo will describe the optional configuration flags specific to that set of templates")
	flag.Var(&templateFeatures, "features", "This commandline parameter is provided for [un]setting the optional features specific to a set of templates. Keywords 'all' and 'none' work as expected. Prepending '-' or '+' indicates that the feature will be unset or set respectively.")
	flag.Var(&templateFeatures, "f", "")
	flag.BoolVar(&zipOutput, "zip", false, "When set, the commandline tool will output a zip file instead of numerous source code files")
	flag.BoolVar(&zipOutput, "z", false, "")
	flag.BoolVar(&forceOverwrite, "quiet", false, "When set, the commandline tool will overwrite generated files without asking")
	flag.BoolVar(&forceOverwrite, "q", false, "")
	flag.BoolVar(&alwaysAsk, "ask", false, "When set, the commandline tool will ask for before overwriting every file. If not set, the tool will ask once and use that answer for all subsequent overwrites")
	flag.BoolVar(&alwaysAsk, "a", false, "")
	flag.BoolVar(&getVersion, "version", false, "Setting this flag will output Levo's version information")
	flag.BoolVar(&getVersion, "v", false, "")
	flag.BoolVar(&example, "example", false, "This flag will cause other flags to be ignored and will produce a directory that contains all of the files needed to form an example workspace")
}

func setupFlagUsage() {
	flag.Usage = func() {
		fmt.Println("levo [options] -config <file_path>")
		fmt.Println("levo [options] -model <model_json> [-model <model_json>] -template <file_path>")
		fmt.Println("levo [options] (-name <model_name> | -names <<model_name>,...>) -schema <file_path> -template <file_path>")
		fmt.Println("levo -template <file_path> -list")
		fmt.Println("levo -example")

		fmt.Println("\nArguments")
		fmt.Printf(printFlagUsage(flag.Lookup("config"), flag.Lookup("c"), "<file_path>"))
		fmt.Printf(printFlagUsage(flag.Lookup("name"), flag.Lookup("n"), "<model_name>"))
		fmt.Printf(printFlagUsage(flag.Lookup("names"), flag.Lookup("N"), "<model_name,...>"))
		fmt.Printf(printFlagUsage(flag.Lookup("model"), flag.Lookup("m"), "<model_def>"))
		fmt.Printf(printFlagUsage(flag.Lookup("schema"), flag.Lookup("s"), "<file_path>"))
		fmt.Printf(printFlagUsage(flag.Lookup("template"), flag.Lookup("t"), "<file_path>"))

		fmt.Println("\nOptions")
		fmt.Printf(printFlagUsage(flag.Lookup("list"), flag.Lookup(""), ""))
		fmt.Printf(printFlagUsage(flag.Lookup("features"), flag.Lookup("f"), "all,none,[+|-]<template_features>"))
		fmt.Printf(printFlagUsage(flag.Lookup("zip"), flag.Lookup("z"), ""))
		fmt.Printf(printFlagUsage(flag.Lookup("quiet"), flag.Lookup("q"), ""))
		fmt.Printf(printFlagUsage(flag.Lookup("ask"), flag.Lookup("a"), ""))
		fmt.Printf(printFlagUsage(flag.Lookup("version"), flag.Lookup("v"), ""))
		fmt.Printf(printFlagUsage(flag.Lookup("project"), flag.Lookup("p"), "<project_name>"))
		fmt.Printf(printFlagUsage(flag.Lookup("package"), flag.Lookup("k"), "<package>"))

		fmt.Println("\nExample")
		fmt.Printf(printFlagUsage(flag.Lookup("example"), nil, ""))
	}
}

func setupFlagUsageTesting() {
	flag.Usage = func() {
		return
	}
}

func printFlagUsage(mainFlag *flag.Flag, shortFlag *flag.Flag, inputType string) string {
	var signature string
	var output string
	if shortFlag != nil {
		signature = "  -" + shortFlag.Name + " (-" + mainFlag.Name + ") " + inputType
	} else {
		signature = "  -" + mainFlag.Name + " " + inputType
	}

	if len(mainFlag.Usage) < 72 {
		output = output + signature + "\n\t" + mainFlag.Usage + "\n"
	} else {
		output = output + signature + "\n\t" + mainFlag.Usage[0:72] + "\n"
		for i := 72; i < len(mainFlag.Usage); {
			if i+72 > len(mainFlag.Usage)-1 {
				output = output + "\t" + mainFlag.Usage[i:] + "\n"
				i += 72
			} else {
				output = output + "\t" + mainFlag.Usage[i:i+72] + "\n"
				i += 72
			}
		}
	}
	return output
}

func parseFlags() {
	flag.Parse()
}

func checkFlags() bool {
	if getVersion {
		return true
	} else if example {
		return true
	} else if configPath != "" && !configFlagGood() {
		fmt.Fprintf(os.Stderr, "When using -config, do not also use -model, -name, -names, -schema, or -template\n")
		flag.Usage()
		return false
	} else if len(model) > 0 && !modelFlagGood() {
		fmt.Fprintf(os.Stderr, "When using -model, -template must also be used\n")
		flag.Usage()
		return false
	} else if modelName != "" && !modelNameFlagGood() {
		fmt.Fprintf(os.Stderr, "When using -name, both -template and -schema must also be used\n")
		flag.Usage()
		return false
	} else if len(modelNames) > 0 && !modelNamesFlagGood() {
		fmt.Fprintf(os.Stderr, "When using -names, both -template and -schema must also be used\n")
		flag.Usage()
		return false
	} else if schemaPath != "" && !schemaFlagGood() {
		fmt.Fprintf(os.Stderr, "When using -schema, -template and one of -name or -names must also be used\n")
		flag.Usage()
		return false
	} else if forceOverwrite && alwaysAsk {
		fmt.Fprintf(os.Stderr, "-force and -ask are mutually exclusive\n")
		flag.Usage()
		return false
	} else if configPath == "" && len(model) <= 0 && modelName == "" && len(modelNames) <= 0 && templatePath == "" && !example {
		flag.Usage()
		return false
	} else if getTemplateFeatures && templatePath == "" {
		fmt.Fprintf(os.Stderr, "-list must be used in conjunction with -template\n")
		flag.Usage()
		return false
	}
	return true
}

func configFlagGood() bool {
	if configPath != "" {
		if len(model) > 0 || modelName != "" || len(modelNames) > 0 || schemaPath != "" || templatePath != "" {
			return false
		}
	}
	return true
}

func modelFlagGood() bool {
	if len(model) > 0 {
		if templatePath == "" {
			return false
		}
	}
	return true
}

func modelNameFlagGood() bool {
	if modelName != "" {
		if templatePath == "" || schemaPath == "" {
			return false
		}
	}
	return true
}

func modelNamesFlagGood() bool {
	if len(modelNames) > 0 {
		if templatePath == "" || schemaPath == "" {
			return false
		}
	}
	return true
}

func schemaFlagGood() bool {
	if schemaPath != "" {
		if templatePath == "" {
			return false
		} else if modelName == "" && len(modelNames) <= 0 {
			return false
		}
	}
	return true
}

func templateFlagGood() bool {
	if templatePath != "" {
		if len(model) <= 0 && modelName == "" && len(modelNames) <= 0 {
			return false
		} else if modelName != "" || len(modelNames) > 0 {
			if schemaPath == "" {
				return false
			}
		}
	}
	return true
}

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
	"flag"
	"fmt"
	"reflect"
	"testing"
)

func TestSetupFlags(testing *testing.T) {
	fmt.Printf("")
	defer func() {
		if r := recover(); r != nil {
			testing.Errorf("Something weird happend with the flags. Stupid flags")
		}
	}()

	if flag.Lookup("config").Name == "config" {
		//yay
	} else {
		setupFlags()
	}

	theFlag := flag.Lookup("config")
	if reflect.DeepEqual(theFlag, flag.Flag{}) == true {
		testing.Errorf("The config flag could not be found")
	}
	if theFlag.Value.String() != "" {
		testing.Errorf("Unusual inital value for config flag")
	}
	if theFlag.Usage != "The full path to your configuration file" {
		testing.Errorf("Unusual usage value for config flag")
	}
}

func TestPrintFlagUsage(testing *testing.T) {
	resetFlags()
	var main int
	var short int
	flag.IntVar(&main, "main", 0, "help message for flagname")
	flag.IntVar(&short, "short", 0, "")
	output := printFlagUsage(flag.Lookup("main"), flag.Lookup("short"), "<num>")
	expecting := "  -short (-main) <num>\n	help message for flagname\n"

	if output != expecting {
		testing.Errorf("Expecting:\n%s\nGot:\n%s\n", expecting, output)
	}

	//long message
	resetFlags()
	flag.IntVar(&main, "main", 0, "This is an unusually long message that explains all the various details about how this flag should (and shouldn't) be used by the person who is using it. The usage message might also contains newlines, but this one doesn't.")
	flag.IntVar(&short, "short", 0, "")
	output = printFlagUsage(flag.Lookup("main"), flag.Lookup("short"), "<num>")
	expecting = "  -short (-main) <num>\n	This is an unusually long message that explains all the various details \n	about how this flag should (and shouldn't) be used by the person who is \n	using it. The usage message might also contains newlines, but this one d\n	oesn't.\n"

	if output != expecting {
		testing.Errorf("Expecting:\n%s\nGot:\n%s\n", expecting, output)
	}
}

func TestCheckFlags(testing *testing.T) {
	defer resetFlags()
	resetFlags()
	ok := checkFlags()
	if ok {
		testing.Errorf("Should have thrown an error with empty args")
	}

	resetFlags()
	flag.Set("config", "")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown an error with empty config")
	}

	resetFlags()
	flag.Set("config", "anything")
	flag.Set("model", "anything")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown an error with config and any other flag")
	}

	resetFlags()
	flag.Set("config", "test-resources/code-gen-config.json")
	ok = checkFlags()
	if !ok {
		testing.Errorf("Good config should not have thrown errors")
	}

	resetFlags()
	flag.Set("model", "Test:Stuff:string")
	flag.Set("template", "path/to/template")
	ok = checkFlags()
	if !ok {
		testing.Errorf("Good model and template should not have thrown errors")
	}

	resetFlags()
	flag.Set("model", "Test:Stuff:string")
	flag.Set("model", "Other:Stuff:string")
	flag.Set("template", "path/to/template")
	ok = checkFlags()
	if !ok {
		testing.Errorf("Good model and template should not have thrown errors")
	}

	resetFlags()
	flag.Set("model", "Test:Stuff")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown error for model missing template")
	}

	resetFlags()
	flag.Set("name", "Test:Stuff")
	flag.Set("template", "path/to/template")
	flag.Set("schema", "path/to/schema")
	ok = checkFlags()
	if !ok {
		testing.Errorf("Good name and template and schema should not have thrown errors")
	}

	resetFlags()
	flag.Set("name", "Test:Stuff")
	flag.Set("schema", "path/to/schema")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown error for name missing template")
	}

	resetFlags()
	flag.Set("name", "Test:Stuff")
	flag.Set("template", "path/to/template")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown error for name missing schema")
	}

	resetFlags()
	flag.Set("names", "Test:Stuff")
	flag.Set("template", "path/to/template")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown error for name missing schema")
	}

	resetFlags()
	flag.Set("schema", "anything")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown error for only schema")
	}

	resetFlags()
	flag.Set("template", "path/to/template")
	ok = checkFlags()
	if !ok {
		testing.Errorf("Valid template should not throw error")
	}

	resetFlags()
	flag.Set("list", "true")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown error for list alone")
	}

	resetFlags()
	flag.Set("list", "true")
	flag.Set("template", "path/to/template")
	ok = checkFlags()
	if !ok {
		testing.Errorf("Valid template and list should not throw error")
	}

	resetFlags()
	flag.Set("force", "true")
	flag.Set("ask", "true")
	ok = checkFlags()
	if ok {
		testing.Errorf("Should have thrown error for force and ask")
	}
}

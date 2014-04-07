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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func helperCommand() *exec.Cmd {
	cmd := exec.Command("./levo")
	return cmd
}

func TestNotOutputExampleWorkspace(testing *testing.T) {
	defer cleanup()
	input := "n\n"
	p := helperCommand()
	p.Args = append(p.Args, "-example")
	p.Stdin = strings.NewReader(input)
	output, err := p.CombinedOutput()
	if err != nil {
		testing.Error(err)
	}
	fmt.Println(string(output))
	if strings.Contains(string(output), "Example files not created") == false {
		testing.Errorf(string(output))
	}
}

func TestOutputExampleWorkspace(testing *testing.T) {
	defer cleanup()
	input := "y\n"
	p := helperCommand()
	p.Args = append(p.Args, "-example")
	p.Stdin = strings.NewReader(input)
	output, err := p.Output()
	if err != nil {
		testing.Error(err)
	}
	if strings.Contains(string(output), "Successfully created") == false {
		testing.Errorf(string(output))
	}
}

func TestCheckForExampleDir(testing *testing.T) {
	var err error
	//Try removing example directory. Don't care if it doesn't work
	os.RemoveAll("example")

	err = checkForExampleDir()
	if err != nil {
		testing.Error("The example directory exists when it shouldn't")
	}

	err = os.Mkdir("example", 0755)
	if err != nil {
		testing.Errorf("Could not create example directory (this is not a code failure, it is a failure of the test itself)")
	} else {
		err = checkForExampleDir()
		if err == nil {
			testing.Error("The example directory doesn't exist when it should")
		}
	}
	err = os.RemoveAll("example")
	if err != nil {
		testing.Error("Could not remove the example directory (this is not a code failure, it is a failure of the test itself)")
	}
}

func TestCreateExampleDirectory(testing *testing.T) {
	var err error
	//Try removing example directory. Don't care if it doesn't work
	os.RemoveAll("example")

	err = createExampleDirectory()
	if err != nil {
		testing.Error("The example directory was not created when it should have been")
	}

	err = checkForExampleDir()
	if err == nil {
		testing.Errorf("The example directory isn't actually there (after a creation) %v", err.Error())
	}

	err = createExampleDirectory()
	if err == nil {
		testing.Error("The example directory claimed to be created when it couldn't have been")
	}

	//theoretically should test a non-writable directory, but that's a pain so we're putting it off

	err = os.RemoveAll("example")
	if err != nil {
		testing.Error("Could not remove the example directory (this is not a code failure, it is a failure of the test itself)")
	}
}

func TestPopulateConfigFiles(testing *testing.T) {
	var err error
	os.RemoveAll("example")

	//try creating config files without first setting up a directory
	err = populateConfigFiles()
	if err == nil {
		testing.Errorf("No error when populating a directory that does not exist")
	}

	//setup a directory
	err = createExampleDirectory()
	if err != nil {
		testing.Error(err.Error())
	}

	//populate the directory with config files
	err = populateConfigFiles()
	if err != nil {
		testing.Errorf("Error when creating config files: %s", err.Error())
	}

	//look for specific files
	fileInfos, err := ioutil.ReadDir("example")
	if err != nil {
		testing.Errorf("Error when reading example dir (OS level error, likely not code related): %s", err.Error())
	}
	//Looking for new files
	var configFile os.FileInfo
	var schemaFile os.FileInfo
	for _, fileInfo := range fileInfos {
		fmt.Printf("")
		if fileInfo.Name() == "config.json" {
			configFile = fileInfo
		}
		if fileInfo.Name() == "schema.json" {
			schemaFile = fileInfo
		}
	}

	//Checking config file
	if configFile == nil {
		testing.Errorf("Config not created")
	} else if configFile.Size() == 0 {
		testing.Errorf("Config created, but is empty")
	} //should check the actual contents of the file

	//Checking schema file
	if schemaFile == nil {
		testing.Errorf("Schema not created")
	} else if schemaFile.Size() == 0 {
		testing.Errorf("Schema created, but is empty")
	} //should check the actual contents of the file
	os.RemoveAll("example")
}

func TestPopulateTemplates(testing *testing.T) {
	var err error
	os.RemoveAll("templates")

	//populate the directory with template files
	err = populateTemplates()
	if err != nil {
		testing.Errorf("Error when creating template files: %s", err.Error())
	}

	//look for specific files
	fileInfos, err := ioutil.ReadDir("templates")
	if err != nil {
		testing.Errorf("Error when reading tempate dir (OS level error, likely not code related): %s", err.Error())
	}
	//Looking for new files
	var templateFile os.FileInfo
	for _, fileInfo := range fileInfos {
		fmt.Printf("")
		if fileInfo.Name() == "_Name_ListView.lt" {
			templateFile = fileInfo
		}
	}

	//Checking config file
	if templateFile == nil {
		testing.Errorf("template not created")
	} else if templateFile.Size() == 0 {
		testing.Errorf("template created, but is empty")
	} //should check the actual contents of the file

	os.RemoveAll("templates")
}

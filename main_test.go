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
	"github.com/cfmobile/levolib"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func cleanup() {
	_, err := os.Stat("Cats.generic")
	if err == nil {
		if err := os.RemoveAll("Cats.generic"); err != nil {
			fmt.Println(err.Error())
		}
	}
	_, err = os.Stat("Dogs.generic")
	if err == nil {
		if err := os.RemoveAll("Dogs.generic"); err != nil {
			fmt.Println(err.Error())
		}
	}
	_, err = os.Stat("People.nongeneric")
	if err == nil {
		if err := os.RemoveAll("People.nongeneric"); err != nil {
			fmt.Println(err.Error())
		}
	}

	_, err = os.Stat("example")
	if err == nil {
		if err := os.RemoveAll("example"); err != nil {
			fmt.Println(err.Error())
		}
	}
	_, err = os.Stat("levo_gen.zip")
	if err == nil {
		if err := os.Remove("levo_gen.zip"); err != nil {
			fmt.Println(err.Error())
		}
	}
	resetFlags()
}

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	setupFlags()
	setupFlagUsageTesting()
}

func printFlags() {
	flag.VisitAll(func(theFlag *flag.Flag) {
		fmt.Printf("%v:	%v\n", theFlag.Name, theFlag.Value.String())
	})
}

func TestMain(testing *testing.T) {
	defer cleanup()

	//test with no filename
	cleanup()
	main()

	//test with invalid filename
	cleanup()
	flag.Set("config", "thisisn'tagoodfilename")
	main()

	//test with valid filename
	cleanup()
	flag.Set("config", "test-resources/code-gen-config.json")
	main()

	//test with example flag
	cleanup()
	flag.Set("example", "true")
	main()

	//test with example flag
	cleanup()
	flag.Set("config", "test-resources/code-gen-config.json")
	flag.Set("z", "true")
	main()
}

func TestProcessArgs(testing *testing.T) {
	defer cleanup()

	//Test with no command line
	resetFlags()
	generatedFiles, err := processArgs()
	if err == nil {
		testing.Errorf("No error when Args is empty")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non-empty generatedFiles returned when parsing empty Args")
	}

	//Test with legit config
	resetFlags()
	flag.Set("config", "test-resources/code-gen-config.json")
	generatedFiles, err = processArgs()
	if err != nil {
		testing.Errorf("Error when processing valid config: %v", err.Error())
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == true {
		testing.Errorf("Empty generatedFiles returned when parsing valid config")
	}

	//Test with legit name and schema
	resetFlags()
	flag.Set("name", "Cats")
	flag.Set("schema", "test-resources/model-schema.json")
	flag.Set("template", "test-resources/templates/_Name_.generic.lt")
	generatedFiles, err = processArgs()
	if err != nil {
		testing.Errorf("Error when processing valid name: %v", err.Error())
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == true {
		testing.Errorf("Empty generatedFiles returned when parsing valid name")
	}

	//Test with legit name and schema
	resetFlags()
	flag.Set("name", "Cats")
	flag.Set("schema", "test-resources/model-schema.json")
	flag.Set("template", "test-resources/workingTemplates")
	generatedFiles, err = processArgs()
	if err != nil {
		testing.Errorf("Error when processing valid name: %v", err.Error())
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == true {
		testing.Errorf("Empty generatedFiles returned when parsing valid name")
	}

	//Test with legit names and schema
	resetFlags()
	flag.Set("names", "Cats,Dogs")
	flag.Set("schema", "test-resources/model-schema.json")
	flag.Set("template", "test-resources/templates/_Name_.generic.lt")
	generatedFiles, err = processArgs()
	if err != nil {
		testing.Errorf("Error when processing valid names")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == true {
		testing.Errorf("Empty generatedFiles returned when parsing valid names")
	}

	//Test with legit model and template
	resetFlags()
	flag.Set("model", "Test:First:string:Second:string")
	flag.Set("template", "test-resources/templates/_Name_.generic.lt")
	generatedFiles, err = processArgs()
	if err != nil {
		testing.Errorf("Error when processing valid model")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == true {
		testing.Errorf("Empty generatedFiles returned when parsing valid model")
	}

	//Test with broken model
	resetFlags()
	flag.Set("model", "Test:First:string:string")
	flag.Set("template", "test-resources/templates/_Name_.generic.lt")
	generatedFiles, err = processArgs()
	if err == nil {
		testing.Errorf("No error when processing invalid model")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non-empty generatedFiles returned when parsing valid config")
	}

	//Test with bad config
	resetFlags()
	flag.Set("config", "test-resources/codonfig.json")
	generatedFiles, err = processArgs()
	if err == nil {
		testing.Errorf("No error when config is bad")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non-empty generatedFiles returned when config is bad")
	}

	//Test with bad config
	resetFlags()
	flag.Set("name", "Garg!")
	flag.Set("schema", "test-resources/model-schema.json")
	flag.Set("template", "test-resources/templates/_Name_.generic.lt")
	generatedFiles, err = processArgs()
	if err == nil {
		testing.Errorf("No error when model name is bad")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non-empty generatedFiles returned when model name is bad")
	}

	//Test with -template and -list
	resetFlags()
	flag.Set("template", "test-resources/templates")
	flag.Set("list", "true")
	_, err = processArgs()
	if err != nil {
		testing.Errorf("Unexpected error: %v", err.Error())
	}

	//Test -list without -template
	resetFlags()
	flag.Set("list", "true")
	_, err = processArgs()
	if err == nil {
		testing.Errorf("No error when -list is used without -template")
	}
}

func TestGenerateFromConfiguration(testing *testing.T) {
	defer cleanup()
	generatedFiles, err := generateFromConfiguration("test-resources/code-gen-config.json")
	if err != nil {
		testing.Errorf("Error returned when processing valid config: " + err.Error())
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) {
		testing.Errorf("Empty generatedFiles returned when processing valid config")
	}

	generatedFiles, err = generateFromConfiguration("test-resources/thisisnotagoodconfig.jslfk")
	if err == nil {
		testing.Errorf("No error returned when passed a non-existant config file")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non-empty generatedFiles returned when parsing a broken config")
	}

	//test an empty config
	generatedFiles, err = generateFromConfiguration("")
	if err == nil {
		testing.Errorf("No error returned when processing empty config")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non empty array of generated files returned from blank generatedFiles")
	}

	//test a config with bad template
	generatedFiles, err = generateFromConfiguration("test-resources/code-gen-config-bad-data.json")
	if err == nil {
		testing.Errorf("No error returned when processing config with broken templates")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non empty array of generated files returned from config with broken templates")
	}
}

func TestProcessRawModel(testing *testing.T) {
	defer cleanup()
	models, err := processRawModel(modelArray{"Test:One:string:Two:string"})
	if err != nil {
		testing.Errorf("Error returned when processing valid model string: " + err.Error())
	} else if reflect.DeepEqual(models, []levo.Model{}) {
		testing.Errorf("Empty models returned when processing valid model string")
	}

	models, err = processRawModel(modelArray{"Test:One::Two:string"})
	if err == nil {
		testing.Errorf("No error returned when processing invalid model")
	} else if reflect.DeepEqual(models, []levo.Model{}) == false {
		testing.Errorf("Non empty generatedFiles returned when processing invalid model: %v", models)
	}

	models, err = processRawModel(modelArray{"Test:::Two:string"})
	if err == nil {
		testing.Errorf("No error returned when processing invalid model")
	} else if reflect.DeepEqual(models, []levo.Model{}) == false {
		testing.Errorf("Non empty generatedFiles returned when processing invalid model: %v", models)
	}

	models, err = processRawModel(modelArray{"Test"})
	if err != nil {
		testing.Errorf("Error returned when processing valid model string: " + err.Error())
	} else if reflect.DeepEqual(models, []levo.Model{}) {
		testing.Errorf("Empty models returned when processing valid model string")
	}
}

func TestProcessModelsFromSchema(testing *testing.T) {
	defer cleanup()
	models, err := processModelsFromSchema("Cats", []string{}, "test-resources/model-schema.json")
	if err != nil {
		testing.Errorf("Error returned when processing valid model from schema: " + err.Error())
	} else if reflect.DeepEqual(models, []levo.Model{}) {
		testing.Errorf("Empty models returned when processing valid model from schema")
	}

	models, err = processModelsFromSchema("", []string{"Cats", "Dogs"}, "test-resources/model-schema.json")
	if err != nil {
		testing.Errorf("Error returned when processing valid model from schema: " + err.Error())
	} else if reflect.DeepEqual(models, []levo.Model{}) {
		testing.Errorf("Empty models returned when processing valid model from schema")
	}

	models, err = processModelsFromSchema("Potatoes", []string{}, "test-resources/model-schema.json")
	if err == nil {
		testing.Errorf("No error returned when processing invalide model from schema")
	} else if reflect.DeepEqual(models, []levo.Model{}) == false {
		testing.Errorf("Non-empty models returned when processing invalid model from schema")
	}

	models, err = processModelsFromSchema("", []string{}, "test-resources/model-schema.json")
	if err == nil {
		testing.Errorf("No error returned when processing invalide model from schema")
	} else if reflect.DeepEqual(models, []levo.Model{}) == false {
		testing.Errorf("Non-empty models returned when processing invalid model from schema")
	}

	models, err = processModelsFromSchema("", []string{}, "test-resources/modema.json")
	if err == nil {
		testing.Errorf("No error returned when processing invalide model from schema")
	} else if reflect.DeepEqual(models, []levo.Model{}) == false {
		testing.Errorf("Non-empty models returned when processing invalid model from schema")
	}
}

func TestGenerateModelsAndTemplates(testing *testing.T) {
	defer cleanup()
	modelsGood, err := processModelsFromSchema("Cats", []string{}, "test-resources/model-schema.json")
	if err != nil {
		testing.Errorf("Error when trying to setup models for testing")
	}

	generatedFiles, err := generateModelsAndTemplates(modelsGood, "test-resources/templates/_Name_.generic.lt")
	if err != nil {
		testing.Errorf("Error when trying to generate valid model/template")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == true {
		testing.Errorf("Empty generated files array when trying to generate valid model/template")
	}

	flag.Set("p", "ProjectName")
	flag.Set("k", "ProjectName")
	generatedFiles, err = generateModelsAndTemplates(modelsGood, "test-resources/templates/_Name_.generic.lt")
	if err != nil {
		testing.Errorf("Error when trying to generate valid model/template")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == true {
		testing.Errorf("Empty generated files array when trying to generate valid model/template")
	}

	generatedFiles, err = generateModelsAndTemplates([]levo.Model{}, "test-resources/templates/_Name_.generic.lt")
	if err != nil {
		testing.Errorf("Error when trying to generate valid model/template")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non-Empty generated files array when trying to generate empty model/template")
	}

	generatedFiles, err = generateModelsAndTemplates(modelsGood, "test-resources/templates/_Name_.gesfric")
	if err == nil {
		testing.Errorf("No error when trying to generate valid model/template")
	} else if reflect.DeepEqual(generatedFiles, []levo.GeneratedFile{}) == false {
		testing.Errorf("Non-Empty generated files array when trying to generate empty model/template")
	}
}

func TestWriteFileData(testing *testing.T) {
	defer cleanup()
	generatedFiles, _ := generateFromConfiguration("test-resources/code-gen-config.json")
	var expectedFileCount = len(generatedFiles)
	err := outputFiles(generatedFiles)
	if err != nil {
		testing.Errorf("Error returned when processing valid generated files")
	}

	generatedFilesInfo, err := ioutil.ReadDir(".")
	if err != nil {
		testing.Errorf("Got filesystem error when looking for files -- likely a system issue")
	}

	for _, generatedFileInfo := range generatedFilesInfo {
		for _, expectedFile := range generatedFiles {
			if generatedFileInfo.Name() == expectedFile.FileName {
				expectedFileCount = expectedFileCount - 1
			}
		}
	}

	if expectedFileCount != 0 {
		if expectedFileCount > 0 {
			testing.Errorf("Not all expected files were present")
		} else if expectedFileCount < 0 {
			testing.Errorf("Duplicate files were present (this is a system error)")
		}
	}

	//test writing garbage files
	generatedFiles = []levo.GeneratedFile{levo.GeneratedFile{FileName: "!#//$%.^#.$%^.as.df$$$", Body: []byte("stuffs")}}
	err = outputFiles(generatedFiles)
	if err == nil {
		testing.Errorf("No error returned when writing an obviously broken filename")
	}
}

func TestOverwritingData(testing *testing.T) {
	defer cleanup()
	generatedFiles, _ := generateFromConfiguration("test-resources/code-gen-config.json")
	err := outputFiles(generatedFiles)
	if err != nil {
		testing.Errorf("Error returned when processing valid generated files: %s", err.Error())
	}

	flag.Set("quiet", "true")
	err = outputFiles(generatedFiles)
	if err != nil {
		testing.Errorf("Error returned when processing valid generated files: %s", err.Error())
	}
}

func TestWriteZipFile(testing *testing.T) {
	defer cleanup()
	generatedFiles, _ := generateFromConfiguration("test-resources/code-gen-config.json")
	err := writeZipFile(generatedFiles)
	if err != nil {
		testing.Errorf("Error returned when processing valid generated files")
	}

	if _, err := os.Stat("levo_gen.zip"); err != nil {
		testing.Error(err.Error())
	}

	//test writing garbage files
	generatedFiles = []levo.GeneratedFile{levo.GeneratedFile{FileName: "!#//$%.^#.$%^.as.df$$$", Body: []byte("stuffs")}}
	err = outputFiles(generatedFiles)
	if err == nil {
		testing.Errorf("No error returned when writing an obviously broken filename")
	}
}

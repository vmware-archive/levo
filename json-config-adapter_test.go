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
	"github.com/cfmobile/levolib"
	"reflect"
	"testing"
)

const TestBasePackage string = "com.test"
const TestLanguage string = "java"
const TestVersion string = "1.0"
const TestProjectName string = "test"
const TestModelName01 string = "People"
const TestModelName02 string = "Cats"
const TestTemplateName01 string = "_Name_.lt"
const TestPropName01 string = "hairy"
const TestPropName02 string = "bald"

var TestInvalidJSONString []byte = []byte("{this is completely {inv:alid json")
var TestValidCorrectConfigString []byte = []byte("{\"TemplaterVersion\":\"" + TestVersion + "\",\"BasePackage\":\"" + TestBasePackage + "\",\"Language\":\"" + TestLanguage + "\",\"Zip\":false,\"ModelSchemaFileName\":\"test-resources/model-schema.json\",\"TemplatesDirectory\":\"test-resources/templates\",\"Mappings\":[{\"ModelNames\":[\"" + TestModelName01 + "\",\"" + TestModelName02 + "\"],\"TemplateNames\":[\"" + TestTemplateName01 + "\"]},{\"ModelNames\":[\"" + TestModelName01 + "\"],\"TemplateNames\":[\"" + TestTemplateName01 + "\"]}]}")
var TestValidIncorrectConfigString []byte = []byte("{\"ModelSchemaFileName\":\"schema.json\",\"TemplatesDirectory\":\"templates/xl-rest_lib-android-v3.0.0\",\"Mappings\":[{\"model_name\":[\"Movie\"],\"templates\":[\"_NAME_Fragment.java\",\"list_item__name_.xml\",\"_NAME_.java\",\"_NAME_Activity.java\",\"_NAME_ContentProvider.java\",\"fragment__name_.xml\",\"activity__name_.xml\",\"_NAME_Validator.java\",\"_NAME_Application.java\",\"_NAME_Table.java\",\"_NAME_ListActivity.java\",\"_NAME_ListValidator.java\",\"Abs_NAME_.java\",\"fragment__name__list.xml\",\"_NAME_ListFragment.java\",\"activity__name__list.xml\",\"Abs_NAME_Table.java\"]},{\"model_name\":[\"Posters\"],\"templates\":[\"_NAME_.java\",\"Abs_NAME_.java\"]},{\"model_name\":[\"Cast\"],\"templates\":[\"_NAME_.java\",\"Abs_NAME_.java\"]},{\"model_name\":[\"Ratings\"],\"templates\":[\"_NAME_.java\",\"Abs_NAME_.java\"]},{\"model_name\":[\"MoviesResponse\",\"ReleaseDates\"],\"templates\":[\"_NAME_.java\",\"Abs_NAME_.java\"]},{\"model_name\":[],\"templates\":[\"_NAME_Application.java\"]}]}")

var TestValidSchemaString []byte = []byte("{\"Project\":\"" + TestProjectName + "\",\"Models\":[{\"Name\":\"" + TestModelName01 + "\",\"Parent\":\"\",\"Properties\":[{\"RemoteIdentifier\":\"" + TestPropName01 + "\",\"PropertyType\":\"string\"}]},{\"Name\":\"" + TestModelName02 + "\",\"Parent\":\"\",\"Properties\":[{\"RemoteIdentifier\":\"" + TestPropName02 + "\",\"PropertyType\":\"string\"}]}]}")
var TestValidIncorrectSchemaString []byte = []byte("{\"Project\":\"" + TestProjectName + "\",\"Models\":[{\"Parent\":\"\",\"Properties\":[{\"RemoteIdentifier\":\"" + TestPropName01 + "\",\"PropertyType\":\"string\"}]},{\"Name\":\"" + TestModelName02 + "\",\"Parent\":\"\",\"Properties\":[{\"RemoteIdentifier\":\"" + TestPropName02 + "\",\"PropertyType\":\"string\"}]}]}")

func TestProcessConfigurationFile(testing *testing.T) {
	fmt.Printf("")

	var configAdapter JSONConfigAdapter = JSONConfigAdapter{}

	returnedContext, err := configAdapter.ProcessConfigurationFile("obviouslynotavalidfilename")
	if err == nil {
		testing.Errorf("No Error when processing non-existent file")
	} else if reflect.DeepEqual(returnedContext, levo.Context{}) == false {
		testing.Errorf("Non empty context returned when processing non-existent file")
	}

	//Test a valid and correct config file
	returnedContext, err = configAdapter.ProcessConfigurationFile("test-resources/code-gen-config.json")

	if err != nil {
		testing.Errorf("ProcessConfigurationFile threw error while processing valid config: %v", err.Error())
	} else if reflect.DeepEqual(returnedContext, levo.Context{}) {
		testing.Errorf("Empty context returned while processing valid configuration file")
	}

	//Check that context has templates
	if len(returnedContext.Templates) <= 0 {
		testing.Errorf("Returned context does not have any templates while processing valid config")
	}

	//Check that context has models
	if len(returnedContext.Schema.Models) <= 0 {
		testing.Errorf("Returned context does not have any models while processing valid config")
	}

	//Check that context has Mappings
	if len(returnedContext.Mappings) <= 0 {
		testing.Errorf("Returned context does not have any mappings between models and templates while processing valid config")
	}

	//Check that context has other values (language, base package, etc)
	if returnedContext.PackageName != TestBasePackage || returnedContext.Language != TestLanguage || returnedContext.TemplaterVersion != TestVersion {
		testing.Errorf("Returned context is missing information while processing valid config")
	}

	//Test an invalid config file
	configAdapter = JSONConfigAdapter{}
	returnedContext, err = configAdapter.ProcessConfigurationFile("test-resources/code-gen-config-invalid.json")
	if err == nil {
		testing.Errorf("Empty error returned while processing invalid json config file")
	}

	//Test a valid but incorrect config file
	configAdapter = JSONConfigAdapter{}
	returnedContext, err = configAdapter.ProcessConfigurationFile("test-resources/code-gen-config-incorrect.json")
	if err == nil {
		testing.Errorf("Empty error returned while processing incorrect json config file")
	}
}

func TestProcessConfigurationString(testing *testing.T) {

	var configAdapter JSONConfigAdapter = JSONConfigAdapter{}

	context, err := configAdapter.ProcessConfigurationString(TestValidCorrectConfigString)
	if err != nil {
		testing.Errorf("Error while processing valid JSON config: %s", err.Error())
	}

	if context.Language == "" {
		testing.Error("Context language empty after processing valid JSON config")
	}

	context, err = configAdapter.ProcessConfigurationString(TestValidIncorrectConfigString)
	if err == nil {
		testing.Error("No error when processing invalid config (valid JSON)")
	}

	if context.Language != "" {
		testing.Errorf("Context language not empty after processing invalid config (valid JSON) config: %s", context.Language)
	}

	context, err = configAdapter.ProcessConfigurationString(TestInvalidJSONString)
	if err == nil {
		testing.Errorf("No error when processing invalid config (invalid JSON)")
	}
}

func TestParseConfigurationString(testing *testing.T) {

	var configAdapter JSONConfigAdapter = JSONConfigAdapter{}

	//Test valid and correct json
	err := configAdapter.ParseConfigurationString(TestValidCorrectConfigString)
	if err != nil {
		testing.Errorf("Error while parsing valid JSON config: ", err.Error())
	}
	if configAdapter.BasePackage != TestBasePackage || configAdapter.Language != TestLanguage || configAdapter.TemplaterVersion != TestVersion {
		testing.Errorf("JSON incorrectly parsed")
	}
	if len(configAdapter.Mappings) != 2 {
		testing.Errorf("Incorrect number of parsed Mappings")
	} else {
		mapping := configAdapter.Mappings[0]
		if len(mapping.ModelNames) != 2 {
			testing.Errorf("Expected %v model names. Got %v", 2, len(mapping.ModelNames))
		} else if mapping.ModelNames[0] != TestModelName01 {
			testing.Errorf("Was expecting ModelName %v. Got %v", TestModelName01, configAdapter.Mappings[0].ModelNames[0])
		}
		if len(mapping.TemplateNames) != 1 {
			testing.Errorf("Expected %v template names. Got %v", 1, len(mapping.TemplateNames))
		} else if mapping.TemplateNames[0] != TestTemplateName01 {
			testing.Errorf("Was expecting TemplateName %v. Got %v", TestTemplateName01, configAdapter.Mappings[0].TemplateNames[0])
		}
	}

	//Test invalid json
	configAdapter = JSONConfigAdapter{}
	err = configAdapter.ParseConfigurationString(TestInvalidJSONString)
	if err == nil {
		testing.Errorf("ParseConfigurationString did not fail when passed invalid JSON")
	}
	//Test valid and incorrect json
	configAdapter = JSONConfigAdapter{}
	err = configAdapter.ParseConfigurationString(TestValidIncorrectConfigString)
	if err == nil {
		testing.Errorf("ParseConfigurationString did not fail when passed incorrect JSON")
	}
}

func TestAddModelsToContext(testing *testing.T) {
	var configAdapter JSONConfigAdapter = JSONConfigAdapter{}
	models := []levo.Model{levo.Model{Name: TestModelName01}, levo.Model{Name: TestModelName02, Parent: TestModelName01}}
	err := configAdapter.addModelsToContext(levo.Schema{Project: TestProjectName, Models: models})
	if err != nil {
		testing.Errorf("Error while adding valid model: ", err.Error())
	}

	//Test invalid model name
	err = configAdapter.addModelsToContext(levo.Schema{Project: TestProjectName, Models: []levo.Model{levo.Model{Name: ""}}})
	if err == nil {
		testing.Errorf("No error when adding model without a name to context")
	}

	//Test invalid parent
	err = configAdapter.addModelsToContext(levo.Schema{Project: TestProjectName, Models: []levo.Model{levo.Model{Name: "AnyRandomName", Parent: "NotARealModel"}}})
	if err == nil {
		testing.Errorf("No error when adding model with invalid parent to context")
	}

	//Test invalid property
	property01 := levo.ModelProperty{RemoteIdentifier: TestPropName01}
	schema := levo.Schema{Project: TestProjectName, Models: []levo.Model{levo.Model{Name: "AnyRandomName02", Parent: "", Properties: []levo.ModelProperty{property01}}}}
	err = configAdapter.addModelsToContext(schema)
	if err == nil {
		testing.Errorf("No error when adding model with invalid parent to context")
	}
}

func TestJSONConfigObjectValidate(testing *testing.T) {
	//TODO Do it
}

func TestSchemaObjectValidate(testing *testing.T) {
	//TODO Do it
}

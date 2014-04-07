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
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cfmobile/levolib"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func init() {
	setupFlags()
	setupFlagUsage()
}

func main() {
	fmt.Printf("")
	parseFlags()
	if !checkFlags() {
		return
	}

	generatedFiles, err := processArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if len(generatedFiles) > 0 {
		if zipOutput {
			err := writeZipFile(generatedFiles)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing zip: ", err.Error())
				return
			}
		} else {
			err := outputFiles(generatedFiles)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing files: ", err.Error())
				return
			}
		}
	}
}

func processArgs() ([]levo.GeneratedFile, error) {
	if example {
		err := outputExampleWorkspace()
		if err != nil && err.Error() != "User Input: n" {
			return []levo.GeneratedFile{}, errors.New("Error creating example files: " + err.Error())
		}
		if err != nil && err.Error() == "User Input: n" {
			return []levo.GeneratedFile{}, errors.New("Example files not created")
		}
		fmt.Println("Successfully created example directory.\nEnter that directory, take a look, and then try 'levo -config config.json'")
		return []levo.GeneratedFile{}, nil
	}

	var err error
	templatePath, err = getUpdatedTemplateRepo(templatePath)
	if err != nil {
		return []levo.GeneratedFile{}, err
	}

	if getTemplateFeatures && templatePath != "" {
		possibleFlags, err := getTemplateFeaturesFromReadMe(templatePath)
		if err != nil {
			return []levo.GeneratedFile{}, err
		}
		for _, flagParts := range possibleFlags {
			fmt.Printf("%v:\n%v\n", flagParts[0], flagParts[1])
		}
		return []levo.GeneratedFile{}, nil
	}

	var generatedFiles []levo.GeneratedFile
	if configPath != "" {
		generatedFiles, err = generateFromConfiguration(configPath)
		if err != nil {
			return []levo.GeneratedFile{}, errors.New("Error reading config: " + err.Error())
		}
		return generatedFiles, nil

	} else if len(model) != 0 {
		models, err := processRawModel(model)
		if err != nil {
			return []levo.GeneratedFile{}, errors.New("Error parsing model string: " + err.Error())
		}
		generatedFiles, err = generateModelsAndTemplates(models, templatePath)
		if err != nil {
			return []levo.GeneratedFile{}, err
		}
		return generatedFiles, nil

	} else if modelName != "" || len(modelNames) > 0 {
		models, err := processModelsFromSchema(modelName, modelNames, schemaPath)
		if err != nil {
			return []levo.GeneratedFile{}, err
		}
		generatedFiles, err = generateModelsAndTemplates(models, templatePath)
		if err != nil {
			return []levo.GeneratedFile{}, err
		}
		return generatedFiles, nil
	} else if templatePath != "" {
		generatedFiles, err = generateModelsAndTemplates(make([]levo.Model, 0), templatePath)
		if err != nil {
			return []levo.GeneratedFile{}, err
		}
		return generatedFiles, nil
	} else {
		return []levo.GeneratedFile{}, errors.New("Unhandled request")
	}
}

func generateFromConfiguration(configFile string) ([]levo.GeneratedFile, error) {
	configAdapter := JSONConfigAdapter{}
	context, err := configAdapter.ProcessConfigurationFile(configFile)
	if err != nil {
		return []levo.GeneratedFile{}, errors.New("Error processing config: " + err.Error())
	}
	generatedFiles, err := levo.ProcessMappings(context)
	if err != nil {
		return []levo.GeneratedFile{}, errors.New("Error generating files from config: " + err.Error())
	}
	return generatedFiles, nil
}

func generateModelsAndTemplates(models []levo.Model, templatePath string) ([]levo.GeneratedFile, error) {
	context := levo.BeginContext()

	if packageString != "" {
		context.PackageName = packageString
	}
	if projectName != "" {
		context.ProjectName = projectName
	}

	for _, newModel := range models {
		addedModel, err := context.AddModelWithName(newModel.Name)
		if err != nil {
			return []levo.GeneratedFile{}, errors.New("Error adding models: " + err.Error())
		}
		for _, newProperty := range newModel.Properties {
			_, err := addedModel.AddProperty(newProperty.RemoteIdentifier, newProperty.LocalIdentifier, newProperty.PropertyType)
			if err != nil {
				return []levo.GeneratedFile{}, errors.New("Error adding properties: " + err.Error())
			}
		}
	}

	templates, err := addTemplatePath(&context, templatePath)
	if err != nil {
		return []levo.GeneratedFile{}, errors.New("Error adding template: " + err.Error())
	}

	for _, templateFeature := range templateFeatures {
		if templateFeature == "all" {
			//set all possible features
			possibleFlags, err := getTemplateFeaturesFromReadMe(templatePath)
			if err != nil {
				return []levo.GeneratedFile{}, err
			}
			for _, possibleFlag := range possibleFlags {
				context.AddTemplateFeature(possibleFlag[0])
			}
		} else if templateFeature == "none" {
			//unset all possible features
			possibleFlags, err := getTemplateFeaturesFromReadMe(templatePath)
			if err != nil {
				return []levo.GeneratedFile{}, err
			}
			for _, possibleFlag := range possibleFlags {
				context.RemoveTemplateFeature(possibleFlag[0])
			}
		} else {
			if templateFeature[0:1] == "-" {
				//unset this flag
				context.RemoveTemplateFeature(templateFeature[1:])
			} else if templateFeature[0:1] == "+" {
				//set this flag
				context.AddTemplateFeature(templateFeature[1:])
			} else {
				context.AddTemplateFeature(templateFeature)
			}
		}
	}

	err = addMappings(&context, templates, models)
	if err != nil {
		return []levo.GeneratedFile{}, errors.New("Error adding mapping: " + err.Error())
	}

	generatedFiles, err := levo.ProcessMappings(context)
	if err != nil {
		return []levo.GeneratedFile{}, errors.New("Error generating files: " + err.Error())
	}
	return generatedFiles, nil
}

func processRawModel(modelsString modelArray) ([]levo.Model, error) {
	models := make([]levo.Model, 0)
	for _, modelString := range modelsString {
		modelObj, err := getModel(modelString)
		if err != nil {
			return []levo.Model{}, errors.New("Error parsing model string: " + err.Error())
		}
		models = append(models, modelObj)
	}
	return models, nil
}

func processModelsFromSchema(modelName string, modelNames []string, schemaPath string) ([]levo.Model, error) {
	schemaAdapter := levo.GetJSONSchemaAdapter()
	schemaObject, err := schemaAdapter.ProcessSchemaFile(schemaPath)
	if err != nil {
		return []levo.Model{}, errors.New("Error while reading schema file: " + err.Error())
	}

	if len(modelNames) <= 0 {
		modelNames = append(modelNames, modelName)
	}

	models := make([]levo.Model, 0)
	for _, requestedName := range modelNames {
		foundInSchema := false
		for _, model := range schemaObject.Models {
			if model.Name == requestedName {
				foundInSchema = true
				models = append(models, model)
			}
		}
		if !foundInSchema {
			return []levo.Model{}, errors.New("Schema " + schemaPath + " does not contain " + requestedName)
		}
	}
	return models, nil
}

func getModel(modelString string) (levo.Model, error) {
	re := regexp.MustCompile("[^A-z0-9_]")
	modelParts := re.Split(modelString, -1)
	properties := make([]levo.ModelProperty, 0)
	modelName := modelParts[0]
	modelParts = modelParts[1:]
	if len(modelParts)%2 != 0 {
		return levo.Model{}, errors.New("Odd number of property parts")
	}
	for i := 0; i < len(modelParts); i += 2 {
		propName := modelParts[i]
		re := regexp.MustCompile(" ")
		localName := re.ReplaceAll([]byte(propName), []byte(""))
		propType := modelParts[i+1]
		if propName == "" || propType == "" {
			return levo.Model{}, errors.New("Name or type of property is empty string")
		}
		property := levo.ModelProperty{RemoteIdentifier: propName, PropertyType: propType, LocalIdentifier: string(localName)}
		properties = append(properties, property)
	}
	return levo.Model{Name: modelName, Properties: properties}, nil
}

func addMappings(context *levo.Context, templates []levo.TemplateInfo, models []levo.Model) error {
	templateNames := make([]string, 0)
	modelNames := make([]string, 0)
	for _, template := range templates {
		// Only map levo templates, not binary files
		if strings.HasSuffix(template.FileName, ".lt") {
			templateNames = append(templateNames, template.FileName)
		}
	}
	for _, model := range models {
		modelNames = append(modelNames, model.Name)
	}
	err := context.AddTemplatesForModelsMapping(templateNames, modelNames)
	if err != nil {
		return err
	}
	return nil
}

func getTemplateFeaturesFromReadMe(templatePath string) ([][]string, error) {
	featuresFromReadMe := make([][]string, 0)
	readMePath := templatePath + "/README.md"
	contents, err := ioutil.ReadFile(readMePath)
	if err != nil {
		return featuresFromReadMe, err
	}
	featuresRegex := regexp.MustCompile("#### *([A-z]+)\n(([^\n]+\n)*)")
	matches := featuresRegex.FindAllStringSubmatch(string(contents), -1)
	for _, featureParts := range matches {
		featuresFromReadMe = append(featuresFromReadMe, featureParts[1:3])
	}
	if len(featuresFromReadMe) == 0 {
		return featuresFromReadMe, errors.New("No template flags found in " + readMePath)
	}
	return featuresFromReadMe, nil
}

func outputFiles(generatedFiles []levo.GeneratedFile) error {
	//come back to here when we're done
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	overWrite := forceOverwrite
	input := bufio.NewReader(os.Stdin)
	for _, generatedFile := range generatedFiles {
		if generatedFile.Directory != "" {
			if err := os.MkdirAll(generatedFile.Directory, 0755); err != nil && !os.IsExist(err) {
				return err
			}
			if err := os.Chdir(generatedFile.Directory); err != nil {
				return err
			}
		}
		_, err := os.Open(generatedFile.FileName)
		//file, err := os.Stat(generatedFile.FileName)
		if err == nil {
			//File already exists
			if overWrite {
				err := writeFile(generatedFile.FileName, []byte(generatedFile.Body))
				if err != nil {
					return err
				}
			} else {
				//Ask the user if they want to overwrite this file
				fmt.Printf("The file %s/%s already exists. Overwrite? (y/n) : ", generatedFile.Directory, generatedFile.FileName)
				answer, err := input.ReadString('\n')
				if err != nil {
					return err
				}
				if answer == "y\n" {
					err := writeFile(generatedFile.FileName, []byte(generatedFile.Body))
					if err != nil {
						return err
					}
					if !alwaysAsk {
						overWrite = true
					}
				}
			}
		} else {
			err := writeFile(generatedFile.FileName, []byte(generatedFile.Body))
			if err != nil {
				return err
			}
		}
		if err := os.Chdir(originalDir); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(fileName string, contents []byte) error {
	if len(contents) >= 14 && string(contents[0:14]) == "<<levobase64>>" {
		headerlessContents := contents[14:]

		decodedContents := make([]byte, base64.StdEncoding.DecodedLen(len(headerlessContents)))
		i, err := base64.StdEncoding.Decode(decodedContents, headerlessContents)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fileName, decodedContents[:i], 0755)
		if err != nil {
			return err
		}
	} else {
		err := ioutil.WriteFile(fileName, contents, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeZipFile(generatedFiles []levo.GeneratedFile) error {
	buffer := bytes.NewBuffer(nil)
	zipWriter := zip.NewWriter(buffer)

	for _, generatedFile := range generatedFiles {
		file, err := zipWriter.Create(generatedFile.Directory + generatedFile.FileName)
		if err != nil {
			return err
		}
		_, err = file.Write(generatedFile.Body)
		if err != nil {
			return err
		}
	}
	err := zipWriter.Close()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("levo_gen.zip", buffer.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

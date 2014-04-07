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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cfmobile/levolib"
	"io/ioutil"
)

type JSONConfigAdapter struct {
	context             levo.Context
	TemplatesDirectory  string
	ModelSchemaFileName string
	TemplaterVersion    string
	Mappings            []modelToTemplateMapping
	BasePackage         string
	Language            string
	TemplateFeatures    []string
	Zip                 bool
}

type modelToTemplateMapping struct {
	ModelNames    []string
	TemplateNames []string
}

type schemaObject struct {
	Project string
	Models  []levo.Model
}

func (self *JSONConfigAdapter) ProcessConfigurationFile(fileName string) (levo.Context, error) {
	fmt.Printf("")
	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return levo.Context{}, err
	}
	return self.ProcessConfigurationString(fileContents)
}

func (self *JSONConfigAdapter) ProcessConfigurationString(configString []byte) (levo.Context, error) {
	fmt.Printf("")

	err := self.ParseConfigurationString(configString)
	if err != nil {
		return levo.Context{}, err
	}

	self.context = levo.BeginContext()
	self.context.PackageName = self.BasePackage
	self.context.Language = self.Language
	self.context.TemplaterVersion = self.TemplaterVersion
	for _, templateFeature := range self.TemplateFeatures {
		self.context.AddTemplateFeature(templateFeature)
	}

	schemaAdapter := levo.GetJSONSchemaAdapter()
	modelSchema, err := schemaAdapter.ProcessSchemaFile(self.ModelSchemaFileName)
	if err != nil {
		return levo.Context{}, err
	}
	if err := self.addModelsToContext(modelSchema); err != nil {
		return levo.Context{}, err
	}

	self.TemplatesDirectory, err = getUpdatedTemplateRepo(self.TemplatesDirectory)
	if err != nil {
		return levo.Context{}, err
	}

	_, err = addTemplatePath(&self.context, self.TemplatesDirectory)
	if err != nil {
		return levo.Context{}, err
	}

	if err := self.addProjectNameToContext(modelSchema); err != nil {
		return levo.Context{}, err
	}

	err = self.addMappingsToContext(self.Mappings)
	if err != nil {
		return levo.Context{}, err
	}

	if err := self.validate(); err == nil {
		return self.context, nil
	} else {
		return levo.Context{}, err
	}
}

func (self *JSONConfigAdapter) ParseConfigurationString(configString []byte) error {
	if err := json.Unmarshal(configString, self); err != nil {
		return err
	}
	if err := self.validate(); err != nil {
		return err
	}
	return nil
}

func (self *JSONConfigAdapter) addProjectNameToContext(modelSchema levo.Schema) error {
	if modelSchema.Project == "" {
		return errors.New("Schema did not define a Project name")
	}
	self.context.ProjectName = modelSchema.Project
	return nil
}

func (self *JSONConfigAdapter) addModelsToContext(modelSchema levo.Schema) error {
	for _, modelFromSchema := range modelSchema.Models {
		model, err := self.context.AddModelWithName(modelFromSchema.Name)
		if err != nil {
			return err
		}
		model.Parent = modelFromSchema.Parent
		for _, propertyFromSchema := range modelFromSchema.Properties {
			if propertyFromSchema.LocalIdentifier == "" {
				propertyFromSchema.LocalIdentifier = propertyFromSchema.RemoteIdentifier
			}
			_, err := model.AddProperty(propertyFromSchema.RemoteIdentifier, propertyFromSchema.LocalIdentifier, propertyFromSchema.PropertyType)
			if err != nil {
				return err
			}
		}
	}

	for _, model := range self.context.Schema.Models {
		if model.Parent != "" {
			var err error
			model.ParentRef, err = self.context.ModelForName(model.Parent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *JSONConfigAdapter) addMappingsToContext(mappings []modelToTemplateMapping) error {
	for _, mapping := range mappings {
		err := self.context.AddTemplatesForModelsMapping(mapping.TemplateNames, mapping.ModelNames)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *JSONConfigAdapter) validate() error {
	if self.BasePackage == "" {
		return errors.New("Configuration did not define a Base Package")
	} else if self.Language == "" {
		return errors.New("Configuration did not define a Language")
	} else if self.TemplaterVersion == "" {
		return errors.New("Configuration did not define a Templater Version")
	}
	//TODO fill this out more
	return nil
}

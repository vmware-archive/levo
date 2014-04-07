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
	"fmt"
	"io/ioutil"
	"os"
)

func outputExampleWorkspace() error {
	fmt.Printf("We are about to create an example workspace. Are you sure? (y/n): ")
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	if string(b) == "n" {
		return errors.New("User Input: n")
	}

	if err := createExampleDirectory(); err != nil {
		return err
	}

	if err := populateConfigFiles(); err != nil {
		return err
	}
	return nil
}

func createExampleDirectory() error {
	//make a directory
	//check that the directory doesn't already exist
	if err := checkForExampleDir(); err != nil {
		return err
	}

	if err := os.Mkdir("example", 0755); err != nil {
		return err
	}
	return nil
}

func checkForExampleDir() error {
	_, err := ioutil.ReadDir("example")
	if err == nil {
		return errors.New("Directory named 'example' already exists")
	}
	return nil
}

func populateConfigFiles() error {
	if err := os.Chdir("example"); err != nil {
		return err
	}

	if err := ioutil.WriteFile("config.json", []byte(configContents), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile("schema.json", []byte(schemaContents), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile("README", []byte(readmeContents), 0755); err != nil {
		return err
	}

	if err := populateTemplates(); err != nil {
		return err
	}
	if err := os.Chdir(".."); err != nil {
		return err
	}
	return nil
}

func populateTemplates() error {
	if err := os.Mkdir("templates", 0755); err != nil {
		return err
	}
	if err := os.Chdir("templates"); err != nil {
		return err
	}

	if err := ioutil.WriteFile("_Name_DescriptionView.lt", []byte(descriptionViewTemplateContents), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile("_Name_ListView.lt", []byte(listViewTemplateContents), 0755); err != nil {
		return err
	}
	if err := os.Chdir(".."); err != nil {
		return err
	}
	return nil
}

var configContents string = `{
  "TemplaterVersion": "1.0",
  "BasePackage": "com.pizzashop",
  "Language": "java",
  "ModelSchemaFileName": "schema.json",
  "TemplatesDirectory": "templates",
  "Mappings": [
    {
      "ModelNames": [
        "Product"
      ],
      "TemplateNames": [
        "_Name_ListView.lt"
      ]
    },
    {
      "ModelNames": [
        "Pizza",
        "Topping"
      ],
      "TemplateNames": [
        "_Name_DescriptionView.lt"
      ]
    }
  ]
}`

var schemaContents string = `{
  "Project": "PizzaShop",
  "Models": [
    {
      "Name": "Product",
      "Parent": "",
      "Properties": [
        {
          "RemoteIdentifier": "Name",
          "PropertyType": "string"
        },
        {
          "RemoteIdentifier": "Price",
          "PropertyType": "string"
        }
      ]
    },
    {
      "Name": "Pizza",
      "Parent": "Product",
      "Properties": [
        {
          "RemoteIdentifier": "Toppings",
          "PropertyType": "[]Topping"
        }
      ]
    },
    {
      "Name": "Topping",
      "Parent": "",
      "Properties": [
        {
          "RemoteIdentifier": "Name",
          "PropertyType": "string"
        }
      ]
    }
  ]
}`

var listViewTemplateContents string = `{{$orig := .}}
{{range .Models}}
<<levo filename:{{.Name}}ListView.java>>
package {{$orig.PackageName}}.models;

{{if hasListType .}}import java.util.List;
{{end}}!>
import com.google.gson.annotations.SerializedName;

public abstract class Abs{{.Name}} {{if ne .Parent ""}}extends {{.Parent}}{{end}} {
   protected static class Fields {
   {{range $index, $prop := .Properties}}!>
      public static final String {{snakecase $prop.RemoteIdentifier | upper}} = "{{$prop.RemoteIdentifier}}";
   {{end}}!>
   }

   {{range $index, $prop := .Properties}}!>
   @SerializedName(Fields.{{snakecase $prop.RemoteIdentifier | upper}}) private {{$prop.PropertyType}} m{{camelcase $prop.LocalIdentifier}};
   {{end}}!>

   {{range $index, $prop := .Properties}}!>
   public {{$prop.PropertyType}} get{{camelcase $prop.LocalIdentifier}}() {
      return m{{camelcase $prop.LocalIdentifier}};
   }

   public void set{{camelcase $prop.LocalIdentifier}}(final {{$prop.PropertyType}} {{$prop.LocalIdentifier}}) {
      m{{camelcase $prop.LocalIdentifier}} = {{$prop.LocalIdentifier}};
   }
   {{end}}!>
}
<<levo>>
{{end}}
`

var descriptionViewTemplateContents string = `{{$orig := .}}
{{range .Models}}
<<levo filename:{{.Name}}DescriptionView.java>>
package {{$orig.PackageName}}.models;

{{if hasListType .}}import java.util.List;
{{end}}!>
import com.google.gson.annotations.SerializedName;

public abstract class Abs{{.Name}} {{if ne .Parent ""}}extends {{.Parent}}{{end}} {
   protected static class Fields {
   {{range $index, $prop := .Properties}}!>
      public static final String {{snakecase $prop.RemoteIdentifier | upper}} = "{{$prop.RemoteIdentifier}}";
   {{end}}!>
   }

   {{range $index, $prop := .Properties}}!>
   @SerializedName(Fields.{{snakecase $prop.RemoteIdentifier | upper}}) private {{$prop.PropertyType}} m{{camelcase $prop.LocalIdentifier}};
   {{end}}!>

   {{range $index, $prop := .Properties}}!>
   public {{$prop.PropertyType}} get{{camelcase $prop.LocalIdentifier}}() {
      return m{{camelcase $prop.LocalIdentifier}};
   }

   public void set{{camelcase $prop.LocalIdentifier}}(final {{$prop.PropertyType}} {{$prop.LocalIdentifier}}) {
      m{{camelcase $prop.LocalIdentifier}} = {{$prop.LocalIdentifier}};
   }
   {{end}}!>
}
<<levo>>
{{end}}
`

var readmeContents string = `
For a quick tutorial on how to use this tool, go to https://github.com/cfmobile/levolib

For documentation, go to https://github.com/cfmobile/levolib/wiki
`

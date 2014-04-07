# Levo
The CLI interface to the [Levolib package](https://github.com/cfmobile/levolib)

- [Installation from Source](https://github.com/cfmobile/levo#installation-from-source)
- [Installing Binaries](https://github.com/cfmobile/levo#installing-binaries)
- [License](https://github.com/cfmobile/levo#license)
- [Quick Tutorial](https://github.com/cfmobile/levo#quick-tutorial)
	1. [Get the example files](https://github.com/cfmobile/levo#1-get-the-example-files)
	2. [Generate a single model](https://github.com/cfmobile/levo#2-generate-a-single-model)
	3. [Generate a model from a schema](https://github.com/cfmobile/levo#3-generate-a-model-from-a-schema)
	4. [Add a new model](https://github.com/cfmobile/levo#4-add-a-new-model)
	5. [Run Levo again](https://github.com/cfmobile/levo#5-run-levo-again)
	6. [Learn more](https://github.com/cfmobile/levo#6-learn-more)
- [File Definition](https://github.com/cfmobile/levo#file-definition)

## Installation from Source
Having first [setup a Go workspace](http://golang.org/doc/code.html), simply run `go get github.com/cfmobile/levo`. Then compile the project (`go build`) and add `levo` to your PATH.

## Installing Binaries
Go to this github repo's [**Releases**](https://github.com/cfmobile/levo/releases) page and download the binary that matches your platform.

## License
This library is licensed under the Apache License, Version 2.0 [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

## Quick Tutorial
Levo blends templates and models to quickly create boilerplate code so you don't have to.

A quick example: I want to create an app for my local pizza shop. The app should use a 3rd party library, but implementing that library involves writing a bunch of boring boilerplate code. Luckily, the library comes with a set of Levo templates, so I can use the `levo` command to generate that boilerplate code automatically.

Let's take a look at how this all works:

### 1. Get the example files
Once you have Levo installed, simply run `levo -example` to output a directory full of example files. These files act as a good starting point for creating a new project using Levo. We'll look at the files in detail later on. For now...

### 2. Generate a single model
You can use the `-model` flag to define a model right on the command line. The `-model` flag accepts a short-form model definition that should be familiar to Rails developers. The format is:

modelName\<*sep*\>RemoteIdentifier\<*sep*\>PropertyType[\<*sep*\>RemoteIdentifier\<*sep*\>PropertyType...]

where \<*sep*\> is any character that is not A-z, 0-9, or _. For example, `-model Product:Name:string:Price:string` or `-model "Product Name:string Price:string"`. **NOTE: If you use spaces as separators, you MUST wrap the whole string in quotes.**

Try generating source code using a model defined on the command line. You'll also need to tell Levo which template you want it to use; pick one from the "templates" directory. For example, you might enter `levo -template templates/_Name_ListView.java -model Beverage:Name:string:Price:float`

A directory called `levo_gen` has been created and it contains a new file! Compare that file with its template counterpart to see how your model was fused with the template to create source code.

### 3. Generate a model from a schema
Instead of writing out the model each time, you can tell Levo to process a model that is defined in a schema file. Select a template from the template directory and a model from schema.json. Then run Levo with the path to the template, the path to the schema, and the name of the model as arguments. For example, you might run `levo -template templates/_Name_ListView.java -schema schema.json -name Pizza`. This will generate a new file in the `levo_gen` directory. Take a look.

### 4. Add a new model
Now let's add a new model to the schema. Make sure you add it to the **Models** array. It can have any content you want, so long as it has a **Name** and **Properties** and it's properties have an **RemoteIdentifier** and a **PropertyType**. Here's an empty template you can use.
```json
{
  "Name": "",
  "Parent": "",
  "Properties": [
    {
      "RemoteIdentifier": "",
      "PropertyType": ""
    }
  ]
}
```

### 5. Run Levo again
Finally, run `levo -template templates/_Name_ListView.java -schema schema.json -name yourModelName`. Your new model has been used to generate a new source code file in the `levo_gen` directory.

### 6. Learn more
That's the end of this tutorial. You can discover more details about these Levo features, and learn about Levo's advanced features, by reading through the documentation in our [github wiki](https://github.com/aaronjarecki/levo/wiki).

## File Definition

#### schema.json
The schema is where we define what makes our app different than any other app using the same templates.
```json
{
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
}
```
The schema defines a **Project Name** and a list of **Models**. All models *must* have a **Name** and **Properties**. They may also define a **Parent**, which they inherit properties from. Each **Property** *must* have an **RemoteIdentifier** and a **PropertyType**. Templates may use both the name of the project and the information about it's models to produce source code files.

#### templates
Templates are stored in a single directory and always have an extension beginning with 'lt'. A single template might produce many source code files. To do this, the output of the template will include headers and footers that divide the output into pieces. The contents of individual files are each wrapped in a header and footer. Headers and footers start with `<<levo`. The header also specifies what the name of the file should be, so a full header looks like `<<levo filename:somefilename.ext>>`.

Here's an example of a template that produces a separate file for each model. Note the header and footer sit within a `range` loop, and have the contents of the file between them.
```java
{{$orig := .}}
{{range .Models}}
<<levo filename:{{.Name}}.java>>
package {{$orig.PackageName}}.models;

{{if hasListType .}}import java.util.List;
{{end}}!>
import com.google.gson.annotations.SerializedName;

public abstract class Abs{{.Name}} {{if ne .Parent ""}}extends {{.Parent}}{{end}} {
   protected static class Fields {
   {{range $index, $prop := .Properties}}!>
      public static final String {{underscoreUppercase $prop.RemoteIdentifier}} = "{{$prop.RemoteIdentifier}}";
   {{end}}!>
   }

   {{range $index, $prop := .Properties}}!>
   @SerializedName(Fields.{{underscoreUppercase $prop.RemoteIdentifier}}) private {{$prop.PropertyType}} m{{camelcase $prop.LocalIdentifier true}};
   {{end}}!>

   {{range $index, $prop := .Properties}}!>
   public {{$prop.PropertyType}} get{{camelcase $prop.LocalIdentifier true}}() {
      return m{{camelcase $prop.LocalIdentifier true}};
   }

   public void set{{camelcase $prop.LocalIdentifier true}}(final {{$prop.PropertyType}} {{$prop.LocalIdentifier}}) {
      m{{camelcase $prop.LocalIdentifier true}} = {{$prop.LocalIdentifier}};
   }
   {{end}}!>
}
<<levo>>
{{end}}
```
At this time, Levo has only one template adapter, able to process templates written in Go. Additional adapters can be written for the tool, making it easy to port templates.

All templates have access to **Project Name**, **Package Name**, an **array of Models**, and **helper functions**. Looking at line *2* and *4* above show examples of how `{{.Models}}` and `{{.PackageName}}` can be used. Line *6* shows a helper function called `{{hasListType}}` being used.

#### config.json
Instead of running the `levo` command many times to produce many source code files, you can pass it a configuration file. The config tells Levo what to do. Most importantly, it indicates the **location of the schema file**, it indicates the **path to the templates directory**, and it provides a **mapping** between models and the templates they should populate. The config also provides a version number and, if applicable, a language and base package.
```json
{
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
        "_Name_ListView.java"
      ]
    },
    {
      "ModelNames": [
        "Pizza",
        "Topping"
      ],
      "TemplateNames": [
        "_Name_DescriptionView.java"
      ]
    }
  ]
}
```

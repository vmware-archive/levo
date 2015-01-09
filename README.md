# Levo
> Code generation for any platform.

# Installation

## From Source
After you've [setup your GO workspace](http://golang.org/doc/code.html) run the following commands:
```
go get github.com/cfmobile/levo
go build
```

## Homebrew

```bash
brew tap cfmobile/homebrew-tap
brew install levo
```

# Simple Usage 

```bash
levo -t <templates> -m <model>
```

e.g. `levo -t template.lt -m "User id:long name:string age:int"`

# Template Example

In the example above, the contents of template.lt looks something like this:

```java
{{$package := .PackageName}}
{{$path := .PackagePath}}
{{range .Models}}
<<levo filename:{{titlecase .Name}}.java directory:src/{{$path}}/models>>
package {{$package}}.models;

import java.util.ArrayList;

public class {{titlecase .Name}} {{if .Parent}}extends {{.Parent}} {{end}}{

	public static class List extends ArrayList<{{titlecase .Name}}> {
		private static final long serialVersionUID = 1L;
	}

	{{range .Properties}}!>
	private {{toJavaType .}} m{{titlecase .LocalIdentifier}};
	{{end}}!>

	{{range .Properties}}!>
	public {{toJavaType .}} get{{titlecase .LocalIdentifier}}() {
		return m{{titlecase .LocalIdentifier}};
	}

	public void set{{titlecase .LocalIdentifier}}({{toJavaType .}} {{camelcase .LocalIdentifier}}) {
		m{{titlecase .LocalIdentifier}} = {{camelcase .LocalIdentifier}};
	}

	{{end}}!>
}
<<levo>>
{{end}}
```

# Existing templates

- [Arca Android](https://github.com/cfmobile/arca-android-templates)
- [Arca iOS](https://github.com/cfmobile/arca-ios-templates)
- [Rails](https://github.com/cfmobile/rails-scaffold-templates)

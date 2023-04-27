# go-dyncapnp

go-dyncapnp is a Go library that provides a wrapper around the original capnproto C++ library with CGO. It is designed to parse .capnp schema files on-the-go, without requiring re-compilation. This makes it easier to work with capnproto schemas in Go projects, especially when the schema is being actively developed and updated.

## Features

- Written in Go, with CGO bindings to the original Cap'n Proto C++ library
- Parses .capnp schema files on-the-go, without re-compilation
- Capable of encoding/decoding packed/plain Cap'n Proto format binary based on the schema
- Catches C++ exceptions and emits them in Go style

## Installation

To use go-dyncapnp in your project, you can simply run:

```
go get github.com/devsisters/go-dyncapnp
```

## Usage

To use go-dyncapnp, you need to provide a map of files containing capnproto schema definitions, along with any necessary imports and the paths to the schema files. You can then call the `ParseFromFiles` function to parse the schemas and obtain a map of `ParsedSchema` objects.

```go
package main

import (
    "fmt"

    "github.com/devsisters/go-dyncapnp"
)

func main() {
	schemaFiles := map[string][]byte{
		"user.capnp": []byte(`
            @0xe4b3f8bd92dc065d;

            struct User {
                id    @0 :UInt64;
                name  @1 :Text;
                email @2 :Text;
            }
        `),
	}

	schemaImports := map[string][]byte{}

	schemaPaths := []string{"user.capnp"}

	schemas, err := dyncapnp.ParseFromFiles(schemaFiles, schemaImports, schemaPaths)
	if err != nil {
		panic(err)
	}

	userFile := schemas["user.capnp"]
	userSchema, err := userFile.Nested("User")
	if err != nil {
		panic(err)
	}

	fmt.Println(userSchema.ShortDisplayName())
	for _, field := range userSchema.AsStruct().Fields() {
		p, err := field.Proto()
		if err != nil {
			panic(err)
		}
		fmt.Printf(" %s @%d %s\n", p.Name(), field.Index(), field.Type().Which())
	}
	// User
	//  id @0 uint64
	//  name @1 text
	//  email @2 text
}
```

## Documentation

For more detailed documentation and additional functions, please visit the [official documentation](https://pkg.go.dev/github.com/devsisters/go-dyncapnp).

## Contributing

Contributions are welcome! If you find any bugs or have any suggestions, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

# Translator

A simple self-hostable translation service, powered by [spaGO](https://github.com/nlpodyssey/spago/).

Some notable features:
- it doesn't rely on external translation services
- no fees, no strings attached
- it's self-contained: just build and run with minimal dependencies
- tiny executable, simple configuration
- REST (OpenAPI) and gRPC API
- it works with [spaGO](https://github.com/nlpodyssey/spago/) models
  - let the program automatically download and convert models from Hugging Face Hub
    (reference: [spaGO Hugging Face Importer](https://github.com/nlpodyssey/spago/tree/main/cmd/huggingfaceimporter))
  - and/or just provide your own models

## Supported models

This project uses [spaGO](https://github.com/nlpodyssey/spago/) machine-learning/NLP
library behind the hood. At present, BART and Marian models for conditional
generation are supported. For more information please refer to
[spaGO BART Machine Translation README section](https://github.com/nlpodyssey/spago/blob/68fc3365bc894f666abd5327cf51eca4964df66d/cmd/bart/README.md#machine-translation).

## Build and Run

The primary intended usage is to run it as a standalone program.
You can get the code and build it like this:

```shell
git clone https://github.com/SpecializedGeneralist/translator.git
cd translator
go mod download
go build -o translator cmd/translator/main.go
```

The `translator` program requires a configuration file to run.
Please refer to the file `sample-configuration.yaml` included with this
project to see an example.

Once you are done with your configuration definition, run:

```shell
./translator -c your-config.yaml
```

The program will first load the configured models from the given path.
If a model is not found, the program will automatically attempt to download it from
Hugging Face models hub, convert it to a spaGO model, and load it as well.

Eventually, the server will start and will be ready to accept requests.
The configured endpoint can be used indifferently for REST (OpenAPI-defined) requests,
or as gRPC service.

The folder `pkg/api` from this project provides the OpenAPI definition file (`api.yaml`)
and also protobuf and gRPC-related definitions and code.

## Use as Go package

This project is a Go module, so you can get and use it from your own code:

```shell
go get -u github.com/SpecializedGeneralist/translator
```

For example, a typical scenario is to import and use the included gRPC client:

```go
import "github.com/SpecializedGeneralist/translator/pkg/api"

// ...

client := api.NewApiClient(conn)

// ...
```

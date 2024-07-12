module github.com/aserto-dev/azm

go 1.22

toolchain go1.22.5

// replace github.com/aserto-dev/go-directory => ../go-directory

require (
	github.com/antlr4-go/antlr/v4 v4.13.1
	github.com/aserto-dev/errors v0.0.8
	github.com/aserto-dev/go-directory v0.31.5
	github.com/deckarep/golang-set/v2 v2.6.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/magefile/mage v1.15.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/nsf/jsondiff v0.0.0-20230430225905-43f6cf3098c1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.33.0
	github.com/samber/lo v1.44.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/text v0.16.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240508200655-46a4cf4ba109.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240624140628-dc46fd24d27d // indirect
)

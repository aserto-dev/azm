module github.com/aserto-dev/azm

go 1.22.9

toolchain go1.23.3

replace github.com/aserto-dev/go-directory => ../go-directory

require (
	github.com/antlr4-go/antlr/v4 v4.13.1
	github.com/aserto-dev/errors v0.0.11
	github.com/aserto-dev/go-directory v0.33.2-0.20241209182914-962e2284a774
	github.com/deckarep/golang-set/v2 v2.6.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/nsf/jsondiff v0.0.0-20230430225905-43f6cf3098c1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.33.0
	github.com/samber/lo v1.47.0
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.68.0
	google.golang.org/protobuf v1.35.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.35.2-20240920164238-5a7b106cbb87.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.24.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20241108190413-2d47ceb2692f // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241118233622-e639e219e697 // indirect
)

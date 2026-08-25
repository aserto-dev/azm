# Private fork of [deckarep/golang-set/v2](https://github.com/deckarep/golang-set)

Removed all the mongodb driver dependencies and implementation

From original ci.yaml

```
go test -v -race ./...
go vet ./...
go test -bench=.
```

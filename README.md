### mitum-currency-extension

*mitum-currency-extension* is the extension of currency model, based on
[*mitum*](https://github.com/spikeekips/mitum) and [*mitum-currency*](https://github.com/spikeekips/mitum-currency).

#### Features,

* account: account address and keypair is not same.
* contract account: account which does not have keys.
* simple transaction: create contract account, deactivate, withdraw.
* *mongodb*: as mitum does, *mongodb* is the primary storage.

#### Installation

> NOTE: at this time, *mitum* and *mitum-currency-extension* is actively developed, so before building mitum-currency-extension, you will be better with building the latest
mitum and mitum-currency source.
> `$ git clone https://github.com/spikeekips/mitum`
> `$ git clone https://github.com/spikeekips/mitum-currency`
>
> and then, add `replace github.com/spikeekips/mitum => <your mitum source directory>` and `replace github.com/spikeekips/mitum-currency => <your mitum-currency source directory>` to `go.mod` of *mitum-currency-extension*.

Build it from source
```sh
$ cd mitum-currency-extension
$ go build -ldflags="-X 'main.Version=v0.0.1'" -o ./mitum-currency-extension ./main.go
```

#### Run

At the first time, you can simply start node with example configuration.

> To start, you need to run *mongodb* on localhost(port, 27017).

```
$ ./mitum-currency-extension node init ./standalone.yml
$ ./mitum-currency-extension node run ./standalone.yml
```

> Please check `$ ./mitum-currency-extension --help` for detailed usage.

#### Test

```sh
$ go clean -testcache; time go test -race -tags 'test' -v -timeout 20m ./... -run .
```

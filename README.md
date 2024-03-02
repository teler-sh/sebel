# sebel

[![GoDoc](https://pkg.go.dev/static/frontend/badge/badge.svg)](http://pkg.go.dev/github.com/teler-sh/sebel)
[![tests](https://github.com/teler-sh/sebel/actions/workflows/tests.yaml/badge.svg)](https://github.com/teler-sh/sebel/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/teler-sh/sebel)](https://goreportcard.com/report/github.com/teler-sh/sebel)

sebel is a Go package that provides functionality for checking SSL/TLS certificates against malicious connections, by identifying and blacklisting certificates used by botnet command and control (C&C) servers.

## Usage

Setting up Sebel instance:

```go
import "github.com/teler-sh/sebel"

// ...

s := sebel.New(Options{/* ... */})
```

> [!NOTE]
> The `Options` parameter is optional. Currently, the only supported option is disabling the SSL blacklist. See [TODO](#TODO).

### Examples

Next, set the transport for the HTTP client you are using:

```go
// initialize Sebel (fetch SSLBL data)
s := sebel.New()

client := &http.Client{
    Transport: s.RoundTripper(http.DefaultTransport),
}

// now, you can use [client.Do], [client.Get], etc. to create requests.

resp, err := client.Get("https://c2.host")
if err != nil && sebel.IsBlacklist(err) {
    // certificate blacklisted
    panic(err)
}
defer resp.Body.Close()
```

Alternatively, for seamless integration without configuring a new client, replace your current default HTTP client with Sebel's `RoundTripper`:

```go
http.DefaultClient.Transport = sebel.New().RoundTripper(http.DefaultTransport)
```

You can also check the certificate later using Sebel's `CheckTLS`.

```go
r, err := http.Get("https://c2.host")
if err != nil {
	panic(err)
}
defer r.Body.Close()

s := sebel.New()

_, err = s.CheckTLS(r.TLS)
if err != nil && sebel.IsBlacklist(err) {
	// certificate blacklisted
	panic(err)
}
```

These examples demonstrate various ways to set up Sebel and integrate it with HTTP clients for SSL/TLS certificate checks.

## TODO

* [ ] Caching SSLBL data under user-specific cache directory.
* [ ] Add `io.Writer` option.
* [ ] ~Add `CheckIP` method.~ Not planned, instead:
* [ ] Add `CheckHost` method.

## Status

> [!CAUTION]
> Sebel has NOT reached 1.0 yet. Therefore, this library is currently not supported and does not offer a stable API; use at your own risk.

There are no guarantees of stability for the APIs in this library, and while they are not expected to change dramatically. API tweaks and bug fixes may occur.

## License

`sebel` is released by [**@dwisiswant0**](https://github.com/dwisiswant0) under the Apache 2.0 license. See [LICENSE](/LICENSE).

The data used in this project are © by [abuse.ch](https://abuse.ch/) under [CC0](https://creativecommons.org/public-domain/cc0/).
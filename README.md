# fasttar [![Go Reference](https://pkg.go.dev/badge/github.com/bored-engineer/fasttar.svg)](https://pkg.go.dev/github.com/bored-engineer/fasttar)
A fast, zero-allocation implementation of Golang's archive/tar

## Security
Some of fasttar's performance gains come from assuming that the input is well-formed. It should never be used in a security sensitive context as it can almost certainly be tricked into parsing data incorrectly or even panicing.

# gobindep

This program lists the dependencies for **Go binaries** that were compiled using **Go modules** support.

Only the dependencies that are actually used by the final executable are included.
This can be compared to the dependencies listed in `go.mod` that might include code that is only used for Tests for example.

## Usage

```bash
gobindep your-go-executable
gobindep --help
```

## Combined with security scanning

This tool can be used for doing a security scan of used dependencies such as **Sonatype nancy**:

* https://github.com/sonatype-nexus-community/nancy

```bash
# Run security check on the dependencies in the go.sum file
$ nancy go.sum
...
Audited dependencies: 65, Vulnerable: 0

# Same thing but use the go introspection tool instead
$ go list -m all | nancy
...
Audited dependencies: 77, Vulnerable: 1

# Finally only check the dependencies we actually use in the final program!
$ gobindep your-go-executable | nancy
...
Audited dependencies: 45, Vulnerable: 0
```

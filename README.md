# vert: A command-line version comparison tool
[![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)
[![Build Status](https://travis-ci.org/Masterminds/vert.svg?branch=master)](https://travis-ci.org/Masterminds/vert)


vert (Version Tester) is a simple command line tool for comparing two or
more versions, or testing versions against fuzzy version rules.

## Basic Usage

Vert takes at least two arguments: A base version, and one or more
versions to compare to the base. The base version is special, in that it
can be a fuzzy version specification rather than an exact version.

```
$ vert 1.0.0 1.0.0
1.0.0
$ echo $?
0
```

When `vert` runs, it will print a normalized string of any version that
matches, and then will set the exit code to the number of version match
failures that it saw.

`vert` understands SemVer v2 versions, so the following will also pass:

```
$ vert v1.0.0 1
1.0.0
$ echo $?
0
```

There are three things to note:

- A leading `v` is ignored.
- Numbers are expanded to a full SemVer string, thus `1` is expanded to
  `1.0.0` and `1.1` is expaneded to `1.1.0`
- The output is normalized to the form `X.Y.Z[-PRERELEASE][+BUILD]`

A failed comparison looks like this:

```
$ vert 1.0.0 1.2.0
$ echo $?
1
```

Failed version comparisons to not print any text unless the given base
version is malformed:

```
$ vert 1.zoo.cheese 1.1.1
Could not parse constraint 1.zoo.cheese
```

Base versions can be fuzzy:

- `vert ">1.0" 1.1`
- `vert "^2" 2.1.3`
- `vert ">1.1.2,<1.3.4" 1.2`

And `vert` understands alpha/beta markers:

```
vert ">1.0.0-alpha.1" 1.0.0-beta.1
1.0.0-beta.1
```

Multiple versions can be compared at once, and using the `-s` flag, you
can even sort the output:

```
$ vert ^1 1.1.1 1.0.1 1.2.3 1.0.2 0.9 2.0
1.1.1
1.0.1
1.2.3
1.0.2
$ echo $?
2
```

In the above, we asked vert for all of the version in the `1.X.Y` range
(`^1`), and then gave it a list of versions, including some outside of
that range.

`vert` returned a list of versions that match. Via the return code, we
can see that two failed to match. To see which failed, we can use the
`-f` flag:

```
$ vert -f ^1 1.1.1 1.0.1 1.2.3 1.0.2 0.9 2.0
0.9.0
2.0.0
```

We can sort output using the `-s` flag:

```
vert -s ^1 1.1.1 1.0.1 1.2.3 1.0.2 0.9 2.0
1.0.1
1.0.2
1.1.1
1.2.3
```

Finally, `vert` can transform `git describe` versions into SemVer,
assuming the Git tags are SemVer:

```
$ vert -g ^1 $(git describe --tags)
1.0.1+32.fef45
```

In the future, we'd like to add more transformations. If you have any
ideas, please let us know in the issue queue.

## Installation

Assuming you have make, [Go](http://golang.org) version 1.5.1 or later and
[Glide](https://github.com/Masterminds/glide) version 0.8.2 or greater, you can
simply run `make`:

```
$ make all
glide install
[INFO] Fetching updates for github.com/Masterminds/semver.
[INFO] Fetching updates for github.com/codegangsta/cli.
[INFO] github.com/Masterminds/semver is already set to version 6333b7bd29aad1d79898ff568fd90a8aa533ae82. Skipping update.
[INFO] github.com/codegangsta/cli is already set to version c31a7975863e7810c92e2e288a9ab074f9a88f29. Skipping update.
[INFO] Setting version for github.com/Masterminds/semver to 6333b7bd29aad1d79898ff568fd90a8aa533ae82.
[INFO] Setting version for github.com/codegangsta/cli to c31a7975863e7810c92e2e288a9ab074f9a88f29.
go test .
ok  	github.com/technosophos/vert	0.006s
go build -o vert vert.go
install -d /usr/local/bin/
install -m 755 ./vert /usr/local/bin/vert
```

This will install into `/usr/local/bin/vert`.

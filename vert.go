package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/urfave/cli"
)

const name = "vert"
const version = "0.1.0"
const description = `Version Tester. Compare versions.

Vert is a tool for comparing two version strings, or comparing a version string
to a version range.

	$ vert ">1.0.0" 1.1.0
	1.1.0

	$ vert "<1.0.0" 1.1.0
	   # No output because nothing matched.

	$ vert -f "<1.0.0" 1.1.0
	1.1.0   # -f returns all failures, rather than matches.

	$ vert ">1.0.0" 1.1.0 1.1.1 1.2.3 0.1.1
	1.1.0
	1.1.1
	1.2.3

Vert can also convert a Git version to a SemVer 2 version.

	$ vert -g ">1" v1.10.0-123-g0239788                                                                                                                                                          1 ↵
	1.10.0+123.g0239788

Note that it assigns the git commit count and hash to the build metadata, not
to a pre-release tag.

See below for information about how to determine the number of failed tests.

EXIT CODES:

vert returns exit codes based on the number of failed matches. There are a few
reserved exit codes:

- 128: The command was not called correctly.
- 256: A version failed to parse, and comparisons could not continue. This will
  occur if the original constraint version cannot be parsed. If any
  subsequent version fails to parse, it will simply be counted as a failure.

Any other error codes indicate the number of failed tests. For example:

	$ vert 1.2.3 1.2.3 1.2.4 1.2.5
	1.2.3
	$ echo $?
	2   # <-- Two tests failed.

BASE VERSIONS:

The base version may be in any of the following formats:

- An exact semantic version number
	- 1.2.3
	- v1.2.3
	- 1.2.3-alpha.1+10212015
- A semantic version range
	- *
	- !=1.0.0
	- >=1.2.3
	- >1.2.3,<1.3.2
	- ~1.2.0
	- ^2.3

VERSIONS:

Other than the base version, all other supplied versions must follow the
SemVer 2 spec. Examples:

	- 1.2.3
	- v1.2.3
	- 1.2.3-alpha.1+10212015
	- v1.2.3-alpha.1+10212015
	- 1 (equivalent to 1.0.0)
`

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = description
	app.Action = func(c *cli.Context) { res := run(c); os.Exit(res) }
	app.Version = version
	app.ArgsUsage = "BASE VERSION [VERSION [VERSION [...]]"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "failed,f",
			Usage: "Show the versions that failed rather than the ones that passed.",
		},
		cli.BoolFlag{
			Name:  "sort,s",
			Usage: "Sort the versions before printing. Without this, versions are returned in the order they were tested.",
		},
		cli.BoolFlag{
			Name:  "git,g",
			Usage: "Assume that (non-base) versions are in Git `git describe --tags` version format, and convert to SemVer.",
		},
	}
	app.Run(os.Args)
}

// context describes the relevant portion of a cli.Context.
//
// This abstraction makes mocking easy.
type context interface {
	Bool(string) bool
	Args() cli.Args
}

// run handles all of the flags and then runs the main action.
func run(c context) int {
	args := c.Args()
	if len(args) < 2 {
		perr("Not enough arguments")
		return 128
	}

	if c.Bool("git") {
		for i := 1; i < len(args); i++ {
			nv, err := git2semver(args[i])
			if err != nil {
				perr("Not a recognize git version: %s", args[i])
				continue
			}
			args[i] = nv.String()
		}
	}

	pass, fail, code := compare(args[0], args[1:])

	out := pass
	if c.Bool("failed") {
		out = fail
	}

	if c.Bool("sort") {
		sort.Sort(semver.Collection(out))
	}

	pvers(out)
	return code
}

// compare compiles a base version comparator, and then compares all cases to it.
//
// It retuns an array of versions that passed, and an array of versions that failed.
func compare(base string, cases []string) ([]*semver.Version, []*semver.Version, int) {
	passed, failed := []*semver.Version{}, []*semver.Version{}

	constraint, err := semver.NewConstraint(base)
	if err != nil {
		perr("Could not parse constraint %s", base)
		return passed, failed, 128
	}

	for _, t := range cases {
		ver, err := semver.NewVersion(t)
		if err != nil {
			failed = append(failed, ver)
			perr("Failed to parse %s", t)
			continue
		}
		if constraint.Check(ver) {
			passed = append(passed, ver)
			continue
		}
		failed = append(failed, ver)
	}

	return passed, failed, len(failed)
}

var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

// pvers prints a list of versions to standard out.
func pvers(vers []*semver.Version) {
	for _, v := range vers {
		fmt.Fprintln(stdout, v.String())
	}
}

// pout prints to stdout.
func pout(msg string, args ...interface{}) {
	fmt.Fprintf(stdout, msg, args...)
	fmt.Fprintln(stdout)
}

// perr prints to stderr.
func perr(msg string, args ...interface{}) {
	fmt.Fprintf(stderr, msg, args...)
	fmt.Fprintln(stderr)
}

// git2semver converts a Git version to a SemVer 2 version
//
// This assumes that the base tag is a semver tag.
//
// v1.2.3-3-afeee becomes 1.2.3+3.afeee
func git2semver(ver string) (*semver.Version, error) {
	va := strings.Split(ver, "-")
	target := va[0]
	if len(va) > 1 {
		md := strings.Join(va[1:], ".")
		target += "+" + md
	}
	return semver.NewVersion(target)
}

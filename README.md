# Introduction
CLI for deleting local git branches. Requires [git CLI](https://git-scm.com/downloads) to be installed.

# Usage
Navigate to a directory where git branches are stored, and execute git-delete-branches

Usage mimics typical Linux commands, like ls. empty argument to select all, and filter by * wildcard.

```
$ git-delete-branches branch*

$ git-delete-branches --help
```

# Build
Requires Go 1.12.9 or higher.


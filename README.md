# Github CLI

A simple tool to list Github trending, issues, pull requests.

#### Usage

```bash
# trending
github-cli trending
github-cli trending go
github-cli trending go --time weekly

# pull
github-cli pulls tikv/tikv
github-cli pull tikv/tikv 12511

# issue
github-cli issues tikv/tikv
github-cli issue tikv/tikv 12509
```

#### Thanks
modified on https://github.com/siddontang/github-cli, with changes:
- config file was not required by default
- use go module
- update trending parsing

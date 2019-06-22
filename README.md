[![GoReportCard](http://goreportcard.com/badge/github.com/asticode/go-astislack)](http://goreportcard.com/report/github.com/asticode/go-astislack)
[![GoDoc](https://godoc.org/github.com/asticode/go-astislack?status.svg)](https://godoc.org/github.com/asticode/go-astislack)

Clear your Slack history easily!

Thanks to `astislack`, mass delete all your Slack messages and files!

# Releases

Check out the [release page](https://github.com/asticode/go-astislack/releases) to download a binary compatible with your OS as well as the `local.toml.dist` file.

# Usage

1) Create a Slack legacy token [here](https://api.slack.com/custom-integrations/legacy-tokens).

2) Copy `local.toml.dist` to `local.toml`

3) Edit `local.toml` and replace `SLACK_LEGACY_TOKEN` with your legacy token.

4) Run

```
$ /path/to/binary -c local.toml -v
```

This should start mass deleting all your Slack messages and files!
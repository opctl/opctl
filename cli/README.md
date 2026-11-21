# Problem Statement
CLI for [opctl](https://opctl.io).


## Interactive containers

When stdin and stdout are terminals and opctl is in the foreground, `opctl run`
forwards input to the most recently started container. Docker and the embedded
runtime support this mode. Press `Ctrl-]` to cycle through running containers and
then release focus. Press it again to focus the most recently started container.
If a focused container exits or its input connection closes, focus moves to the
most recent container that can still accept input.

`Ctrl-C` goes to the focused container. With no container focused, the live call
graph resumes unless you passed `--no-progress`. Press `Ctrl-C` to stop the op
gracefully, or press it again to force the CLI to exit.

When input disconnects, opctl sends Ctrl-D. Programs reading in canonical mode
with the default EOF key get end-of-input. Raw-mode programs receive the literal
bytes and may keep running until the op is killed.

Inside an interactive container, use `timeout --foreground` so the child can use
the terminal. A background child can receive `SIGTTOU` and stop before the timeout
can kill it.

# Support
join us on
[![Slack](https://img.shields.io/badge/slack-opctl-E01563.svg)](https://join.slack.com/t/opctl/shared_invite/zt-51zodvjn-Ul_UXfkhqYLWZPQTvNPp5w)
or [open an issue](https://github.com/opctl/opctl/issues)


# Contributing
see [CONTRIBUTING.md](CONTRIBUTING.md)

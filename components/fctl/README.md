# FCTL

## Add commands

### Write usage

Follow guidelines : http://docopt.org/


# TUI (Text User Interface) using Bubble Tea

## Debug mode

```bash
# 1. Run the program
> go run ./ stacks list -o dynamic

# Use `log := helpers.NewLogger("<prefix>")`
# Surroung your logs and happy debugging
# log.Log("Hello world") must be a string.
# It logs to a file named ./debug.log in the current directory

> tail -f debug.log

```

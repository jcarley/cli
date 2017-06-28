# Automatic Updates

Any version of the CLI before version 4.0.0 will no longer automatically update. In order to obtain version 4.0.0 you will have to manually install the CLI. Automatic updates will work again once you have upgraded to 4.0.0 or greater.

Once downloaded, the CLI will attempt to automatically update itself when a new version becomes available. This ensures you are always running a compatible version of the Datica CLI. However you can always check out the latest releases on the [releases page](https://github.com/daticahealth/cli/releases).

To ensure your CLI can automatically update itself, be sure to put the binary in a location where you have **write access** without the need for sudo or escalated privileges.

# Supported Platforms

Since version 2.0.0, the following platforms and architectures are supported by the Datica CLI.

| OS | Architecture |
|----|--------------|
| Darwin (Mac OS X) | 64-bit |
| Linux | 64-bit |
| Windows | 64-bit |

# Global Scope

Datica CLI commands can be run anywhere on your system with two exceptions. The [`datica init`](#init) command expects to be inside of either a git repository or a directory that you intend to be a git repository, as it will set up a git repository if one does not exist. Additionally, the [`datica git-remote`](#git-remote) command is used to manage git remotes and must be run inside of a git repository in order to work.

If you have more than one environment, you must specify which environment to use with the global `-E` flag.

Let's say you have initialized a code project for two environments (ex. "staging" and "production") named `mysandbox` and `myprod`. You have two options to specify which environment to run a command against.

First, you can tell the CLI which environment you want to use with the global option `-E` or `--env` (see [Global Options](#global-options)). Your command might start like this

```
datica -E myprod ...
```

If you don't set the `-E` flag, then the CLI picks one of your environments and prompts you to continue with this environment. This concept of scope will make it easier for Datica customers with multiple environments to use the CLI!

# Bash Autocompletion

One feature we've found helpful on \*Nix systems is autocompletion in bash. To enable this feature, head over to the github repo and download the `datica_autocomplete` file. If you use a Mac, you will need to install bash-completion with `brew install bash-completion` or `source` the `datica_autocomplete` file each time you start up a terminal. Store this file locally in `/etc/bash_completion.d/` or (`/usr/local/etc/bash_completion.d/` on a Mac). Completion will be available when you restart your terminal. Now type `datica ` and hit tab twice to see the list of available commands. **Please note** that autocompletion only works one level deep. The CLI will not autocomplete or suggest completions when you type `datica db ` and then hit tab twice. It currently only works when you have just `datica ` typed into your terminal. This is a feature we are looking into expanding in the future.

Note: you may have to add `source /etc/bash_completion.d/datica_autocomplete` (`/usr/local/etc/bash_completion.d/datica_autocomplete`) in your `~/.bashrc` (`~/.bash_profile`) file.

# Global Options

The following table outlines all global options available in the CLI. Global options are always set after the word `datica` and before any commands. Rather than setting these each time, you may also set an environment variable with the appropriate value which will automatically be used.

| Short Name | Long Name | Description | Environment Variable |
|------------|-----------|-------------|----------------------|
| &nbsp; | --email | Your Datica email that you login to the Dashboard with | DATICA_EMAIL |
| -U | --username | [DEPRECATED] Your Datica username that you login to the Dashboard with. Please use --email instead | DATICA_USERNAME |
| -P | --password | Your Datica password that you login to the Dashboard with | DATICA_PASSWORD |
| -E | --env | The name of the environment for which this command will be run. | DATICA_ENV |

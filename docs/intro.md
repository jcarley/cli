# Automatic Updates

Once downloaded, the CLI will attempt to automatically update itself when a new version becomes available. This ensures you are always running a compatible version of the Catalyze CLI. However you can always check out the latest releases on the [releases page](https://github.com/catalyzeio/cli/releases).

To ensure your CLI can automatically update itself, be sure to put the binary in a location where you have **write access** without the need for sudo or escalated privileges.

# Supported Platforms

Since version 2.0.0, the following platforms and architectures are supported by the Catalyze CLI.

| OS | Architecture |
|----|--------------|
| Darwin (Mac OS X) | 64-bit, 32-bit |
| Linux | 64-bit, 32-bit |
| Windows | 64-bit, 32-bit |

# Global Scope

The CLI now supports the concept of scope. Previous to version 2.0.0, all commands had to be run within an associated local git repo. Now, the only time you need to be in a local git repo is when you associate to a new environment. After the initial association, CLI commands can be run from any directory. If you have more than one environment, you must specify which environment to use with the global `-E` flag.

Let's say you have associated to two environments named `mysandbox` and `myprod`. You have two options to specify which environment to run a command against.

First, you can tell the CLI which environment you want to use with the global option `-E` or `--env` (see [Global Options](#global-options)). Your command might start like this

```
catalyze -E myprod ...
```

If you don't set the `-E` flag, then the CLI takes the first environment you associated and prompts you to continue with this environment. This concept of scope will make it easier for Catalyze customers with multiple environments to use the CLI!

# Environment Aliases


When you associate an environment from within a local git repo, you typically run the following command:

```
catalyze associate "My Health Tech Company Production" app01
```

Where `My Health Tech Company Production` is the name of your environment. However with the concept of [scope](#global-scope) and being able to specify which environment to use on a command by command basis with the `-E` global option, that is a lot to type! This is where environment aliases come in handy.

When you associate an environment and you want to pick a shorter name to reference the environment by, simply add a `-a` flag to the command. Let's try the command again calling it `prod` this time:

```
catalyze associate "My Health Tech Company Production" app01 -a prod
```

Now when you run the [associated](#associated) command, you will see the alias as well as the actual environment name.

When using aliases, there are a couple things to keep in mind. Aliases are only local and never leave your local machine. If you alias this environment `prod`, a coworker can alias the environment `healthtech-prod` with no ramifications. Second, after setting an alias you will never reference the environment by its actual name with the CLI. You will always use the alias for flags, arguments, options, etc.

To change or remove an alias, you must [disassociate](#disassociate) and then [reassociate](#associate) with a new alias.

# Bash Autocompletion

One feature we've found helpful on \*Nix systems is autocompletion in bash. To enable this feature, head over to the github repo and download the `catalyze_autocomplete` file. If you use a Mac, you will need to install bash-completion with `brew install bash-completion` or `source` the `catalyze_autocomplete` file each time you start up a terminal. Store this file locally in `/etc/bash_completion.d/` or (`/usr/local/etc/bash_completion.d/` on a Mac). Completion will be available when you restart your terminal. Now type `catalyze ` and hit tab twice to see the list of available commands. **Please note** that autocompletion only works one level deep. The CLI will not autocomplete or suggest completions when you type `catalyze db ` and then hit tab twice. It currently only works when you have just `catalyze ` typed into your terminal. This is a feature we are looking into expanding in the future.

Note: you may have to add `source /etc/bash_completion.d/catalyze_autocomplete` (`/usr/local/etc/bash_completion.d/catalyze_autocomplete`) in your `~/.bashrc` (`~/.bash_profile`) file.

# Global Options

The following table outlines all global options available in the CLI. Global options are always set after the word `catalyze` and before any commands. Rather than setting these each time, you may also set an environment variable with the appropriate value which will automatically be used.

| Short Name | Long Name | Description | Environment Variable |
|------------|-----------|-------------|----------------------|
| -U | --username | Your catalyze username that you login to the Dashboard with | CATALYZE_USERNAME |
| -P | --password | Your catalyze password that you login to the Dashboard with | CATALYZE_PASSWORD |
| -E | --env | The local alias of the environment in which this command will be run. Read more about [environment aliases](#environment-aliases) | CATALYZE_ENV |

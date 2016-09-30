# Global Scope

The CLI now supports the concept of scope. Previous to version 2.0.0, all commands had to be run within an associated local git repo. Now, the only time you need to be in a local git repo is when you associate to a new environment. After the initial association, CLI commands can be run from any directory. If you have more than one environment, you must specify which environment to use with the global `-E` flag.

Let's say you have associated to two environments named `mysandbox` and `myprod`. You have two options to specify which environment to run a command against.

First, you can tell the CLI which environment you want to use with the global option `-E` or `--env` (see [Global Options](#global-options)). Your command might start like this

```
catalyze -E myprod ...
```

If you don't set the `-E` flag, then the CLI takes the first environment you associated and prompts you to continue with this environment. This concept of scope will make it easier for Catalyze customers with multiple environments to use the CLI!

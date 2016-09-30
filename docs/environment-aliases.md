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

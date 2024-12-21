# xo - aliasing on steroids

`xo` is a command-line application for managing your projects from anywhere in your terminal by importing projects and adding commands to this app. `xo` becomes your own CLI app which you can customize and add commands for all of your projects. The most time-saving part of this tool is that:

- you can run your commands from any path in your terminal
- you can open your project from any path in your terminal
- you can share xofile with anyone for easy understanding of setup of your project

## Getting started

By default `xo` doesn't see know anything about your workspace. You have to do:

```
xo init
```

which will create a file called `xo.json` with defaults. 

Here is an example `xo.json` file which this project uses:

```
{
    "name": "xoproject",
    "commands": [
        {
            "name": "run",
            "cmd": "go run main.go",
            "env": [],
            "help": "run the xo project"
        },
        {
            "name": "install",
            "cmd": "go build -o xo main.go && sudo mv xo /usr/local/bin/xo && chmod +x /usr/local/bin/xo",
            "help": "build and install the xo binary in your PATH"
        },
        {
            "name": "pack",
            "cmd": "go build -o xo main.go",
            "help": "pack xo binary for release"
        }
    ]
}
```

once you have added your commands you can `import` this project to `xo` by doing:

```
xo import
```

now lets run the info command to see if it got imported:

```
xo ! xoproject
```

the `!` command displays the project details and the `xoproject` after the `!` is the project name. This produces the following output:

```
+---------+----------------------------------------------+
| XO      |                                              |
+---------+----------------------------------------------+
| Command | Help                                         |
+---------+----------------------------------------------+
| install | build and install the xo binary in your PATH |
| pack    | pack xo binary for release                   |
| run     | run the xo project                           |
+---------+----------------------------------------------+
```

once you have imported the project you now can run your project from anywhere:

```
xo xoproject run
```

take a look at this [blog post]() for more info on this tool.


## Future of this tool

`xo` will help developers and teams to setup whole workspace with a single command.
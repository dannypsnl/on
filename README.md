# On

This is a command line tool helps you avoid keep typing the same prefix when you are operating the very similar commands

### Demo:

[![asciicast](https://asciinema.org/a/5s6FVyl2WsU7UeM3BYrfK3gdu.svg)](https://asciinema.org/a/5s6FVyl2WsU7UeM3BYrfK3gdu)

### Key Binds

- `<C-a>`: control + A, type this then write done new contexts and `<Enter>`, you would see the context be apply on to old context

    For example:
    ```
    on(kubectl)>
    # <C-a>get<Enter>
    on(kubectl get)>
    ```
- `<C-c>`: control + C, type this would pop out the last element in command context

    For example:
    ```
    on(kubectl get)>
    # <C-c>
    on(kubectl)>
    ```

At here you can see we avoid typing `kubectl` again and again

### Completion

Completion would find out all command, alias, and executable for command `on`

To install the completion could use:

```
# for zsh
$ make zsh-completion
# or for bash
$ make bash-completion
```

#!/bin/bash

complete -W "$(compgen -abc | grep -Ev "\_")" on

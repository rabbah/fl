#!/bin/bash

# Check the first argument
case "$1" in
    # test all available go modules
    "test")
        pushd src > /dev/null
        go test ./...
        popd > /dev/null
        ;;
    # run with the passed arguments
    # E.g. ./run.sh run -l -v prompt...
    "run")
        pushd src > /dev/null
        shift
        go run fl.go $@
        popd > /dev/null
        ;;
    # run the expected 'ls -l' prompt with ONLY input flags (prompt is provided)
    # E.g. ./run.sh rls -l -v
    "rls")
        pushd src > /dev/null
        PROMPT="show the items in the current directory in long list format without backtics or a reference to the shell"
        shift
        go run fl.go $@ $PROMPT
        popd > /dev/null
        ;;
    *)
        echo "Invalid argument."
        echo "Usage:"
        echo "  $0 {test|run [args]|rls [flags]}"
        exit 1
        ;;
esac

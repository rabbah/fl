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
        PROMPT="show the hidden items in this directory in long list format"
        # the below prompts should respectively yield 'ls -a -l ..' 'ls -a -l ~', but it has trouble executing without error
        # PROMPT="show the hidden items in the directory above in long list format"
        # PROMPT="show the hidden items in my home directory in long list format. don't include backtics or a reference to the shell. Separate flags into multiple flags. Don't use tilde expansion."
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

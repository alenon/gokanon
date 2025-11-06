#!/usr/bin/env bash

# Bash completion script for gokanon

_gokanon_completions() {
    local cur prev words cword
    _init_completion || return

    # Main commands
    local commands="run list compare export stats trend check flamegraph serve delete doctor interactive help"

    # If we're at the first argument, complete commands
    if [ $cword -eq 1 ]; then
        COMPREPLY=($(compgen -W "$commands" -- "$cur"))
        return
    fi

    # Get the main command
    local command="${words[1]}"

    # Command-specific completions
    case "$command" in
        run)
            local run_opts="-bench -pkg -profile -benchtime -count -timeout -cpu -v"
            COMPREPLY=($(compgen -W "$run_opts" -- "$cur"))
            ;;
        compare)
            if [[ "$cur" == -* ]]; then
                COMPREPLY=($(compgen -W "--latest -format" -- "$cur"))
            else
                # Complete with run IDs (would need to call gokanon list)
                COMPREPLY=()
            fi
            ;;
        export)
            if [[ "$prev" == "-format" ]]; then
                COMPREPLY=($(compgen -W "html csv markdown json" -- "$cur"))
            elif [[ "$cur" == -* ]]; then
                COMPREPLY=($(compgen -W "--latest -format -output" -- "$cur"))
            fi
            ;;
        stats|trend)
            COMPREPLY=($(compgen -W "-last -format" -- "$cur"))
            ;;
        check)
            COMPREPLY=($(compgen -W "--latest -threshold -format" -- "$cur"))
            ;;
        serve)
            COMPREPLY=($(compgen -W "-port -open" -- "$cur"))
            ;;
        flamegraph)
            if [[ "$cur" == -* ]]; then
                COMPREPLY=($(compgen -W "-port -open" -- "$cur"))
            fi
            ;;
        delete)
            # Could complete with run IDs
            COMPREPLY=()
            ;;
        *)
            COMPREPLY=()
            ;;
    esac
}

complete -F _gokanon_completions gokanon

#!/usr/bin/env bash

# Bash completion script for gokanon

_gokanon_completions() {
    local cur prev words cword
    _init_completion || return

    # Main commands
    local commands="run list compare export stats trend check flamegraph serve delete baseline doctor interactive completion help"

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
            local run_opts="-bench -pkg -profile -storage -benchtime -count -timeout -cpu -v"
            COMPREPLY=($(compgen -W "$run_opts" -- "$cur"))
            ;;
        compare)
            if [[ "$cur" == -* ]]; then
                COMPREPLY=($(compgen -W "--latest --baseline -storage -format" -- "$cur"))
            else
                # Complete with run IDs (would need to call gokanon list)
                COMPREPLY=()
            fi
            ;;
        export)
            if [[ "$prev" == "-format" ]]; then
                COMPREPLY=($(compgen -W "html csv markdown json" -- "$cur"))
            elif [[ "$cur" == -* ]]; then
                COMPREPLY=($(compgen -W "--latest -format -output -storage" -- "$cur"))
            fi
            ;;
        stats|trend)
            COMPREPLY=($(compgen -W "-last -storage -format" -- "$cur"))
            ;;
        check)
            COMPREPLY=($(compgen -W "--latest -threshold -storage -format" -- "$cur"))
            ;;
        serve)
            COMPREPLY=($(compgen -W "-port -storage -open" -- "$cur"))
            ;;
        flamegraph)
            if [[ "$cur" == -* ]]; then
                COMPREPLY=($(compgen -W "-port -storage -open" -- "$cur"))
            fi
            ;;
        baseline)
            # Handle baseline subcommands
            if [ $cword -eq 2 ]; then
                COMPREPLY=($(compgen -W "save list show delete" -- "$cur"))
            else
                local subcommand="${words[2]}"
                case "$subcommand" in
                    save)
                        COMPREPLY=($(compgen -W "-name -run -desc -storage" -- "$cur"))
                        ;;
                    list)
                        COMPREPLY=($(compgen -W "-storage" -- "$cur"))
                        ;;
                    show)
                        COMPREPLY=($(compgen -W "-name -storage" -- "$cur"))
                        ;;
                    delete)
                        COMPREPLY=($(compgen -W "-name -storage" -- "$cur"))
                        ;;
                esac
            fi
            ;;
        completion)
            if [ $cword -eq 2 ]; then
                COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur"))
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

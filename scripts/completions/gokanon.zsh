#compdef gokanon

# Zsh completion script for gokanon

_gokanon() {
    local -a commands
    commands=(
        'run:Run benchmarks and save results'
        'list:List all saved benchmark results'
        'compare:Compare two benchmark results'
        'export:Export comparison results to various formats'
        'stats:Show statistical analysis of multiple runs'
        'trend:Analyze performance trends over time'
        'check:Check performance against thresholds'
        'flamegraph:View CPU/memory flame graphs'
        'serve:Start interactive web dashboard'
        'delete:Delete a benchmark result'
        'doctor:Run diagnostics'
        'interactive:Start interactive mode'
        'help:Show help message'
    )

    local -a run_opts
    run_opts=(
        '-bench[Benchmark pattern]:pattern:'
        '-pkg[Package pattern]:pattern:'
        '-profile[Enable profiling]:types:(cpu mem cpu,mem)'
        '-benchtime[Benchmark duration]:duration:'
        '-count[Run count]:count:'
        '-timeout[Test timeout]:duration:'
        '-cpu[CPU counts]:counts:'
        '-v[Verbose output]'
    )

    local -a export_formats
    export_formats=(
        'html:HTML format'
        'csv:CSV format'
        'markdown:Markdown format'
        'json:JSON format'
    )

    _arguments -C \
        '1: :->command' \
        '*:: :->args'

    case $state in
        command)
            _describe 'gokanon command' commands
            ;;
        args)
            case $words[1] in
                run)
                    _arguments $run_opts
                    ;;
                compare)
                    _arguments \
                        '--latest[Compare latest two runs]' \
                        '-format[Output format]:format:(table json)'
                    ;;
                export)
                    _arguments \
                        '--latest[Export latest comparison]' \
                        '-format[Export format]:format:->formats' \
                        '-output[Output file]:file:_files'
                    ;;
                stats|trend)
                    _arguments \
                        '-last[Number of runs]:count:' \
                        '-format[Output format]:format:(table json)'
                    ;;
                check)
                    _arguments \
                        '--latest[Check latest two runs]' \
                        '-threshold[Threshold percentage]:threshold:' \
                        '-format[Output format]:format:(table json)'
                    ;;
                serve)
                    _arguments \
                        '-port[Server port]:port:' \
                        '-open[Open browser automatically]'
                    ;;
                flamegraph)
                    _arguments \
                        '-port[Server port]:port:' \
                        '-open[Open browser automatically]'
                    ;;
            esac
            ;;
    esac

    if [[ $state == formats ]]; then
        _describe 'export format' export_formats
    fi
}

_gokanon "$@"

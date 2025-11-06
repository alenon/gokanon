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
        'baseline:Manage baseline benchmarks'
        'doctor:Run diagnostics'
        'interactive:Start interactive mode'
        'completion:Install shell completion scripts'
        'help:Show help message'
    )

    local -a run_opts
    run_opts=(
        '-bench[Benchmark pattern]:pattern:'
        '-pkg[Package pattern]:pattern:'
        '-profile[Enable profiling]:types:(cpu mem cpu,mem)'
        '-storage[Storage directory]:directory:_files -/'
        '-benchtime[Benchmark duration]:duration:'
        '-count[Run count]:count:'
        '-timeout[Test timeout]:duration:'
        '-cpu[CPU counts]:counts:'
        '-v[Verbose output]'
    )

    local -a baseline_subcommands
    baseline_subcommands=(
        'save:Save a benchmark run as a baseline'
        'list:List all saved baselines'
        'show:Show details of a specific baseline'
        'delete:Delete a baseline'
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
                        '--baseline[Compare against baseline]:baseline:' \
                        '-storage[Storage directory]:directory:_files -/' \
                        '-format[Output format]:format:(table json)'
                    ;;
                export)
                    _arguments \
                        '--latest[Export latest comparison]' \
                        '-format[Export format]:format:->formats' \
                        '-output[Output file]:file:_files' \
                        '-storage[Storage directory]:directory:_files -/'
                    ;;
                stats|trend)
                    _arguments \
                        '-last[Number of runs]:count:' \
                        '-storage[Storage directory]:directory:_files -/' \
                        '-format[Output format]:format:(table json)'
                    ;;
                check)
                    _arguments \
                        '--latest[Check latest two runs]' \
                        '-threshold[Threshold percentage]:threshold:' \
                        '-storage[Storage directory]:directory:_files -/' \
                        '-format[Output format]:format:(table json)'
                    ;;
                serve)
                    _arguments \
                        '-port[Server port]:port:' \
                        '-storage[Storage directory]:directory:_files -/' \
                        '-open[Open browser automatically]'
                    ;;
                flamegraph)
                    _arguments \
                        '-port[Server port]:port:' \
                        '-storage[Storage directory]:directory:_files -/' \
                        '-open[Open browser automatically]'
                    ;;
                baseline)
                    case $words[2] in
                        save)
                            _arguments \
                                '-name[Baseline name]:name:' \
                                '-run[Run ID]:run_id:' \
                                '-desc[Description]:description:' \
                                '-storage[Storage directory]:directory:_files -/'
                            ;;
                        list)
                            _arguments \
                                '-storage[Storage directory]:directory:_files -/'
                            ;;
                        show)
                            _arguments \
                                '-name[Baseline name]:name:' \
                                '-storage[Storage directory]:directory:_files -/'
                            ;;
                        delete)
                            _arguments \
                                '-name[Baseline name]:name:' \
                                '-storage[Storage directory]:directory:_files -/'
                            ;;
                        *)
                            _describe 'baseline subcommand' baseline_subcommands
                            ;;
                    esac
                    ;;
                completion)
                    _arguments '1:shell:(bash zsh fish)'
                    ;;
            esac
            ;;
    esac

    if [[ $state == formats ]]; then
        _describe 'export format' export_formats
    fi
}

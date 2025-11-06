# Fish completion script for gokanon

# Main commands
complete -c gokanon -f -n __fish_use_subcommand -a run -d "Run benchmarks and save results"
complete -c gokanon -f -n __fish_use_subcommand -a list -d "List all saved benchmark results"
complete -c gokanon -f -n __fish_use_subcommand -a compare -d "Compare two benchmark results"
complete -c gokanon -f -n __fish_use_subcommand -a export -d "Export comparison results"
complete -c gokanon -f -n __fish_use_subcommand -a stats -d "Show statistical analysis"
complete -c gokanon -f -n __fish_use_subcommand -a trend -d "Analyze performance trends"
complete -c gokanon -f -n __fish_use_subcommand -a check -d "Check performance thresholds"
complete -c gokanon -f -n __fish_use_subcommand -a flamegraph -d "View flame graphs"
complete -c gokanon -f -n __fish_use_subcommand -a serve -d "Start web dashboard"
complete -c gokanon -f -n __fish_use_subcommand -a delete -d "Delete a benchmark result"
complete -c gokanon -f -n __fish_use_subcommand -a baseline -d "Manage baseline benchmarks"
complete -c gokanon -f -n __fish_use_subcommand -a doctor -d "Run diagnostics"
complete -c gokanon -f -n __fish_use_subcommand -a interactive -d "Start interactive mode"
complete -c gokanon -f -n __fish_use_subcommand -a completion -d "Install shell completion scripts"
complete -c gokanon -f -n __fish_use_subcommand -a help -d "Show help message"

# run command options
complete -c gokanon -n "__fish_seen_subcommand_from run" -o bench -d "Benchmark pattern"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o pkg -d "Package pattern"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o profile -d "Enable profiling" -a "cpu mem cpu,mem"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o storage -d "Storage directory" -r
complete -c gokanon -n "__fish_seen_subcommand_from run" -o benchtime -d "Benchmark duration"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o count -d "Run count"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o timeout -d "Test timeout"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o cpu -d "CPU counts"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o v -d "Verbose output"

# compare command options
complete -c gokanon -n "__fish_seen_subcommand_from compare" -l latest -d "Compare latest two runs"
complete -c gokanon -n "__fish_seen_subcommand_from compare" -l baseline -d "Compare against baseline" -r
complete -c gokanon -n "__fish_seen_subcommand_from compare" -o storage -d "Storage directory" -r
complete -c gokanon -n "__fish_seen_subcommand_from compare" -o format -d "Output format" -a "table json"

# export command options
complete -c gokanon -n "__fish_seen_subcommand_from export" -l latest -d "Export latest comparison"
complete -c gokanon -n "__fish_seen_subcommand_from export" -o format -d "Export format" -a "html csv markdown json"
complete -c gokanon -n "__fish_seen_subcommand_from export" -o output -d "Output file" -r
complete -c gokanon -n "__fish_seen_subcommand_from export" -o storage -d "Storage directory" -r

# stats and trend command options
complete -c gokanon -n "__fish_seen_subcommand_from stats trend" -o last -d "Number of runs"
complete -c gokanon -n "__fish_seen_subcommand_from stats trend" -o storage -d "Storage directory" -r
complete -c gokanon -n "__fish_seen_subcommand_from stats trend" -o format -d "Output format" -a "table json"

# check command options
complete -c gokanon -n "__fish_seen_subcommand_from check" -l latest -d "Check latest two runs"
complete -c gokanon -n "__fish_seen_subcommand_from check" -o threshold -d "Threshold percentage"
complete -c gokanon -n "__fish_seen_subcommand_from check" -o storage -d "Storage directory" -r
complete -c gokanon -n "__fish_seen_subcommand_from check" -o format -d "Output format" -a "table json"

# serve and flamegraph command options
complete -c gokanon -n "__fish_seen_subcommand_from serve flamegraph" -o port -d "Server port"
complete -c gokanon -n "__fish_seen_subcommand_from serve flamegraph" -o storage -d "Storage directory" -r
complete -c gokanon -n "__fish_seen_subcommand_from serve flamegraph" -o open -d "Open browser automatically"

# baseline command - subcommands
complete -c gokanon -f -n "__fish_seen_subcommand_from baseline; and not __fish_seen_subcommand_from save list show delete" -a save -d "Save a benchmark run as baseline"
complete -c gokanon -f -n "__fish_seen_subcommand_from baseline; and not __fish_seen_subcommand_from save list show delete" -a list -d "List all saved baselines"
complete -c gokanon -f -n "__fish_seen_subcommand_from baseline; and not __fish_seen_subcommand_from save list show delete" -a show -d "Show details of a baseline"
complete -c gokanon -f -n "__fish_seen_subcommand_from baseline; and not __fish_seen_subcommand_from save list show delete" -a delete -d "Delete a baseline"

# baseline save options
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from save" -o name -d "Baseline name" -r
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from save" -o run -d "Run ID to save" -r
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from save" -o desc -d "Baseline description" -r
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from save" -o storage -d "Storage directory" -r

# baseline list options
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from list" -o storage -d "Storage directory" -r

# baseline show options
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from show" -o name -d "Baseline name" -r
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from show" -o storage -d "Storage directory" -r

# baseline delete options
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from delete" -o name -d "Baseline name" -r
complete -c gokanon -n "__fish_seen_subcommand_from baseline; and __fish_seen_subcommand_from delete" -o storage -d "Storage directory" -r

# completion command options
complete -c gokanon -f -n "__fish_seen_subcommand_from completion" -a "bash zsh fish" -d "Shell type"

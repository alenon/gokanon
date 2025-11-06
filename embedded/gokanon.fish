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
complete -c gokanon -f -n __fish_use_subcommand -a doctor -d "Run diagnostics"
complete -c gokanon -f -n __fish_use_subcommand -a interactive -d "Start interactive mode"
complete -c gokanon -f -n __fish_use_subcommand -a help -d "Show help message"

# run command options
complete -c gokanon -n "__fish_seen_subcommand_from run" -o bench -d "Benchmark pattern"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o pkg -d "Package pattern"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o profile -d "Enable profiling" -a "cpu mem cpu,mem"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o benchtime -d "Benchmark duration"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o count -d "Run count"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o timeout -d "Test timeout"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o cpu -d "CPU counts"
complete -c gokanon -n "__fish_seen_subcommand_from run" -o v -d "Verbose output"

# compare command options
complete -c gokanon -n "__fish_seen_subcommand_from compare" -l latest -d "Compare latest two runs"
complete -c gokanon -n "__fish_seen_subcommand_from compare" -o format -d "Output format" -a "table json"

# export command options
complete -c gokanon -n "__fish_seen_subcommand_from export" -l latest -d "Export latest comparison"
complete -c gokanon -n "__fish_seen_subcommand_from export" -o format -d "Export format" -a "html csv markdown json"
complete -c gokanon -n "__fish_seen_subcommand_from export" -o output -d "Output file" -r

# stats and trend command options
complete -c gokanon -n "__fish_seen_subcommand_from stats trend" -o last -d "Number of runs"
complete -c gokanon -n "__fish_seen_subcommand_from stats trend" -o format -d "Output format" -a "table json"

# check command options
complete -c gokanon -n "__fish_seen_subcommand_from check" -l latest -d "Check latest two runs"
complete -c gokanon -n "__fish_seen_subcommand_from check" -o threshold -d "Threshold percentage"
complete -c gokanon -n "__fish_seen_subcommand_from check" -o format -d "Output format" -a "table json"

# serve and flamegraph command options
complete -c gokanon -n "__fish_seen_subcommand_from serve flamegraph" -o port -d "Server port"
complete -c gokanon -n "__fish_seen_subcommand_from serve flamegraph" -o open -d "Open browser automatically"

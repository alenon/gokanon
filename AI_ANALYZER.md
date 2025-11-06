# AI Analyzer for GoKanon

GoKanon now includes an AI-powered analyzer that provides intelligent insights and optimization suggestions for your benchmark results. The AI analyzer uses free AI services to analyze profiling data and comparison results.

## Features

- **Enhanced Profile Analysis**: AI-powered suggestions for CPU and memory optimization
- **Comparison Insights**: Intelligent analysis of performance changes between runs
- **Multiple AI Providers**: Supports Ollama, OpenAI, Anthropic Claude, Google Gemini, Groq, and any OpenAI-compatible service
- **Zero Configuration**: Works out of the box with sensible defaults

## Supported AI Providers

### 1. Ollama (Recommended for Local/Private Use)

**Pros:**
- Completely free
- No API keys required
- Private - your data stays local
- No rate limits
- Supports many models (Llama 3.2, Mistral, Gemma, etc.)

**Setup:**
1. Install Ollama: https://ollama.ai/
2. Pull a model: `ollama pull llama3.2`
3. Ensure Ollama is running: `ollama serve`

**Configuration:**
```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=ollama
export GOKANON_AI_MODEL=llama3.2  # Optional, auto-detected
```

### 2. OpenAI (GPT-4o, GPT-4-turbo)

**Pros:**
- State-of-the-art performance
- Fast inference
- Well-documented API
- Supports latest GPT-4o and GPT-4-turbo models

**Cons:**
- Requires paid API key (usage-based pricing)
- Data sent to OpenAI servers

**Setup:**
1. Sign up at https://platform.openai.com
2. Get your API key from https://platform.openai.com/api-keys
3. Set environment variables

**Configuration:**
```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai
export GOKANON_AI_API_KEY=sk-your-openai-key
export GOKANON_AI_MODEL=gpt-4o  # Optional: gpt-4o, gpt-4-turbo, gpt-4, gpt-3.5-turbo
```

### 3. Anthropic Claude (Sonnet 4.5, Haiku 4.5)

**Pros:**
- Excellent reasoning and analysis capabilities
- Claude Sonnet 4.5 is one of the best models available
- Fast inference with Haiku 4.5
- Strong at understanding complex code patterns

**Cons:**
- Requires paid API key (usage-based pricing)
- Data sent to Anthropic servers

**Setup:**
1. Sign up at https://console.anthropic.com
2. Get your API key
3. Set environment variables

**Configuration:**
```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=anthropic  # or "claude"
export GOKANON_AI_API_KEY=sk-ant-your-anthropic-key
export GOKANON_AI_MODEL=claude-sonnet-4-5-20250929  # Optional: claude-haiku-4-5 for faster/cheaper
```

**Pricing (as of 2025):**
- Claude Sonnet 4.5: $3/$15 per million input/output tokens
- Claude Haiku 4.5: $1/$5 per million input/output tokens

### 4. Groq (Fast Cloud Inference)

**Pros:**
- Very fast inference (fastest in the market)
- Free tier available
- No local installation needed
- Supports Llama 3.3 70B and other models

**Cons:**
- Rate limits on free tier
- Data sent to Groq servers

**Setup:**
1. Sign up at https://console.groq.com
2. Get your API key
3. Set environment variable

**Configuration:**
```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=groq
export GOKANON_AI_API_KEY=your-groq-key
export GOKANON_AI_MODEL=llama-3.3-70b-versatile  # Optional
```

### 5. Google Gemini (Gemini 2.5 Flash, 2.0 Flash)

**Pros:**
- Very affordable pricing (33% cheaper than previous versions)
- Latest Gemini 2.5 Flash with excellent performance
- Multimodal support (text, images, audio, video)
- Large context windows (up to 1M tokens)
- Free tier available through Google AI Studio

**Cons:**
- Requires Google API key
- Data sent to Google servers

**Setup:**
1. Get API key from Google AI Studio: https://aistudio.google.com/apikey
2. Or use Google Cloud Console
3. Set environment variables

**Configuration:**
```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=gemini
export GOKANON_AI_API_KEY=your-google-api-key
export GOKANON_AI_MODEL=gemini-2.5-flash  # Optional: gemini-2.0-flash, gemini-2.5-pro
```

**Pricing (as of 2025):**
- Gemini 2.0 Flash: $0.10/$0.40 per million input/output tokens
- Gemini 2.5 Flash: $0.30/$2.50 per million input/output tokens
- Gemini 2.5 Pro: Higher pricing for more advanced capabilities

**Available Models:**
- `gemini-2.5-flash` - Latest, best balance of speed and quality (default)
- `gemini-2.5-pro` - Most powerful, adaptive thinking
- `gemini-2.0-flash` - Faster, more affordable option

### 6. OpenAI-Compatible (Cursor, LM Studio, LocalAI, etc.)

**Pros:**
- Works with any OpenAI-compatible API
- Supports Cursor AI, LM Studio, LocalAI, vLLM, and more
- Flexible for custom setups
- API key optional for local services

**Use Cases:**
- Cursor AI with custom models
- LM Studio running local models
- LocalAI for on-premise deployments
- vLLM or text-generation-inference servers

**Setup:**
1. Set up your OpenAI-compatible service (e.g., LM Studio, Cursor)
2. Get the API endpoint URL
3. Configure GoKanon

**Configuration:**
```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai-compatible  # or "custom"
export GOKANON_AI_BASE_URL=http://localhost:1234  # Your service URL
export GOKANON_AI_MODEL=your-model-name
export GOKANON_AI_API_KEY=optional-if-needed
```

**Example for LM Studio:**
```bash
# Start LM Studio and load a model, then:
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai-compatible
export GOKANON_AI_BASE_URL=http://localhost:1234/v1
export GOKANON_AI_MODEL=local-model
```

## Usage

### Basic Usage

Enable AI analysis by setting the environment variable:

```bash
export GOKANON_AI_ENABLED=true
```

Then run benchmarks with profiling:

```bash
gokanon run -profile=cpu,memory
```

The AI analyzer will automatically enhance the profile summary with intelligent suggestions.

### Comparison Analysis

When comparing benchmark runs:

```bash
gokanon compare --latest
```

If AI is enabled, you'll see an "AI Analysis" section explaining the performance changes.

## Configuration

Configure the AI analyzer using environment variables:

### Required Settings

```bash
# Enable AI analysis (default: false)
export GOKANON_AI_ENABLED=true
```

### Optional Settings

```bash
# AI provider (default: ollama)
# Options: "ollama", "openai", "anthropic" (or "claude"), "gemini", "groq", "openai-compatible" (or "custom")
export GOKANON_AI_PROVIDER=ollama

# Model name (auto-detected if not set)
# Examples:
#   ollama: llama3.2, mistral, gemma
#   openai: gpt-4o, gpt-4-turbo, gpt-3.5-turbo
#   anthropic: claude-sonnet-4-5-20250929, claude-haiku-4-5
#   gemini: gemini-2.5-flash, gemini-2.0-flash, gemini-2.5-pro
#   groq: llama-3.3-70b-versatile, mixtral-8x7b-32768
export GOKANON_AI_MODEL=llama3.2

# Base URL for the AI service (auto-detected if not set)
# Defaults:
#   ollama: http://localhost:11434
#   openai: https://api.openai.com
#   anthropic: https://api.anthropic.com
#   gemini: https://generativelanguage.googleapis.com
#   groq: https://api.groq.com/openai/v1
#   openai-compatible: http://localhost:8080 (must be set by user)
export GOKANON_AI_BASE_URL=http://localhost:11434

# API key (required for cloud providers, optional for local)
# Required for: openai, anthropic, gemini, groq
# Not required for: ollama, some openai-compatible services
export GOKANON_AI_API_KEY=your-api-key
```

## Example Workflows

### Local Development with Ollama

```bash
# 1. Install and start Ollama
ollama serve &
ollama pull llama3.2

# 2. Enable AI analysis
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=ollama

# 3. Run benchmarks with profiling
gokanon run -profile=cpu,memory

# 4. View AI-enhanced suggestions in the output
```

### CI/CD with Groq

```bash
# 1. Set up Groq API key in your CI environment
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=groq
export GOKANON_AI_API_KEY=$GROQ_API_KEY

# 2. Run benchmarks
gokanon run -profile=cpu,memory

# 3. Compare with baseline
gokanon compare --latest
```

### Production Analysis with OpenAI

```bash
# 1. Set up OpenAI API key
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai
export GOKANON_AI_API_KEY=sk-your-openai-key
export GOKANON_AI_MODEL=gpt-4o

# 2. Run benchmarks with profiling
gokanon run -profile=cpu,memory

# 3. Get high-quality AI analysis
gokanon compare --latest
```

### Advanced Analysis with Claude

```bash
# 1. Set up Anthropic API key
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=anthropic
export GOKANON_AI_API_KEY=sk-ant-your-anthropic-key
export GOKANON_AI_MODEL=claude-sonnet-4-5-20250929

# 2. Run benchmarks
gokanon run -profile=cpu,memory

# 3. Claude provides detailed reasoning about performance patterns
```

### Affordable Analysis with Gemini

```bash
# 1. Set up Google API key
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=gemini
export GOKANON_AI_API_KEY=your-google-api-key
export GOKANON_AI_MODEL=gemini-2.5-flash

# 2. Run benchmarks with profiling
gokanon run -profile=cpu,memory

# 3. Get cost-effective AI analysis with large context windows
gokanon compare --latest
```

### Local LM Studio Setup

```bash
# 1. Install and start LM Studio (https://lmstudio.ai)
# 2. Load a model (e.g., Llama 3 8B)
# 3. Start the local server (default port 1234)

# 4. Configure GoKanon
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai-compatible
export GOKANON_AI_BASE_URL=http://localhost:1234/v1
export GOKANON_AI_MODEL=local-model

# 5. Run benchmarks
gokanon run -profile=cpu,memory
```

### Comparing Performance Changes

```bash
# Run baseline
gokanon run -profile=cpu,memory

# Make code changes
# ...

# Run new benchmark
gokanon run -profile=cpu,memory

# Compare with AI insights (works with any provider)
export GOKANON_AI_ENABLED=true
gokanon compare --latest
```

## What the AI Analyzer Does

### Profile Analysis

When you run benchmarks with profiling enabled, the AI analyzer:

1. **Analyzes CPU hotspots**: Identifies functions consuming significant CPU time
2. **Detects memory issues**: Finds excessive allocations and potential leaks
3. **Evaluates hot paths**: Examines critical execution paths
4. **Generates actionable suggestions**: Provides specific optimization recommendations with:
   - Function names affected
   - Issue descriptions
   - Concrete suggestions
   - Expected impact

### Comparison Analysis

When comparing benchmark runs, the AI analyzer:

1. **Identifies significant changes**: Highlights important improvements or regressions
2. **Explains possible causes**: Suggests reasons for performance changes
3. **Provides context**: Evaluates whether changes are concerning
4. **Recommends next steps**: Suggests follow-up actions

## Example Output

### Profile Analysis

```
--- Profile Summary ---

Top CPU Functions:
  1. mypackage.ProcessData (45.2%)
  2. runtime.mallocgc (12.8%)
  3. mypackage.ComputeHash (8.3%)

Suggestions:
  [HIGH] CPU - mypackage.ProcessData
    Issue: Function consumes 45.2% of CPU time
    Suggestion: This function performs multiple string concatenations in a loop.
                Consider using strings.Builder or bytes.Buffer for better performance.
    Impact: Could improve overall performance by 30-40%

  [MEDIUM] Memory - mypackage.ComputeHash
    Issue: High allocation rate detected
    Suggestion: Function allocates a new slice on each call. Consider using sync.Pool
                to reuse buffers across calls.
    Impact: Could reduce allocation pressure and GC overhead significantly
```

### Comparison Analysis

```
--- AI Analysis ---

Performance Improvements:
The benchmarks show a 23% improvement in ProcessData, which is excellent. This
improvement is likely due to the recent optimization of string handling using
strings.Builder instead of concatenation.

Concerns:
ComputeHash shows a 5% regression. While minor, this function is called frequently
and the cumulative impact could be significant. The regression might be caused by
the additional boundary checks in the new implementation.

Recommendations:
1. The ProcessData optimization is solid and should be kept
2. Review the ComputeHash implementation - consider benchmarking it in isolation
3. Monitor the overall trend over the next few runs to confirm stability
```

## Troubleshooting

### "AI analysis failed"

**Ollama:**
- Check if Ollama is running: `curl http://localhost:11434/api/version`
- Ensure the model is pulled: `ollama list`
- Try pulling the model: `ollama pull llama3.2`

**Groq:**
- Verify your API key is set correctly
- Check your Groq quota at https://console.groq.com
- Ensure you have network connectivity

### "Failed to initialize AI analyzer"

- Check environment variables are set correctly
- Verify the provider name is either "ollama" or "groq"
- For Groq, ensure GOKANON_AI_API_KEY is set

### Slow Analysis

- **Ollama**: First run may be slow as the model loads; subsequent runs are faster
- **Groq**: Usually fast, but may be rate-limited on free tier
- Consider using smaller models for faster analysis (e.g., `llama3.2` instead of larger models)

## Privacy and Data

- **Ollama**: All data stays on your machine. No data is sent to external services.
- **Groq**: Benchmark data is sent to Groq's API. Review Groq's privacy policy at https://groq.com/privacy-policy

## Performance Impact

The AI analyzer has minimal performance impact:

- **Profile enhancement**: Adds 1-3 seconds to benchmark runs with profiling
- **Comparison analysis**: Adds 1-2 seconds to comparison commands
- **No impact on benchmark execution**: AI analysis happens after benchmarks complete

## Disabling AI Analysis

To disable AI analysis:

```bash
unset GOKANON_AI_ENABLED
# or
export GOKANON_AI_ENABLED=false
```

GoKanon will work normally without AI analysis, falling back to built-in suggestions.

## Future Enhancements

Planned features:
- Support for more AI providers (Anthropic Claude, OpenAI)
- Custom prompts and analysis templates
- Historical trend analysis with AI insights
- Integration with web dashboard
- Caching of AI responses for repeated analyses

## Contributing

We welcome contributions to improve the AI analyzer! Areas of interest:
- Additional AI provider integrations
- Improved prompt engineering for better suggestions
- Better parsing of AI responses
- Performance optimizations

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

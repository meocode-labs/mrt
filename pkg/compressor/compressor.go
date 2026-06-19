package compressor

import (
	"bufio"
	"bytes"
	"math"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

type Compressor struct {
	mu           sync.RWMutex
	tokenCount   int64
	lineCount    int64
	duplicateCnt int64
	ansiRegex    *regexp.Regexp
	config       CompressionConfig
}

type CompressionConfig struct {
	StripANSI      bool
	CollapsePaths  bool
	RemoveDupLines bool
	Skeletonize    bool
	TargetProfile  string
}

type CompressionResult struct {
	OriginalLen   int
	CompressedLen int
	TokensRemoved int64
	LinesRemoved  int64
	Duplicates    int64
	ReductionRate float64
}

var (
	ansiEscapeRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	pathRegex       = regexp.MustCompile(`/home/[^/\s]+`)
	progressRegex   = regexp.MustCompile(`(\d+)%`)
	duplicateTracker sync.Map
)

const (
	TOKEN_PRICE_USD = 0.00001
)

func NewCompressor(config CompressionConfig) *Compressor {
	if config.StripANSI {
		config.StripANSI = true
	}
	return &Compressor{
		config: config,
	}
}

func (c *Compressor) Compress(input string) (string, CompressionResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	originalLen := len(input)
	result := CompressionResult{OriginalLen: originalLen}

	lines := c.splitLines(input)
	compressedLines := c.processLines(lines)

	result.Duplicates = c.duplicateCnt
	result.LinesRemoved = c.lineCount
	result.TokensRemoved = c.estimateTokensRemoved(input, compressedLines)

	output := c.joinLines(compressedLines)
	result.CompressedLen = len(output)

	if originalLen > 0 {
		result.ReductionRate = float64(originalLen-result.CompressedLen) / float64(originalLen) * 100
	}

	c.tokenCount += result.TokensRemoved

	return output, result
}

func (c *Compressor) splitLines(input string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func (c *Compressor) processLines(lines []string) []string {
	var result []string
	seen := make(map[string]int)

	for _, line := range lines {
		processed := c.processLine(line)

		if processed == "" {
			continue
		}

		if c.config.RemoveDupLines {
			normalized := c.normalizeForDedup(processed)
			if count, exists := seen[normalized]; exists {
				seen[normalized]++
				c.duplicateCnt++
				continue
			}
			seen[normalized] = 1
		}

		result = append(result, processed)
		c.lineCount++
	}

	return result
}

func (c *Compressor) processLine(line string) string {
	if c.config.StripANSI {
		line = c.stripANSI(line)
	}

	if c.config.CollapsePaths {
		line = c.collapsePaths(line)
	}

	line = c.collapseProgressBars(line)

	return strings.TrimRight(line, " \t")
}

func (c *Compressor) stripANSI(input string) string {
	return ansiEscapeRegex.ReplaceAllString(input, "")
}

func (c *Compressor) collapsePaths(input string) string {
	return pathRegex.ReplaceAllStringFunc(input, func(match string) string {
		parts := strings.Split(match, "/")
		if len(parts) >= 3 {
			return "~/" + parts[len(parts)-2] + "/" + parts[len(parts)-1]
		}
		return match
	})
}

func (c *Compressor) collapseProgressBars(input string) string {
	matches := progressRegex.FindAllStringIndex(input, -1)
	if len(matches) <= 1 {
		return input
	}

	var result bytes.Buffer
	lastEnd := 0

	seenPercents := make(map[string]bool)
	for _, match := range matches {
		result.WriteString(input[lastEnd:match[0]])

		percent := input[match[0]:match[1]]
		if seenPercents[percent] {
			result.WriteString("...")
			c.duplicateCnt++
		} else {
			result.WriteString(percent)
			seenPercents[percent] = true
		}
		lastEnd = match[1]
	}
	result.WriteString(input[lastEnd:])

	return result.String()
}

func (c *Compressor) normalizeForDedup(line string) string {
	line = ansiEscapeRegex.ReplaceAllString(line, "")
	line = progressRegex.ReplaceAllString(line, "#")
	line = strings.TrimSpace(line)
	return line
}

func (c *Compressor) estimateTokensRemoved(original string, compressed []string) int64 {
	originalTokens := c.countTokens(original)
	compressedTokens := int64(0)

	for _, line := range compressed {
		compressedTokens += c.countTokens(line)
	}

	removed := originalTokens - compressedTokens
	if removed < 0 {
		return 0
	}
	return removed
}

func (c *Compressor) countTokens(text string) int64 {
	var count int64
	var token strings.Builder
	inToken := false

	for _, r := range text {
		if unicode.IsSpace(r) {
			if inToken {
				count++
				inToken = false
				token.Reset()
			}
			continue
		}

		token.WriteRune(r)
		inToken = true

		if token.Len() >= 4 {
			count++
			inToken = false
			token.Reset()
		}
	}

	if inToken {
		count++
	}

	return int64(math.Max(1, float64(count)/4))
}

func (c *Compressor) SkeletonizeCode(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	inBlock := false

	sigRegex := regexp.MustCompile(`^(func|def|class|struct|interface|enum|pub|priv|export)\s+`)
	blockStart := regexp.MustCompile(`\{|\:$`)
	blockEnd := regexp.MustCompile(`^\s*\}|\s*pass\s*$`)

	for _, line := range lines {
		if blockStart.MatchString(line) && !strings.Contains(line, "{}") && !strings.HasSuffix(line, ": {}") {
			result = append(result, line)
			result = append(result, "    ...")
			inBlock = true
			continue
		}

		if blockEnd.MatchString(line) && inBlock {
			inBlock = false
			continue
		}

		if inBlock {
			continue
		}

		if sigRegex.MatchString(line) {
			result = append(result, line)
			result = append(result, "    ...")
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func (c *Compressor) GetStats() (tokens, lines, duplicates int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tokenCount, c.lineCount, c.duplicateCnt
}

func (c *Compressor) CalculateCostSaved() float64 {
	tokens, _, _ := c.GetStats()
	return float64(tokens) * TOKEN_PRICE_USD
}

func StripANSIColors(input string) string {
	return ansiEscapeRegex.ReplaceAllString(input, "")
}

func CountTokens(input string) int64 {
	var count int64
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		count++
	}

	return int64(math.Max(1, float64(count)/4))
}

func RemoveDuplicateLines(input string) string {
	var lines []string
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || seen[line] {
			continue
		}
		seen[line] = true
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func CollapseVerbosePaths(input string) string {
	return pathRegex.ReplaceAllStringFunc(input, func(match string) string {
		parts := strings.Split(match, "/")
		if len(parts) >= 3 {
			return "~/" + parts[len(parts)-2] + "/" + parts[len(parts)-1]
		}
		return match
	})
}

func NormalizeProgressBars(input string) string {
	return progressRegex.ReplaceAllString(input, "[#]")
}

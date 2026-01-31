package expression

import (
	"fmt"
	"regexp"
	"strings"
)

// Parser parses expression strings.
type Parser struct {
	expressionPattern *regexp.Regexp
}

// NewParser creates a new expression parser.
func NewParser() *Parser {
	return &Parser{
		expressionPattern: regexp.MustCompile(`\$\{\{\s*(.+?)\s*\}\}`),
	}
}

// Parse parses a string that may contain expressions.
func (p *Parser) Parse(input string) (*Expression, error) {
	expr := &Expression{Raw: input}

	matches := p.expressionPattern.FindAllStringSubmatchIndex(input, -1)
	if len(matches) == 0 {
		// No expressions, just a literal
		expr.Segments = []Segment{LiteralSegment{Value: input}}
		return expr, nil
	}

	lastEnd := 0
	for _, match := range matches {
		// Add literal segment before this expression
		if match[0] > lastEnd {
			expr.Segments = append(expr.Segments, LiteralSegment{
				Value: input[lastEnd:match[0]],
			})
		}

		// Parse the expression content
		exprContent := input[match[2]:match[3]]
		ref, err := p.parseReference(exprContent)
		if err != nil {
			return nil, fmt.Errorf("invalid expression %q: %w", exprContent, err)
		}
		expr.Segments = append(expr.Segments, ref)

		lastEnd = match[1]
	}

	// Add trailing literal if any
	if lastEnd < len(input) {
		expr.Segments = append(expr.Segments, LiteralSegment{
			Value: input[lastEnd:],
		})
	}

	return expr, nil
}

// parseReference parses a reference like "databases.main.url" or "dependents.*.routes.*.url | join ','"
func (p *Parser) parseReference(content string) (ReferenceSegment, error) {
	parts := strings.Split(content, "|")
	pathStr := strings.TrimSpace(parts[0])

	ref := ReferenceSegment{
		Path: p.parsePath(pathStr),
	}

	// Parse pipe functions
	for i := 1; i < len(parts); i++ {
		pipeStr := strings.TrimSpace(parts[i])
		pf, err := p.parsePipeFunc(pipeStr)
		if err != nil {
			return ref, err
		}
		ref.Pipe = append(ref.Pipe, pf)
	}

	return ref, nil
}

// parsePath parses a dotted path like "databases.main.url"
func (p *Parser) parsePath(pathStr string) []string {
	// Handle array access notation
	pathStr = strings.ReplaceAll(pathStr, "[", ".")
	pathStr = strings.ReplaceAll(pathStr, "]", "")

	return strings.Split(pathStr, ".")
}

// parsePipeFunc parses a pipe function like "join ','" or "first"
func (p *Parser) parsePipeFunc(pipeStr string) (PipeFunc, error) {
	// Split by space, first part is function name
	parts := strings.Fields(pipeStr)
	if len(parts) == 0 {
		return PipeFunc{}, fmt.Errorf("empty pipe function")
	}

	pf := PipeFunc{Name: parts[0]}

	// Parse arguments (handle quoted strings)
	if len(parts) > 1 {
		for _, arg := range parts[1:] {
			// Remove quotes if present
			arg = strings.Trim(arg, `"'`)
			pf.Args = append(pf.Args, arg)
		}
	}

	return pf, nil
}

// ContainsExpression checks if a string contains ${{ }} expressions.
func ContainsExpression(s string) bool {
	return strings.Contains(s, "${{") && strings.Contains(s, "}}")
}

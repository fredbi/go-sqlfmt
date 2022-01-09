package group

type (
	CommaStyle uint8

	Option func(*options)

	options struct {
		IndentLevel          int
		commaStyle           CommaStyle
		hasCommaBefore       bool
		hasParenthesisBefore bool
		hasCastBefore        bool
		indentSize           int
	}
)

const (
	CommaStyleLeft CommaStyle = iota
	CommaStyleRight
)

func defaultOptions(opts ...Option) *options {
	o := &options{
		commaStyle: CommaStyleLeft,
		indentSize: 2,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

func WithIndentLevel(level int) Option {
	return func(opts *options) {
		opts.IndentLevel = level
	}
}

func WithCommaStyle(style CommaStyle) Option {
	return func(opts *options) {
		opts.commaStyle = style
	}
}

// WithHasCommaBefore instructs the group about the comma-specific indentation context.
func WithHasCommaBefore(enabled bool) Option {
	return func(opts *options) {
		opts.hasCommaBefore = enabled
	}
}

// WithHasParenthesisBefore instructs the group about the parenthesis-specific indentation context.
func WithHasParenthesisBefore(enabled bool) Option {
	return func(opts *options) {
		opts.hasParenthesisBefore = enabled
	}
}

// WithHasCastBefore instructs the group about the oerator-specific indentation context.
func WithHasCastBefore(enabled bool) Option {
	return func(opts *options) {
		opts.hasCastBefore = enabled
	}
}

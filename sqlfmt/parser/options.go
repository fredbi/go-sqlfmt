package parser

import "github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"

type (
	// Option for the parser.
	Option func(*options)

	options struct {
		groupOptions []group.Option
		afterComma   bool // TODO: generalize to any contextual information to pass to ReIndenters
	}
)

func (o *options) ToGroupOptions() []group.Option {
	res := make([]group.Option, 0, len(o.groupOptions)+1)
	res = append(res, o.groupOptions...)
	res = append(res, group.WithHasCommaBefore(o.afterComma))

	return res
}

func (o *options) CloneWithOptions(opts ...Option) *options {
	c := *o

	for _, apply := range opts {
		apply(&c)
	}

	return &c
}

func defaultOptions(opts ...Option) *options {
	o := &options{}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

// WithGroupOptions specifies some grouping options.
func WithGroupOptions(groupOptions ...group.Option) Option {
	return func(opts *options) {
		opts.groupOptions = groupOptions
	}
}

func withOptions(o *options) Option {
	return func(opts *options) {
		if o == nil {
			o = defaultOptions()
		}
		*opts = *o
	}
}

// withAfterComma produces some formatting context about the position of
// a group following a comma or not.
func withAfterComma(afterComma bool) Option {
	return func(opts *options) {
		opts.afterComma = afterComma
	}
}

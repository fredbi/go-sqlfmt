package parser

// TODO: make an augmented token that knows about grouping?

import (
	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"
)

var _ group.Reindenter = lexer.Token{}

func (r *Retriever) getGroupFromTokens(firstToken lexer.Token, tokenSource []group.Reindenter) group.Reindenter {
	if groupBuilder := getGroupBuilder(firstToken); groupBuilder != nil {
		return groupBuilder(tokenSource, r.ToGroupOptions()...)
	}

	return nil
}

type groupBuilder func([]group.Reindenter, ...group.Option) group.Reindenter

func getGroupBuilder(token lexer.Token) groupBuilder {
	switch token.Type {
	case lexer.SELECT:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewSelect(tokenSource, opts...)
		}
	case lexer.FROM:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewFrom(tokenSource, opts...)
		}
	case lexer.JOIN, lexer.INNER, lexer.OUTER, lexer.LEFT, lexer.RIGHT, lexer.NATURAL, lexer.CROSS, lexer.LATERAL:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewJoin(tokenSource, opts...)
		}
	case lexer.WHERE:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewWhere(tokenSource, opts...)
		}
	case lexer.ANDGROUP:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewAndGroup(tokenSource, opts...)
		}
	case lexer.ORGROUP:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewOrGroup(tokenSource, opts...)
		}
	case lexer.GROUP:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewGroupBy(tokenSource, opts...)
		}
	case lexer.ORDER:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewOrderBy(tokenSource, opts...)
		}
	case lexer.HAVING:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewHaving(tokenSource, opts...)
		}
	case lexer.LIMIT, lexer.OFFSET, lexer.FETCH:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewLimitClause(tokenSource, opts...)
		}
	case lexer.UNION, lexer.INTERSECT, lexer.EXCEPT:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewTieClause(tokenSource, opts...)
		}
	case lexer.UPDATE:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewUpdate(tokenSource, opts...)
		}
	case lexer.SET:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewSet(tokenSource, opts...)
		}
	case lexer.RETURNING:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewReturning(tokenSource, opts...)
		}
	case lexer.LOCK:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewLock(tokenSource, opts...)
		}
	case lexer.INSERT:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewInsert(tokenSource, opts...)
		}
	case lexer.VALUES:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewValues(tokenSource, opts...)
		}
	case lexer.DELETE:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewDelete(tokenSource, opts...)
		}
	case lexer.WITH:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			return group.NewWith(tokenSource, opts...)
		}
	case lexer.CASE:
		// TODO: find a more elegant way to include the end token in the group
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			// endKeyWord of CASE group("END") has to be included in the group, so it is appended to result
			endToken := lexer.MakeToken(lexer.END, "END", lexer.WithOptionsFrom(token))
			tokenSource = append(tokenSource, endToken)

			return group.NewCase(tokenSource, opts...)
		}
	case lexer.STARTPARENTHESIS:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			// endKeyWord of subQuery group (")") has to be included in the group, so it is appended to result
			endToken := lexer.MakeToken(lexer.ENDPARENTHESIS, ")", lexer.WithOptionsFrom(token))
			tokenSource = append(tokenSource, endToken)

			if _, isSubQuery := tokenSource[1].(*group.Select); isSubQuery {
				return group.NewSubquery(tokenSource, opts...)
			}

			return group.NewParenthesis(tokenSource, opts...)
		}
	case lexer.FUNCTION:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			endToken := lexer.MakeToken(lexer.ENDPARENTHESIS, ")", lexer.WithOptionsFrom(token))
			tokenSource = append(tokenSource, endToken)

			return group.NewFunction(tokenSource, opts...)
		}
	case lexer.TYPE:
		return func(tokenSource []group.Reindenter, opts ...group.Option) group.Reindenter {
			endToken := lexer.MakeToken(lexer.ENDPARENTHESIS, ")", lexer.WithOptionsFrom(token))
			tokenSource = append(tokenSource, endToken)

			return group.NewTypeCast(tokenSource, opts...)
		}
	default:
		return nil
	}
}

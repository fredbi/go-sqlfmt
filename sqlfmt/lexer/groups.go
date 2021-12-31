package lexer

// EndTokenTypes returns the corresponding end of clause tokens, or nil if this token does not define a group.
func (t Token) EndTokenTypes() map[TokenType]struct{} {
	switch t.Type {
	case SELECT:
		return map[TokenType]struct{}{
			FROM:  {},
			UNION: {},
			// TODO: need INTERSECT???
		}
	case FROM:
		return map[TokenType]struct{}{
			WHERE: {},
			INNER: {}, OUTER: {}, LEFT: {}, RIGHT: {}, JOIN: {},
			NATURAL: {}, CROSS: {}, LATERAL: {},
			UNION: {}, INTERSECT: {},
			ORDER: {}, GROUP: {},
			LIMIT: {}, OFFSET: {}, FETCH: {}, EXCEPT: {},
			ENDPARENTHESIS: {},
		}
	case CASE:
		return map[TokenType]struct{}{END: {}}
	case JOIN, INNER, OUTER, LEFT, RIGHT, NATURAL, CROSS:
		return map[TokenType]struct{}{
			WHERE: {},
			ORDER: {}, GROUP: {},
			LIMIT: {}, OFFSET: {}, FETCH: {}, EXCEPT: {},
			ANDGROUP: {}, ORGROUP: {},
			LEFT: {}, RIGHT: {}, INNER: {}, OUTER: {},
			NATURAL: {}, CROSS: {}, LATERAL: {},
			UNION: {}, INTERSECT: {},
			ENDPARENTHESIS: {},
		}
	case WHERE:
		return map[TokenType]struct{}{
			GROUP: {}, ORDER: {},
			LIMIT: {}, OFFSET: {}, FETCH: {}, EXCEPT: {},
			ANDGROUP: {}, OR: {},
			UNION:          {},
			INTERSECT:      {},
			RETURNING:      {},
			ENDPARENTHESIS: {},
		}
	case ANDGROUP:
		return map[TokenType]struct{}{
			GROUP: {}, ORDER: {},
			LIMIT: {}, OFFSET: {}, FETCH: {}, EXCEPT: {},
			UNION: {}, INTERSECT: {},
			ANDGROUP: {}, ORGROUP: {},
			ENDPARENTHESIS: {},
		}
	case ORGROUP:
		return map[TokenType]struct{}{
			GROUP: {}, ORDER: {},
			LIMIT: {}, OFFSET: {}, FETCH: {}, EXCEPT: {},
			UNION: {}, INTERSECT: {},
			ANDGROUP: {}, ORGROUP: {},
			ENDPARENTHESIS: {},
		}
	case GROUP:
		return map[TokenType]struct{}{
			ORDER: {},
			LIMIT: {}, FETCH: {}, OFFSET: {}, EXCEPT: {},
			UNION: {}, INTERSECT: {},
			HAVING:         {},
			ENDPARENTHESIS: {},
		}
	case HAVING:
		return map[TokenType]struct{}{
			ORDER: {},
			LIMIT: {}, OFFSET: {}, FETCH: {}, EXCEPT: {},
			UNION: {}, INTERSECT: {},
			ENDPARENTHESIS: {},
		}
	case ORDER:
		return map[TokenType]struct{}{
			GROUP: {},
			LIMIT: {}, FETCH: {}, OFFSET: {}, EXCEPT: {},
			UNION: {}, INTERSECT: {},
			ENDPARENTHESIS: {},
		}
	case LIMIT, FETCH, OFFSET:
		return map[TokenType]struct{}{
			UNION: {}, INTERSECT: {},
			EXCEPT:         {},
			ENDPARENTHESIS: {},
		}
	case STARTPARENTHESIS:
		return map[TokenType]struct{}{ENDPARENTHESIS: {}}
	case UNION, INTERSECT, EXCEPT:
		return map[TokenType]struct{}{SELECT: {}}
	case UPDATE:
		return map[TokenType]struct{}{
			WHERE:     {},
			SET:       {},
			RETURNING: {},
		}
	case SET:
		return map[TokenType]struct{}{
			WHERE:     {},
			RETURNING: {},
		}
	case RETURNING:
		return map[TokenType]struct{}{
			EOF: {},
		}
	case DELETE:
		return map[TokenType]struct{}{
			WHERE: {},
			FROM:  {},
		}
	case INSERT:
		return map[TokenType]struct{}{
			VALUES: {},
		}
	case VALUES:
		return map[TokenType]struct{}{
			UPDATE:    {},
			RETURNING: {},
		}
	case FUNCTION:
		return map[TokenType]struct{}{ENDPARENTHESIS: {}}
	case TYPE:
		return map[TokenType]struct{}{ENDPARENTHESIS: {}}
	case LOCK:
		return map[TokenType]struct{}{
			EOF: {},
		}
	case WITH:
		return map[TokenType]struct{}{
			EOF: {},
		}
	default:
		return nil
	}
}

// token types that contain the keyword to make subGroup.
var (
	TokenTypesOfGroupMaker []TokenType
	TokenTypesOfJoinMaker  []TokenType
	TokenTypeOfTieClause   []TokenType
	TokenTypeOfLimitClause []TokenType
)

func init() {
	// TODO: those defs should belong to the parser pkg

	TokenTypesOfGroupMaker = []TokenType{
		SELECT, CASE, FROM, WHERE, ORDER, GROUP, LIMIT,
		ANDGROUP, ORGROUP, HAVING,
		UNION, EXCEPT, INTERSECT,
		FUNCTION,
		STARTPARENTHESIS,
		TYPE,
	}
	TokenTypesOfJoinMaker = []TokenType{
		JOIN, INNER, OUTER, LEFT, RIGHT, NATURAL, CROSS, LATERAL,
	}
	TokenTypeOfTieClause = []TokenType{UNION, INTERSECT, EXCEPT}
	TokenTypeOfLimitClause = []TokenType{LIMIT, FETCH, OFFSET}
}

// IsJoinStart determines if ttype is included in TokenTypesOfJoinMaker.
func (t Token) IsJoinStart() bool {
	for _, v := range TokenTypesOfJoinMaker {
		if t.Type == v {
			return true
		}
	}

	return false
}

// IsTieClauseStart determines if ttype is included in TokenTypesOfTieClause.
func (t Token) IsTieClauseStart() bool {
	for _, v := range TokenTypeOfTieClause {
		if t.Type == v {
			return true
		}
	}

	return false
}

// IsLimitClauseStart determines ttype is included in TokenTypesOfLimitClause.
func (t Token) IsLimitClauseStart() bool {
	for _, v := range TokenTypeOfLimitClause {
		if t.Type == v {
			return true
		}
	}

	return false
}

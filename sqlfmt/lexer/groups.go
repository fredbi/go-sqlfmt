package lexer

// EndTokenTypes returns the corresponding end of clause tokens, or nil if this token does not define a group.
func (t Token) EndTokenTypes() map[TokenType]struct{} {
	switch t.Type {
	case SELECT:
		return map[TokenType]struct{}{
			FROM:  {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
		}
	case FROM:
		return map[TokenType]struct{}{
			WHERE: {},
			INNER: {}, OUTER: {}, LEFT: {}, RIGHT: {}, JOIN: {},
			NATURAL: {}, CROSS: {}, LATERAL: {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
			ORDER: {}, GROUP: {},
			LIMIT: {}, OFFSET: {}, FETCH: {},
			ENDPARENTHESIS: {},
		}
	case CASE:
		return map[TokenType]struct{}{END: {}}
	case JOIN, INNER, OUTER, LEFT, RIGHT, NATURAL, CROSS:
		return map[TokenType]struct{}{
			WHERE: {},
			ORDER: {}, GROUP: {},
			LIMIT: {}, OFFSET: {}, FETCH: {},
			ANDGROUP: {}, ORGROUP: {},
			LEFT: {}, RIGHT: {}, INNER: {}, OUTER: {},
			NATURAL: {}, CROSS: {}, LATERAL: {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
			ENDPARENTHESIS: {},
		}
	case WHERE:
		return map[TokenType]struct{}{
			GROUP: {}, ORDER: {},
			LIMIT: {}, OFFSET: {}, FETCH: {},
			ANDGROUP: {}, OR: {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
			RETURNING:      {},
			ENDPARENTHESIS: {},
		}
	case ANDGROUP: // is this useful?
		return map[TokenType]struct{}{
			GROUP: {}, ORDER: {},
			LIMIT: {}, OFFSET: {}, FETCH: {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
			ANDGROUP: {}, ORGROUP: {},
			ENDPARENTHESIS: {},
		}
	case ORGROUP: // is this useful?
		return map[TokenType]struct{}{
			GROUP: {}, ORDER: {},
			LIMIT: {}, OFFSET: {}, FETCH: {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
			ANDGROUP: {}, ORGROUP: {},
			ENDPARENTHESIS: {},
		}
	case GROUP:
		return map[TokenType]struct{}{
			ORDER: {},
			LIMIT: {}, FETCH: {}, OFFSET: {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
			HAVING:         {},
			ENDPARENTHESIS: {},
		}
	case HAVING:
		return map[TokenType]struct{}{
			ORDER: {},
			LIMIT: {}, OFFSET: {}, FETCH: {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
			ENDPARENTHESIS: {},
		}
	case ORDER:
		return map[TokenType]struct{}{
			GROUP: {},
			LIMIT: {}, FETCH: {}, OFFSET: {},
			UNION: {}, INTERSECT: {}, EXCEPT: {},
			ENDPARENTHESIS: {},
		}
	case LIMIT, FETCH, OFFSET:
		return map[TokenType]struct{}{
			UNION: {}, INTERSECT: {}, EXCEPT: {},
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

func GroupMakers() map[TokenType]struct{} {
	return map[TokenType]struct{}{
		SELECT: {},
		CASE:   {}, FROM: {},
		WHERE: {},
		ORDER: {}, GROUP: {},
		LIMIT: {},
		UNION: {}, INTERSECT: {}, EXCEPT: {},
		ANDGROUP: {}, ORGROUP: {}, HAVING: {},
		FUNCTION:         {},
		STARTPARENTHESIS: {},
		TYPE:             {},
	}
}

func JoinMakers() map[TokenType]struct{} {
	return map[TokenType]struct{}{
		JOIN: {}, INNER: {}, OUTER: {},
		LEFT: {}, RIGHT: {},
		NATURAL: {}, CROSS: {}, LATERAL: {},
	}
}

func TieMakers() map[TokenType]struct{} {
	return map[TokenType]struct{}{
		UNION: {}, INTERSECT: {}, EXCEPT: {},
	}
}

func LimitMakers() map[TokenType]struct{} {
	return map[TokenType]struct{}{
		LIMIT: {}, FETCH: {}, OFFSET: {},
	}
}

// IsKeywordInSelect returns true if token is a keyword in select group.
func (t Token) IsKeywordInSelect() bool {
	return t.Type == SELECT ||
		t.Type == EXISTS ||
		t.Type == DISTINCT ||
		t.Type == DISTINCTROW ||
		t.Type == INTO ||
		t.Type == AS ||
		t.Type == GROUP ||
		t.Type == ORDER ||
		t.Type == BY ||
		t.Type == ON ||
		t.Type == RETURNING ||
		t.Type == SET ||
		t.Type == UPDATE
}

/*
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
*/

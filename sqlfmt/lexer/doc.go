// Package lexer is a lexer for SQL.
//
// By default, it is equipped with all tokens to parse Postgres SQL.
// Postgres default covers the uuid-ossp extension.
//
// The package may optionally register the postgis extension.
//
// TODO:
//  * [x] operators: +, *, /, -, <>, !=, ->, @,  ||, ,...
//  * [x] reserved values (e.g. TRUE, FALSE, TIMESTAMP, Infinity, -Infinity, NaN)
//  * [x] literals, e.g. &U(xxx), B(xyz)
//  * [x] multi-token types (DOUBLE PRECISION, CHARACTER VARYING
//  * [x] register extensions
//  * [x] postgis types and functions
//  * [x] replace maps by prefix keys
//  * sql comments
//  * ambiguity when functions are called without parenthesis (e.g. current_timestamp() vs current_timestamp)
//  * postgres advanced quoting ($$, nested quoting...): means that some newlines bear some parsing semantics OR add special "composite" multiline quoted string token?
//  * DOMAIN
//  * more tests on edge cases (eof, unterminated quotes, etc)
//  * more standard postgres extensions e.g. tgm...
//
// Known unsupported constructs:
//   * Two string constants that are only separated by whitespace with at least one newline are concatenated and effectively treated as if the string had been written as one constant.
// .
//
// Reference lexemes: https://github.com/postgres/postgres/blob/master/src/include/parser/kwlist.h
//
// NOTES:
// * postgres' lexer manages to avoid any backtracking. We do some, when scanning operators (max stack: 4-5 runes)
package lexer

// Package lexer is a lexer for SQL.
//
// By default, it is equipped with all tokens to parse Postgres SQL.
//
// TODO:
//  * [x] operators: +, *, /, -, <>, !=, ->, @,  ||, ,...
//  * [x] reserved values (e.g. TRUE, FALSE, TIMESTAMP, Infinity, -Infinity, NaN)
//  * [x] literals, e.g. &U(xxx), B(xyz)
//  * [x] multi-token types (DOUBLE PRECISION, CHARACTER VARYING
//  * register extensions
//  * sql comments
//  * postgis types and functions
//  * ambiguity when functions are called without parenthesis (e.g. current_timestamp() vs current_timestamp)
//  * postgres advanced quoting
//  * replace maps by prefix keys
//  * DOMAIN
//
// Known unsupported constructs:
//   * Two string constants that are only separated by whitespace with at least one newline are concatenated and effectively treated as if the string had been written as one constant.
// .
package lexer

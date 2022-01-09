# SQL formatting rules
We need to rationalize the indentation rules.

## Casing rules
   * Casing rules are handled by the lexer - there is no contextualization
   * By default all known lexemes are upper-cased (keywords, functions, types and special values such as NULL, TRUE, ...)
   * Identifiers, quoted or not, are not recased by default.

## Coloring rules
   * Coloring rules are handled by the lexer - there is no contextualization

## Indentation rules
* We need groups that know how to:
  * prepend a line feed to the group

* Groups preceded by a line feed:
    * `DO ...`
    * subqueries: `(SELECT ...)`
    * `WHERE ...`
    * FROM ...
    * CASE ... END
    * first column in a SELECT (TODO)
    * first member in a WHERE clause (TODO)
    * multiple VALUES() (TODO)

  Ex:
  ```
  WITH alias AS (
    SELECT
      x,
      y
    FROM
      t
    WHERE
      x IS NULL AND
      y >= 2
    ), another AS (
    SELECT
      b,
      CASE
        WHEN x = 1 THEN x+1
        ELSE x
      END AS c
    FROM
      'aa'
    )
    SELECT
      alias.a,
      another.*
    FROM
      alias
    JOIN another ON
      alias.a = another.b

* Spacing rules:
    * default indentation increment: double space (TODO: configurable)
    * operators, identifiers and other strings: single space (TODO: configurable): `a + b`
    * `.` is not an operator but a character in some identifier
    * cast operator: no spacing, e.g. `a::JSONB`
    * function calls: `FUNC(x, y)`, nested: `MIN(MAX(a))` (TODO: flex)
    * general: single space `WITH {ident} AS ...`
 
* SQL expression layout rules:
    * left-justified style:
      ```
      a
      AND x
      AND y
      OR z
      ```
    * right-justified style:
      ```
      a AND
      x AND
      y OR
      z
      ```
    * parenthesis in expressions (`WHERE`, `ON`)
      ```
      a AND 
      x AND (
        y OR
        z
      )
      ```
      or:
      ```
      a
      AND x
      AND (
        y
        OR z
      )
      ```
* Comma layout rules:
    * left-justified style:
      ```
       , a
       , b
       [...]
      ```
    * right-justified style:
      ```
       a,
       b,
       [...]
      ```

* Flex indentation (TODO)
  Flex-indentation refers to dynamically deciding whether to pack a group on a single line or to split it over different lines.

  Short clauses such as `WITH x AS (SELECT now())` could be packed on a single line,
  whereas longer ones would be split. The cut-off width should be configurable so we may end up in a different formatting optimized for a given terminal width.

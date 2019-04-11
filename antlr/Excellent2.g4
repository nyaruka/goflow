grammar Excellent2;

// rebuild with % antlr -Dlanguage=Go Excellent2.g4 -o ../excellent/gen -package gen -visitor

import LexUnicode;

COMMA: ',';
LPAREN: '(';
RPAREN: ')';
LBRACK: '[';
RBRACK: ']';

DOT: '.';

PLUS: '+';
MINUS: '-';
TIMES: '*';
DIVIDE: '/';
EXPONENT: '^';

EQ: '=';
NEQ: '!=';

LTE: '<=';
LT: '<';
GTE: '>=';
GT: '>';

AMPERSAND: '&';

TEXT: '"' (~["] | '\\"')* '"';
NUMBER: [0-9]+ ('.' [0-9]+)?;

TRUE: [Tt][Rr][Uu][Ee];
FALSE: [Ff][Aa][Ll][Ss][Ee];
NULL: [Nn][Uu][Ll][Ll];

NAME: UnicodeLetter+ (UnicodeLetter | UnicodeDigit | '_')*;
// variable names, e.g. contact.name or function names, e.g. SUM

WS: [ \t\n\r]+ -> skip; // ignore whitespace

ERROR: .;

parse: expression EOF;

atom:
	atom LPAREN parameters? RPAREN	# functionCall
	| atom DOT NAME					# dotLookup
	| atom LBRACK expression RBRACK	# arrayLookup
	| NAME							# contextReference
	| TEXT							# textLiteral
	| NUMBER						# numberLiteral
	| TRUE							# true
	| FALSE							# false
	| NULL							# null;

expression:
	atom												# atomReference
	| MINUS expression									# negation
	| expression EXPONENT expression					# exponent
	| expression op = (TIMES | DIVIDE) expression		# multiplicationOrDivision
	| expression op = (PLUS | MINUS) expression			# additionOrSubtraction
	| expression op = (LTE | LT | GTE | GT) expression	# comparison
	| expression op = (EQ | NEQ) expression				# equality
	| expression AMPERSAND expression					# concatenation
	| LPAREN expression RPAREN							# parentheses;

parameters: expression (COMMA expression)* # functionParameters;
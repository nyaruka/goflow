grammar Excellent3;

// rebuild with % antlr -Dlanguage=Go Excellent3.g4 -o ../excellent/gen -package gen -visitor

import LexUnicode;

COMMA: ',';
LPAREN: '(';
RPAREN: ')';
LBRACK: '[';
RBRACK: ']';

DOT: '.';
ARROW: '=>';

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
INTEGER: [0-9]+;
DECIMAL: [0-9]+ '.' [0-9]+;

TRUE: [Tt][Rr][Uu][Ee];
FALSE: [Ff][Aa][Ll][Ss][Ee];
NULL: [Nn][Uu][Ll][Ll];

NAME: (UnicodeLetter | '_')+ (UnicodeLetter | UnicodeDigit | '_')*;

WS: [ \t\n\r]+ -> skip; // ignore whitespace

ERROR: .;

parse: expression EOF;

expression:
	atom												# atomReference
	| MINUS expression									# negation
	| expression EXPONENT expression					# exponent
	| expression op = (TIMES | DIVIDE) expression		# multiplicationOrDivision
	| expression op = (PLUS | MINUS) expression			# additionOrSubtraction
	| expression op = (LTE | LT | GTE | GT) expression	# comparison
	| expression op = (EQ | NEQ) expression				# equality
	| expression AMPERSAND expression					# concatenation
	| LPAREN nameList RPAREN ARROW expression			# anonFunction
	| TEXT												# textLiteral
	| (INTEGER | DECIMAL)								# numberLiteral
	| TRUE												# true
	| FALSE												# false
	| NULL												# null;

// a subset of expressions which can be followed by (), [] or .
atom:
	atom LPAREN parameters? RPAREN	# functionCall
	| atom DOT (NAME | INTEGER)		# dotLookup
	| atom LBRACK expression RBRACK	# arrayLookup
	| LPAREN expression RPAREN		# parentheses
	| NAME							# contextReference;

parameters: expression (COMMA expression)* # functionParameters;

nameList: NAME (COMMA NAME)*;
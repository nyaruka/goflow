grammar Excellent1;

// rebuild with % antlr -Dlanguage=Go Excellent1.g4 -o ../flows/definition/legacy/gen -package gen -visitor

import LexUnicode;

COMMA: ',';
LPAREN: '(';
RPAREN: ')';

PLUS: '+';
MINUS: '-';
TIMES: '*';
DIVIDE: '/';
EXPONENT: '^';

EQ: '=';
NEQ: '<>';

LTE: '<=';
LT: '<';
GTE: '>=';
GT: '>';

AMPERSAND: '&';

DECIMAL: [0-9]+ ('.' [0-9]+)?;
STRING: '"' (~["] | '""')* '"';

TRUE: [Tt][Rr][Uu][Ee];
FALSE: [Ff][Aa][Ll][Ss][Ee];

NAME:
	UnicodeLetter+ (UnicodeLetter | UnicodeDigit | '_' | '.')*;

WS: [ \t\n\r]+ -> skip; // ignore whitespace

ERROR: .;

parse: expression EOF;

expression:
	fnname LPAREN parameters? RPAREN					# functionCall
	| MINUS expression									# negation
	| expression EXPONENT expression					# exponentExpression
	| expression op = (TIMES | DIVIDE) expression		# multiplicationOrDivisionExpression
	| expression op = (PLUS | MINUS) expression			# additionOrSubtractionExpression
	| expression op = (LTE | LT | GTE | GT) expression	# comparisonExpression
	| expression op = (EQ | NEQ) expression				# equalityExpression
	| expression AMPERSAND expression					# concatenation
	| STRING											# stringLiteral
	| DECIMAL											# decimalLiteral
	| TRUE												# true
	| FALSE												# false
	| NAME												# contextReference
	| LPAREN expression RPAREN							# parentheses;

fnname: NAME | TRUE | FALSE;

parameters: expression (COMMA expression)* # functionParameters;
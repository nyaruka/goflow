grammar ContactQL;

// rebuild with % antlr -Dlanguage=Go ContactQL.g4 -o ../contactql/gen -package gen -visitor

import LexUnicode;

fragment HAS: [Hh][Aa][Ss];
fragment IS: [Ii][Ss];

LPAREN: '(';
RPAREN: ')';
AND: [Aa][Nn][Dd];
OR: [Oo][Rr];
COMPARATOR: (
		'='
		| '!='
		| '~'
		| '>='
		| '<='
		| '>'
		| '<'
		| HAS
		| IS
	);
TEXT: (
		UnicodeLetter
		| UnicodeDigit
		| '_'
		| '.'
		| '-'
		| '+'
		| '/'
		| '\''
		| '@'
		| ':'
	)+;
STRING: '"' (~["] | '\\"')* '"';

WS: [ \t\n\r]+ -> skip; // ignore whitespace

ERROR: .;

parse: expression EOF;

expression:
	expression AND expression	# combinationAnd
	| expression expression		# combinationImpicitAnd
	| expression OR expression	# combinationOr
	| LPAREN expression RPAREN	# expressionGrouping
	| TEXT COMPARATOR literal	# condition
	| literal					# implicitCondition;

literal: TEXT # textLiteral | STRING # stringLiteral;
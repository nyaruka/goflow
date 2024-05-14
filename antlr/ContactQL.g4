grammar ContactQL;

import LexUnicode;

// Lexer rules
fragment HAS: [Hh][Aa][Ss];
fragment IS: [Ii][Ss];
fragment PROPTYPE: (UnicodeLetter)+;
fragment PROPKEY: (UnicodeLetter | UnicodeDigit | '_')+;

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
STRING: '"' (~["] | '\\"')* '"';
PROPERTY: (PROPTYPE '.')? PROPKEY;
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

WS: [ \t\n\r]+ -> skip; // ignore whitespace

ERROR: .;

// Parser rules
parse: expression EOF;

expression:
	expression AND expression		# combinationAnd
	| expression expression			# combinationImpicitAnd
	| expression OR expression		# combinationOr
	| LPAREN expression RPAREN		# expressionGrouping
	| PROPERTY COMPARATOR literal	# condition
	| literal						# implicitCondition;

literal:
	PROPERTY	# textLiteral // it's not really a property, just indistinguishable by lexer
	| TEXT		# textLiteral
	| STRING	# stringLiteral;
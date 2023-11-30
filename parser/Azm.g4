grammar Azm;

relation
    :   unionRel EOF
    ;

permission
    :   unionPerm EOF           # ToUnionPerm
    |   intersectionPerm EOF    # ToIntersectionPerm
    |   exclusionPerm EOF       # ToExclusionPerm
    ;

unionRel
    :   rel ('|' rel)*
    ;

unionPerm
    :   perm ('|' perm)*
    ;

intersectionPerm
    :   perm ('&' perm)*
    ;

exclusionPerm
    :   perm '-' perm
    ;

rel
    :   singleRel       # ToSingleRel
    |   wildcardRel     # ToWildcardRel
    |   subjectRel      # ToSubjectRel
    |   arrowRel        # ToArrowRel
    ;

perm
    :   singleRel       # ToSinglePerm
    |   arrowRel        # ToArrowPerm
    ;

singleRel
    :   ID
    ;

subjectRel
    :   ID HASH ID
    ;

wildcardRel
    :   ID COLON ASTERISK
    ;

arrowRel
    :   ID ARROW ID
    ;

ARROW:
    '-''>' ;

HASH:
    '#' ;

COLON:
    ':' ;

ASTERISK:
    '*';

ID: [a-z][a-z0-9._]*[a-z0-9] ;

WS: [ \t\n\r\f]+ -> skip ;

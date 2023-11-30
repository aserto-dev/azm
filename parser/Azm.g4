grammar Azm;

relation
    :   rel ('|' rel)*  EOF
    ;

permission
    :   union EOF           #UnionPerm
    |   intersection EOF    #IntersectionPerm
    |   exclusion EOF       #ExclusionPerm
    ;

union
    :   perm ('|' perm)*
    ;

intersection
    :   perm ('&' perm)*
    ;

exclusion
    :   perm '-' perm
    ;

rel
    :   single      #SingleRel
    |   wildcard    #WildcardRel
    |   subject     #SubjectRel
    |   arrow       #ArrowRel
    ;

perm
    :   single      #SinglePerm
    |   arrow       #ArrowPerm
    ;

single
    :   ID
    ;

subject
    :   ID HASH ID
    ;

wildcard
    :   ID COLON ASTERISK
    ;

arrow
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

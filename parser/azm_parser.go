// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type AzmParser struct {
	*antlr.BaseParser
}

var AzmParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func azmParserInit() {
	staticData := &AzmParserStaticData
	staticData.LiteralNames = []string{
		"", "'|'", "'&'", "'-'", "", "'#'", "':'", "'*'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "ARROW", "HASH", "COLON", "ASTERISK", "ID", "WS",
	}
	staticData.RuleNames = []string{
		"relation", "permission", "unionPerm", "intersectionPerm", "exclusionPerm",
		"rel", "perm", "single", "subject", "wildcard", "arrow",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 9, 88, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 4,
		2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 10,
		1, 0, 1, 0, 1, 0, 5, 0, 26, 8, 0, 10, 0, 12, 0, 29, 9, 0, 1, 0, 1, 0, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 42, 8, 1, 1, 2,
		1, 2, 1, 2, 5, 2, 47, 8, 2, 10, 2, 12, 2, 50, 9, 2, 1, 3, 1, 3, 1, 3, 5,
		3, 55, 8, 3, 10, 3, 12, 3, 58, 9, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5,
		1, 5, 1, 5, 3, 5, 68, 8, 5, 1, 6, 1, 6, 3, 6, 72, 8, 6, 1, 7, 1, 7, 1,
		8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10,
		1, 10, 0, 0, 11, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 0, 0, 85, 0, 22,
		1, 0, 0, 0, 2, 41, 1, 0, 0, 0, 4, 43, 1, 0, 0, 0, 6, 51, 1, 0, 0, 0, 8,
		59, 1, 0, 0, 0, 10, 67, 1, 0, 0, 0, 12, 71, 1, 0, 0, 0, 14, 73, 1, 0, 0,
		0, 16, 75, 1, 0, 0, 0, 18, 79, 1, 0, 0, 0, 20, 83, 1, 0, 0, 0, 22, 27,
		3, 10, 5, 0, 23, 24, 5, 1, 0, 0, 24, 26, 3, 10, 5, 0, 25, 23, 1, 0, 0,
		0, 26, 29, 1, 0, 0, 0, 27, 25, 1, 0, 0, 0, 27, 28, 1, 0, 0, 0, 28, 30,
		1, 0, 0, 0, 29, 27, 1, 0, 0, 0, 30, 31, 5, 0, 0, 1, 31, 1, 1, 0, 0, 0,
		32, 33, 3, 4, 2, 0, 33, 34, 5, 0, 0, 1, 34, 42, 1, 0, 0, 0, 35, 36, 3,
		6, 3, 0, 36, 37, 5, 0, 0, 1, 37, 42, 1, 0, 0, 0, 38, 39, 3, 8, 4, 0, 39,
		40, 5, 0, 0, 1, 40, 42, 1, 0, 0, 0, 41, 32, 1, 0, 0, 0, 41, 35, 1, 0, 0,
		0, 41, 38, 1, 0, 0, 0, 42, 3, 1, 0, 0, 0, 43, 48, 3, 12, 6, 0, 44, 45,
		5, 1, 0, 0, 45, 47, 3, 12, 6, 0, 46, 44, 1, 0, 0, 0, 47, 50, 1, 0, 0, 0,
		48, 46, 1, 0, 0, 0, 48, 49, 1, 0, 0, 0, 49, 5, 1, 0, 0, 0, 50, 48, 1, 0,
		0, 0, 51, 56, 3, 12, 6, 0, 52, 53, 5, 2, 0, 0, 53, 55, 3, 12, 6, 0, 54,
		52, 1, 0, 0, 0, 55, 58, 1, 0, 0, 0, 56, 54, 1, 0, 0, 0, 56, 57, 1, 0, 0,
		0, 57, 7, 1, 0, 0, 0, 58, 56, 1, 0, 0, 0, 59, 60, 3, 12, 6, 0, 60, 61,
		5, 3, 0, 0, 61, 62, 3, 12, 6, 0, 62, 9, 1, 0, 0, 0, 63, 68, 3, 14, 7, 0,
		64, 68, 3, 18, 9, 0, 65, 68, 3, 16, 8, 0, 66, 68, 3, 20, 10, 0, 67, 63,
		1, 0, 0, 0, 67, 64, 1, 0, 0, 0, 67, 65, 1, 0, 0, 0, 67, 66, 1, 0, 0, 0,
		68, 11, 1, 0, 0, 0, 69, 72, 3, 14, 7, 0, 70, 72, 3, 20, 10, 0, 71, 69,
		1, 0, 0, 0, 71, 70, 1, 0, 0, 0, 72, 13, 1, 0, 0, 0, 73, 74, 5, 8, 0, 0,
		74, 15, 1, 0, 0, 0, 75, 76, 5, 8, 0, 0, 76, 77, 5, 5, 0, 0, 77, 78, 5,
		8, 0, 0, 78, 17, 1, 0, 0, 0, 79, 80, 5, 8, 0, 0, 80, 81, 5, 6, 0, 0, 81,
		82, 5, 7, 0, 0, 82, 19, 1, 0, 0, 0, 83, 84, 5, 8, 0, 0, 84, 85, 5, 4, 0,
		0, 85, 86, 5, 8, 0, 0, 86, 21, 1, 0, 0, 0, 6, 27, 41, 48, 56, 67, 71,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// AzmParserInit initializes any static state used to implement AzmParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewAzmParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func AzmParserInit() {
	staticData := &AzmParserStaticData
	staticData.once.Do(azmParserInit)
}

// NewAzmParser produces a new parser instance for the optional input antlr.TokenStream.
func NewAzmParser(input antlr.TokenStream) *AzmParser {
	AzmParserInit()
	this := new(AzmParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &AzmParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "Azm.g4"

	return this
}

// AzmParser tokens.
const (
	AzmParserEOF      = antlr.TokenEOF
	AzmParserT__0     = 1
	AzmParserT__1     = 2
	AzmParserT__2     = 3
	AzmParserARROW    = 4
	AzmParserHASH     = 5
	AzmParserCOLON    = 6
	AzmParserASTERISK = 7
	AzmParserID       = 8
	AzmParserWS       = 9
)

// AzmParser rules.
const (
	AzmParserRULE_relation         = 0
	AzmParserRULE_permission       = 1
	AzmParserRULE_unionPerm        = 2
	AzmParserRULE_intersectionPerm = 3
	AzmParserRULE_exclusionPerm    = 4
	AzmParserRULE_rel              = 5
	AzmParserRULE_perm             = 6
	AzmParserRULE_single           = 7
	AzmParserRULE_subject          = 8
	AzmParserRULE_wildcard         = 9
	AzmParserRULE_arrow            = 10
)

// IRelationContext is an interface to support dynamic dispatch.
type IRelationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRel() []IRelContext
	Rel(i int) IRelContext
	EOF() antlr.TerminalNode

	// IsRelationContext differentiates from other interfaces.
	IsRelationContext()
}

type RelationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelationContext() *RelationContext {
	var p = new(RelationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_relation
	return p
}

func InitEmptyRelationContext(p *RelationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_relation
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_relation

	return p
}

func (s *RelationContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationContext) AllRel() []IRelContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelContext); ok {
			len++
		}
	}

	tst := make([]IRelContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelContext); ok {
			tst[i] = t.(IRelContext)
			i++
		}
	}

	return tst
}

func (s *RelationContext) Rel(i int) IRelContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelContext)
}

func (s *RelationContext) EOF() antlr.TerminalNode {
	return s.GetToken(AzmParserEOF, 0)
}

func (s *RelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RelationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterRelation(s)
	}
}

func (s *RelationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitRelation(s)
	}
}

func (s *RelationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitRelation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Relation() (localctx IRelationContext) {
	localctx = NewRelationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, AzmParserRULE_relation)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(22)
		p.Rel()
	}
	p.SetState(27)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__0 {
		{
			p.SetState(23)
			p.Match(AzmParserT__0)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(24)
			p.Rel()
		}

		p.SetState(29)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(30)
		p.Match(AzmParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPermissionContext is an interface to support dynamic dispatch.
type IPermissionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UnionPerm() IUnionPermContext
	EOF() antlr.TerminalNode
	IntersectionPerm() IIntersectionPermContext
	ExclusionPerm() IExclusionPermContext

	// IsPermissionContext differentiates from other interfaces.
	IsPermissionContext()
}

type PermissionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPermissionContext() *PermissionContext {
	var p = new(PermissionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_permission
	return p
}

func InitEmptyPermissionContext(p *PermissionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_permission
}

func (*PermissionContext) IsPermissionContext() {}

func NewPermissionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PermissionContext {
	var p = new(PermissionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_permission

	return p
}

func (s *PermissionContext) GetParser() antlr.Parser { return s.parser }

func (s *PermissionContext) UnionPerm() IUnionPermContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionPermContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionPermContext)
}

func (s *PermissionContext) EOF() antlr.TerminalNode {
	return s.GetToken(AzmParserEOF, 0)
}

func (s *PermissionContext) IntersectionPerm() IIntersectionPermContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntersectionPermContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntersectionPermContext)
}

func (s *PermissionContext) ExclusionPerm() IExclusionPermContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExclusionPermContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExclusionPermContext)
}

func (s *PermissionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PermissionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PermissionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterPermission(s)
	}
}

func (s *PermissionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitPermission(s)
	}
}

func (s *PermissionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitPermission(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Permission() (localctx IPermissionContext) {
	localctx = NewPermissionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, AzmParserRULE_permission)
	p.SetState(41)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(32)
			p.UnionPerm()
		}
		{
			p.SetState(33)
			p.Match(AzmParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(35)
			p.IntersectionPerm()
		}
		{
			p.SetState(36)
			p.Match(AzmParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(38)
			p.ExclusionPerm()
		}
		{
			p.SetState(39)
			p.Match(AzmParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnionPermContext is an interface to support dynamic dispatch.
type IUnionPermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPerm() []IPermContext
	Perm(i int) IPermContext

	// IsUnionPermContext differentiates from other interfaces.
	IsUnionPermContext()
}

type UnionPermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnionPermContext() *UnionPermContext {
	var p = new(UnionPermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_unionPerm
	return p
}

func InitEmptyUnionPermContext(p *UnionPermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_unionPerm
}

func (*UnionPermContext) IsUnionPermContext() {}

func NewUnionPermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionPermContext {
	var p = new(UnionPermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_unionPerm

	return p
}

func (s *UnionPermContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionPermContext) AllPerm() []IPermContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPermContext); ok {
			len++
		}
	}

	tst := make([]IPermContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPermContext); ok {
			tst[i] = t.(IPermContext)
			i++
		}
	}

	return tst
}

func (s *UnionPermContext) Perm(i int) IPermContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPermContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPermContext)
}

func (s *UnionPermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionPermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionPermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterUnionPerm(s)
	}
}

func (s *UnionPermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitUnionPerm(s)
	}
}

func (s *UnionPermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitUnionPerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) UnionPerm() (localctx IUnionPermContext) {
	localctx = NewUnionPermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, AzmParserRULE_unionPerm)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(43)
		p.Perm()
	}
	p.SetState(48)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__0 {
		{
			p.SetState(44)
			p.Match(AzmParserT__0)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(45)
			p.Perm()
		}

		p.SetState(50)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIntersectionPermContext is an interface to support dynamic dispatch.
type IIntersectionPermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPerm() []IPermContext
	Perm(i int) IPermContext

	// IsIntersectionPermContext differentiates from other interfaces.
	IsIntersectionPermContext()
}

type IntersectionPermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntersectionPermContext() *IntersectionPermContext {
	var p = new(IntersectionPermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_intersectionPerm
	return p
}

func InitEmptyIntersectionPermContext(p *IntersectionPermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_intersectionPerm
}

func (*IntersectionPermContext) IsIntersectionPermContext() {}

func NewIntersectionPermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntersectionPermContext {
	var p = new(IntersectionPermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_intersectionPerm

	return p
}

func (s *IntersectionPermContext) GetParser() antlr.Parser { return s.parser }

func (s *IntersectionPermContext) AllPerm() []IPermContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPermContext); ok {
			len++
		}
	}

	tst := make([]IPermContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPermContext); ok {
			tst[i] = t.(IPermContext)
			i++
		}
	}

	return tst
}

func (s *IntersectionPermContext) Perm(i int) IPermContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPermContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPermContext)
}

func (s *IntersectionPermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntersectionPermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntersectionPermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterIntersectionPerm(s)
	}
}

func (s *IntersectionPermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitIntersectionPerm(s)
	}
}

func (s *IntersectionPermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitIntersectionPerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) IntersectionPerm() (localctx IIntersectionPermContext) {
	localctx = NewIntersectionPermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, AzmParserRULE_intersectionPerm)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(51)
		p.Perm()
	}
	p.SetState(56)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__1 {
		{
			p.SetState(52)
			p.Match(AzmParserT__1)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(53)
			p.Perm()
		}

		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExclusionPermContext is an interface to support dynamic dispatch.
type IExclusionPermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPerm() []IPermContext
	Perm(i int) IPermContext

	// IsExclusionPermContext differentiates from other interfaces.
	IsExclusionPermContext()
}

type ExclusionPermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExclusionPermContext() *ExclusionPermContext {
	var p = new(ExclusionPermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_exclusionPerm
	return p
}

func InitEmptyExclusionPermContext(p *ExclusionPermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_exclusionPerm
}

func (*ExclusionPermContext) IsExclusionPermContext() {}

func NewExclusionPermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExclusionPermContext {
	var p = new(ExclusionPermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_exclusionPerm

	return p
}

func (s *ExclusionPermContext) GetParser() antlr.Parser { return s.parser }

func (s *ExclusionPermContext) AllPerm() []IPermContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPermContext); ok {
			len++
		}
	}

	tst := make([]IPermContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPermContext); ok {
			tst[i] = t.(IPermContext)
			i++
		}
	}

	return tst
}

func (s *ExclusionPermContext) Perm(i int) IPermContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPermContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPermContext)
}

func (s *ExclusionPermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExclusionPermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExclusionPermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterExclusionPerm(s)
	}
}

func (s *ExclusionPermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitExclusionPerm(s)
	}
}

func (s *ExclusionPermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitExclusionPerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) ExclusionPerm() (localctx IExclusionPermContext) {
	localctx = NewExclusionPermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, AzmParserRULE_exclusionPerm)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(59)
		p.Perm()
	}
	{
		p.SetState(60)
		p.Match(AzmParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(61)
		p.Perm()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRelContext is an interface to support dynamic dispatch.
type IRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsRelContext differentiates from other interfaces.
	IsRelContext()
}

type RelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelContext() *RelContext {
	var p = new(RelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_rel
	return p
}

func InitEmptyRelContext(p *RelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_rel
}

func (*RelContext) IsRelContext() {}

func NewRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelContext {
	var p = new(RelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_rel

	return p
}

func (s *RelContext) GetParser() antlr.Parser { return s.parser }

func (s *RelContext) CopyAll(ctx *RelContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *RelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SingleRelContext struct {
	RelContext
}

func NewSingleRelContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SingleRelContext {
	var p = new(SingleRelContext)

	InitEmptyRelContext(&p.RelContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelContext))

	return p
}

func (s *SingleRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SingleRelContext) Single() ISingleContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleContext)
}

func (s *SingleRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterSingleRel(s)
	}
}

func (s *SingleRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitSingleRel(s)
	}
}

func (s *SingleRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitSingleRel(s)

	default:
		return t.VisitChildren(s)
	}
}

type SubjectRelContext struct {
	RelContext
}

func NewSubjectRelContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SubjectRelContext {
	var p = new(SubjectRelContext)

	InitEmptyRelContext(&p.RelContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelContext))

	return p
}

func (s *SubjectRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubjectRelContext) Subject() ISubjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISubjectContext)
}

func (s *SubjectRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterSubjectRel(s)
	}
}

func (s *SubjectRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitSubjectRel(s)
	}
}

func (s *SubjectRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitSubjectRel(s)

	default:
		return t.VisitChildren(s)
	}
}

type ArrowRelContext struct {
	RelContext
}

func NewArrowRelContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ArrowRelContext {
	var p = new(ArrowRelContext)

	InitEmptyRelContext(&p.RelContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelContext))

	return p
}

func (s *ArrowRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrowRelContext) Arrow() IArrowContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrowContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrowContext)
}

func (s *ArrowRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterArrowRel(s)
	}
}

func (s *ArrowRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitArrowRel(s)
	}
}

func (s *ArrowRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitArrowRel(s)

	default:
		return t.VisitChildren(s)
	}
}

type WildcardRelContext struct {
	RelContext
}

func NewWildcardRelContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *WildcardRelContext {
	var p = new(WildcardRelContext)

	InitEmptyRelContext(&p.RelContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelContext))

	return p
}

func (s *WildcardRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WildcardRelContext) Wildcard() IWildcardContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWildcardContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWildcardContext)
}

func (s *WildcardRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterWildcardRel(s)
	}
}

func (s *WildcardRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitWildcardRel(s)
	}
}

func (s *WildcardRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitWildcardRel(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Rel() (localctx IRelContext) {
	localctx = NewRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, AzmParserRULE_rel)
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSingleRelContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(63)
			p.Single()
		}

	case 2:
		localctx = NewWildcardRelContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(64)
			p.Wildcard()
		}

	case 3:
		localctx = NewSubjectRelContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(65)
			p.Subject()
		}

	case 4:
		localctx = NewArrowRelContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(66)
			p.Arrow()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPermContext is an interface to support dynamic dispatch.
type IPermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPermContext differentiates from other interfaces.
	IsPermContext()
}

type PermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPermContext() *PermContext {
	var p = new(PermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_perm
	return p
}

func InitEmptyPermContext(p *PermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_perm
}

func (*PermContext) IsPermContext() {}

func NewPermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PermContext {
	var p = new(PermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_perm

	return p
}

func (s *PermContext) GetParser() antlr.Parser { return s.parser }

func (s *PermContext) CopyAll(ctx *PermContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *PermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ArrowPermContext struct {
	PermContext
}

func NewArrowPermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ArrowPermContext {
	var p = new(ArrowPermContext)

	InitEmptyPermContext(&p.PermContext)
	p.parser = parser
	p.CopyAll(ctx.(*PermContext))

	return p
}

func (s *ArrowPermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrowPermContext) Arrow() IArrowContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrowContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrowContext)
}

func (s *ArrowPermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterArrowPerm(s)
	}
}

func (s *ArrowPermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitArrowPerm(s)
	}
}

func (s *ArrowPermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitArrowPerm(s)

	default:
		return t.VisitChildren(s)
	}
}

type SinglePermContext struct {
	PermContext
}

func NewSinglePermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SinglePermContext {
	var p = new(SinglePermContext)

	InitEmptyPermContext(&p.PermContext)
	p.parser = parser
	p.CopyAll(ctx.(*PermContext))

	return p
}

func (s *SinglePermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SinglePermContext) Single() ISingleContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleContext)
}

func (s *SinglePermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterSinglePerm(s)
	}
}

func (s *SinglePermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitSinglePerm(s)
	}
}

func (s *SinglePermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitSinglePerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Perm() (localctx IPermContext) {
	localctx = NewPermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, AzmParserRULE_perm)
	p.SetState(71)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSinglePermContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(69)
			p.Single()
		}

	case 2:
		localctx = NewArrowPermContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(70)
			p.Arrow()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISingleContext is an interface to support dynamic dispatch.
type ISingleContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsSingleContext differentiates from other interfaces.
	IsSingleContext()
}

type SingleContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySingleContext() *SingleContext {
	var p = new(SingleContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_single
	return p
}

func InitEmptySingleContext(p *SingleContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_single
}

func (*SingleContext) IsSingleContext() {}

func NewSingleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SingleContext {
	var p = new(SingleContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_single

	return p
}

func (s *SingleContext) GetParser() antlr.Parser { return s.parser }

func (s *SingleContext) ID() antlr.TerminalNode {
	return s.GetToken(AzmParserID, 0)
}

func (s *SingleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SingleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SingleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterSingle(s)
	}
}

func (s *SingleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitSingle(s)
	}
}

func (s *SingleContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitSingle(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Single() (localctx ISingleContext) {
	localctx = NewSingleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, AzmParserRULE_single)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(73)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISubjectContext is an interface to support dynamic dispatch.
type ISubjectContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	HASH() antlr.TerminalNode

	// IsSubjectContext differentiates from other interfaces.
	IsSubjectContext()
}

type SubjectContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubjectContext() *SubjectContext {
	var p = new(SubjectContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_subject
	return p
}

func InitEmptySubjectContext(p *SubjectContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_subject
}

func (*SubjectContext) IsSubjectContext() {}

func NewSubjectContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubjectContext {
	var p = new(SubjectContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_subject

	return p
}

func (s *SubjectContext) GetParser() antlr.Parser { return s.parser }

func (s *SubjectContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(AzmParserID)
}

func (s *SubjectContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(AzmParserID, i)
}

func (s *SubjectContext) HASH() antlr.TerminalNode {
	return s.GetToken(AzmParserHASH, 0)
}

func (s *SubjectContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubjectContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubjectContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterSubject(s)
	}
}

func (s *SubjectContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitSubject(s)
	}
}

func (s *SubjectContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitSubject(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Subject() (localctx ISubjectContext) {
	localctx = NewSubjectContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, AzmParserRULE_subject)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(75)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(76)
		p.Match(AzmParserHASH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(77)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWildcardContext is an interface to support dynamic dispatch.
type IWildcardContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	ASTERISK() antlr.TerminalNode

	// IsWildcardContext differentiates from other interfaces.
	IsWildcardContext()
}

type WildcardContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWildcardContext() *WildcardContext {
	var p = new(WildcardContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_wildcard
	return p
}

func InitEmptyWildcardContext(p *WildcardContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_wildcard
}

func (*WildcardContext) IsWildcardContext() {}

func NewWildcardContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WildcardContext {
	var p = new(WildcardContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_wildcard

	return p
}

func (s *WildcardContext) GetParser() antlr.Parser { return s.parser }

func (s *WildcardContext) ID() antlr.TerminalNode {
	return s.GetToken(AzmParserID, 0)
}

func (s *WildcardContext) COLON() antlr.TerminalNode {
	return s.GetToken(AzmParserCOLON, 0)
}

func (s *WildcardContext) ASTERISK() antlr.TerminalNode {
	return s.GetToken(AzmParserASTERISK, 0)
}

func (s *WildcardContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WildcardContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WildcardContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterWildcard(s)
	}
}

func (s *WildcardContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitWildcard(s)
	}
}

func (s *WildcardContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitWildcard(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Wildcard() (localctx IWildcardContext) {
	localctx = NewWildcardContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, AzmParserRULE_wildcard)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(79)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(80)
		p.Match(AzmParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(81)
		p.Match(AzmParserASTERISK)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArrowContext is an interface to support dynamic dispatch.
type IArrowContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	ARROW() antlr.TerminalNode

	// IsArrowContext differentiates from other interfaces.
	IsArrowContext()
}

type ArrowContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrowContext() *ArrowContext {
	var p = new(ArrowContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_arrow
	return p
}

func InitEmptyArrowContext(p *ArrowContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_arrow
}

func (*ArrowContext) IsArrowContext() {}

func NewArrowContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrowContext {
	var p = new(ArrowContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_arrow

	return p
}

func (s *ArrowContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrowContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(AzmParserID)
}

func (s *ArrowContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(AzmParserID, i)
}

func (s *ArrowContext) ARROW() antlr.TerminalNode {
	return s.GetToken(AzmParserARROW, 0)
}

func (s *ArrowContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrowContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrowContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterArrow(s)
	}
}

func (s *ArrowContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitArrow(s)
	}
}

func (s *ArrowContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitArrow(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Arrow() (localctx IArrowContext) {
	localctx = NewArrowContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, AzmParserRULE_arrow)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(83)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(84)
		p.Match(AzmParserARROW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(85)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

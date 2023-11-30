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
		"relation", "permission", "unionRel", "unionPerm", "intersectionPerm",
		"exclusionPerm", "rel", "perm", "singleRel", "subjectRel", "wildcardRel",
		"arrowRel",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 9, 91, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 4,
		2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 10,
		2, 11, 7, 11, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 3, 1, 37, 8, 1, 1, 2, 1, 2, 1, 2, 5, 2, 42, 8, 2, 10, 2,
		12, 2, 45, 9, 2, 1, 3, 1, 3, 1, 3, 5, 3, 50, 8, 3, 10, 3, 12, 3, 53, 9,
		3, 1, 4, 1, 4, 1, 4, 5, 4, 58, 8, 4, 10, 4, 12, 4, 61, 9, 4, 1, 5, 1, 5,
		1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 71, 8, 6, 1, 7, 1, 7, 3, 7, 75,
		8, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1,
		11, 1, 11, 1, 11, 1, 11, 1, 11, 0, 0, 12, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 0, 0, 87, 0, 24, 1, 0, 0, 0, 2, 36, 1, 0, 0, 0, 4, 38, 1, 0,
		0, 0, 6, 46, 1, 0, 0, 0, 8, 54, 1, 0, 0, 0, 10, 62, 1, 0, 0, 0, 12, 70,
		1, 0, 0, 0, 14, 74, 1, 0, 0, 0, 16, 76, 1, 0, 0, 0, 18, 78, 1, 0, 0, 0,
		20, 82, 1, 0, 0, 0, 22, 86, 1, 0, 0, 0, 24, 25, 3, 4, 2, 0, 25, 26, 5,
		0, 0, 1, 26, 1, 1, 0, 0, 0, 27, 28, 3, 6, 3, 0, 28, 29, 5, 0, 0, 1, 29,
		37, 1, 0, 0, 0, 30, 31, 3, 8, 4, 0, 31, 32, 5, 0, 0, 1, 32, 37, 1, 0, 0,
		0, 33, 34, 3, 10, 5, 0, 34, 35, 5, 0, 0, 1, 35, 37, 1, 0, 0, 0, 36, 27,
		1, 0, 0, 0, 36, 30, 1, 0, 0, 0, 36, 33, 1, 0, 0, 0, 37, 3, 1, 0, 0, 0,
		38, 43, 3, 12, 6, 0, 39, 40, 5, 1, 0, 0, 40, 42, 3, 12, 6, 0, 41, 39, 1,
		0, 0, 0, 42, 45, 1, 0, 0, 0, 43, 41, 1, 0, 0, 0, 43, 44, 1, 0, 0, 0, 44,
		5, 1, 0, 0, 0, 45, 43, 1, 0, 0, 0, 46, 51, 3, 14, 7, 0, 47, 48, 5, 1, 0,
		0, 48, 50, 3, 14, 7, 0, 49, 47, 1, 0, 0, 0, 50, 53, 1, 0, 0, 0, 51, 49,
		1, 0, 0, 0, 51, 52, 1, 0, 0, 0, 52, 7, 1, 0, 0, 0, 53, 51, 1, 0, 0, 0,
		54, 59, 3, 14, 7, 0, 55, 56, 5, 2, 0, 0, 56, 58, 3, 14, 7, 0, 57, 55, 1,
		0, 0, 0, 58, 61, 1, 0, 0, 0, 59, 57, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0, 60,
		9, 1, 0, 0, 0, 61, 59, 1, 0, 0, 0, 62, 63, 3, 14, 7, 0, 63, 64, 5, 3, 0,
		0, 64, 65, 3, 14, 7, 0, 65, 11, 1, 0, 0, 0, 66, 71, 3, 16, 8, 0, 67, 71,
		3, 20, 10, 0, 68, 71, 3, 18, 9, 0, 69, 71, 3, 22, 11, 0, 70, 66, 1, 0,
		0, 0, 70, 67, 1, 0, 0, 0, 70, 68, 1, 0, 0, 0, 70, 69, 1, 0, 0, 0, 71, 13,
		1, 0, 0, 0, 72, 75, 3, 16, 8, 0, 73, 75, 3, 22, 11, 0, 74, 72, 1, 0, 0,
		0, 74, 73, 1, 0, 0, 0, 75, 15, 1, 0, 0, 0, 76, 77, 5, 8, 0, 0, 77, 17,
		1, 0, 0, 0, 78, 79, 5, 8, 0, 0, 79, 80, 5, 5, 0, 0, 80, 81, 5, 8, 0, 0,
		81, 19, 1, 0, 0, 0, 82, 83, 5, 8, 0, 0, 83, 84, 5, 6, 0, 0, 84, 85, 5,
		7, 0, 0, 85, 21, 1, 0, 0, 0, 86, 87, 5, 8, 0, 0, 87, 88, 5, 4, 0, 0, 88,
		89, 5, 8, 0, 0, 89, 23, 1, 0, 0, 0, 6, 36, 43, 51, 59, 70, 74,
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
	AzmParserRULE_unionRel         = 2
	AzmParserRULE_unionPerm        = 3
	AzmParserRULE_intersectionPerm = 4
	AzmParserRULE_exclusionPerm    = 5
	AzmParserRULE_rel              = 6
	AzmParserRULE_perm             = 7
	AzmParserRULE_singleRel        = 8
	AzmParserRULE_subjectRel       = 9
	AzmParserRULE_wildcardRel      = 10
	AzmParserRULE_arrowRel         = 11
)

// IRelationContext is an interface to support dynamic dispatch.
type IRelationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UnionRel() IUnionRelContext
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

func (s *RelationContext) UnionRel() IUnionRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionRelContext)
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
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(24)
		p.UnionRel()
	}
	{
		p.SetState(25)
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

func (s *PermissionContext) CopyAll(ctx *PermissionContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *PermissionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PermissionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ToExclusionPermContext struct {
	PermissionContext
}

func NewToExclusionPermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToExclusionPermContext {
	var p = new(ToExclusionPermContext)

	InitEmptyPermissionContext(&p.PermissionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PermissionContext))

	return p
}

func (s *ToExclusionPermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToExclusionPermContext) ExclusionPerm() IExclusionPermContext {
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

func (s *ToExclusionPermContext) EOF() antlr.TerminalNode {
	return s.GetToken(AzmParserEOF, 0)
}

func (s *ToExclusionPermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToExclusionPerm(s)
	}
}

func (s *ToExclusionPermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToExclusionPerm(s)
	}
}

func (s *ToExclusionPermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToExclusionPerm(s)

	default:
		return t.VisitChildren(s)
	}
}

type ToIntersectionPermContext struct {
	PermissionContext
}

func NewToIntersectionPermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToIntersectionPermContext {
	var p = new(ToIntersectionPermContext)

	InitEmptyPermissionContext(&p.PermissionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PermissionContext))

	return p
}

func (s *ToIntersectionPermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToIntersectionPermContext) IntersectionPerm() IIntersectionPermContext {
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

func (s *ToIntersectionPermContext) EOF() antlr.TerminalNode {
	return s.GetToken(AzmParserEOF, 0)
}

func (s *ToIntersectionPermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToIntersectionPerm(s)
	}
}

func (s *ToIntersectionPermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToIntersectionPerm(s)
	}
}

func (s *ToIntersectionPermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToIntersectionPerm(s)

	default:
		return t.VisitChildren(s)
	}
}

type ToUnionPermContext struct {
	PermissionContext
}

func NewToUnionPermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToUnionPermContext {
	var p = new(ToUnionPermContext)

	InitEmptyPermissionContext(&p.PermissionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PermissionContext))

	return p
}

func (s *ToUnionPermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToUnionPermContext) UnionPerm() IUnionPermContext {
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

func (s *ToUnionPermContext) EOF() antlr.TerminalNode {
	return s.GetToken(AzmParserEOF, 0)
}

func (s *ToUnionPermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToUnionPerm(s)
	}
}

func (s *ToUnionPermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToUnionPerm(s)
	}
}

func (s *ToUnionPermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToUnionPerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Permission() (localctx IPermissionContext) {
	localctx = NewPermissionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, AzmParserRULE_permission)
	p.SetState(36)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		localctx = NewToUnionPermContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(27)
			p.UnionPerm()
		}
		{
			p.SetState(28)
			p.Match(AzmParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewToIntersectionPermContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(30)
			p.IntersectionPerm()
		}
		{
			p.SetState(31)
			p.Match(AzmParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewToExclusionPermContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(33)
			p.ExclusionPerm()
		}
		{
			p.SetState(34)
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

// IUnionRelContext is an interface to support dynamic dispatch.
type IUnionRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRel() []IRelContext
	Rel(i int) IRelContext

	// IsUnionRelContext differentiates from other interfaces.
	IsUnionRelContext()
}

type UnionRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnionRelContext() *UnionRelContext {
	var p = new(UnionRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_unionRel
	return p
}

func InitEmptyUnionRelContext(p *UnionRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_unionRel
}

func (*UnionRelContext) IsUnionRelContext() {}

func NewUnionRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionRelContext {
	var p = new(UnionRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_unionRel

	return p
}

func (s *UnionRelContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionRelContext) AllRel() []IRelContext {
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

func (s *UnionRelContext) Rel(i int) IRelContext {
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

func (s *UnionRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterUnionRel(s)
	}
}

func (s *UnionRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitUnionRel(s)
	}
}

func (s *UnionRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitUnionRel(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) UnionRel() (localctx IUnionRelContext) {
	localctx = NewUnionRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, AzmParserRULE_unionRel)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(38)
		p.Rel()
	}
	p.SetState(43)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__0 {
		{
			p.SetState(39)
			p.Match(AzmParserT__0)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(40)
			p.Rel()
		}

		p.SetState(45)
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
	p.EnterRule(localctx, 6, AzmParserRULE_unionPerm)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(46)
		p.Perm()
	}
	p.SetState(51)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__0 {
		{
			p.SetState(47)
			p.Match(AzmParserT__0)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(48)
			p.Perm()
		}

		p.SetState(53)
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
	p.EnterRule(localctx, 8, AzmParserRULE_intersectionPerm)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(54)
		p.Perm()
	}
	p.SetState(59)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__1 {
		{
			p.SetState(55)
			p.Match(AzmParserT__1)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(56)
			p.Perm()
		}

		p.SetState(61)
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
	p.EnterRule(localctx, 10, AzmParserRULE_exclusionPerm)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(62)
		p.Perm()
	}
	{
		p.SetState(63)
		p.Match(AzmParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(64)
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

type ToWildcardRelContext struct {
	RelContext
}

func NewToWildcardRelContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToWildcardRelContext {
	var p = new(ToWildcardRelContext)

	InitEmptyRelContext(&p.RelContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelContext))

	return p
}

func (s *ToWildcardRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToWildcardRelContext) WildcardRel() IWildcardRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWildcardRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWildcardRelContext)
}

func (s *ToWildcardRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToWildcardRel(s)
	}
}

func (s *ToWildcardRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToWildcardRel(s)
	}
}

func (s *ToWildcardRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToWildcardRel(s)

	default:
		return t.VisitChildren(s)
	}
}

type ToSingleRelContext struct {
	RelContext
}

func NewToSingleRelContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToSingleRelContext {
	var p = new(ToSingleRelContext)

	InitEmptyRelContext(&p.RelContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelContext))

	return p
}

func (s *ToSingleRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToSingleRelContext) SingleRel() ISingleRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleRelContext)
}

func (s *ToSingleRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToSingleRel(s)
	}
}

func (s *ToSingleRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToSingleRel(s)
	}
}

func (s *ToSingleRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToSingleRel(s)

	default:
		return t.VisitChildren(s)
	}
}

type ToSubjectRelContext struct {
	RelContext
}

func NewToSubjectRelContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToSubjectRelContext {
	var p = new(ToSubjectRelContext)

	InitEmptyRelContext(&p.RelContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelContext))

	return p
}

func (s *ToSubjectRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToSubjectRelContext) SubjectRel() ISubjectRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubjectRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISubjectRelContext)
}

func (s *ToSubjectRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToSubjectRel(s)
	}
}

func (s *ToSubjectRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToSubjectRel(s)
	}
}

func (s *ToSubjectRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToSubjectRel(s)

	default:
		return t.VisitChildren(s)
	}
}

type ToArrowRelContext struct {
	RelContext
}

func NewToArrowRelContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToArrowRelContext {
	var p = new(ToArrowRelContext)

	InitEmptyRelContext(&p.RelContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelContext))

	return p
}

func (s *ToArrowRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToArrowRelContext) ArrowRel() IArrowRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrowRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrowRelContext)
}

func (s *ToArrowRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToArrowRel(s)
	}
}

func (s *ToArrowRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToArrowRel(s)
	}
}

func (s *ToArrowRelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToArrowRel(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Rel() (localctx IRelContext) {
	localctx = NewRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, AzmParserRULE_rel)
	p.SetState(70)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		localctx = NewToSingleRelContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(66)
			p.SingleRel()
		}

	case 2:
		localctx = NewToWildcardRelContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(67)
			p.WildcardRel()
		}

	case 3:
		localctx = NewToSubjectRelContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(68)
			p.SubjectRel()
		}

	case 4:
		localctx = NewToArrowRelContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(69)
			p.ArrowRel()
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

type ToArrowPermContext struct {
	PermContext
}

func NewToArrowPermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToArrowPermContext {
	var p = new(ToArrowPermContext)

	InitEmptyPermContext(&p.PermContext)
	p.parser = parser
	p.CopyAll(ctx.(*PermContext))

	return p
}

func (s *ToArrowPermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToArrowPermContext) ArrowRel() IArrowRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrowRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrowRelContext)
}

func (s *ToArrowPermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToArrowPerm(s)
	}
}

func (s *ToArrowPermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToArrowPerm(s)
	}
}

func (s *ToArrowPermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToArrowPerm(s)

	default:
		return t.VisitChildren(s)
	}
}

type ToSinglePermContext struct {
	PermContext
}

func NewToSinglePermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ToSinglePermContext {
	var p = new(ToSinglePermContext)

	InitEmptyPermContext(&p.PermContext)
	p.parser = parser
	p.CopyAll(ctx.(*PermContext))

	return p
}

func (s *ToSinglePermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToSinglePermContext) SingleRel() ISingleRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleRelContext)
}

func (s *ToSinglePermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterToSinglePerm(s)
	}
}

func (s *ToSinglePermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitToSinglePerm(s)
	}
}

func (s *ToSinglePermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case AzmVisitor:
		return t.VisitToSinglePerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *AzmParser) Perm() (localctx IPermContext) {
	localctx = NewPermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, AzmParserRULE_perm)
	p.SetState(74)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext()) {
	case 1:
		localctx = NewToSinglePermContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(72)
			p.SingleRel()
		}

	case 2:
		localctx = NewToArrowPermContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(73)
			p.ArrowRel()
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

// ISingleRelContext is an interface to support dynamic dispatch.
type ISingleRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsSingleRelContext differentiates from other interfaces.
	IsSingleRelContext()
}

type SingleRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySingleRelContext() *SingleRelContext {
	var p = new(SingleRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_singleRel
	return p
}

func InitEmptySingleRelContext(p *SingleRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_singleRel
}

func (*SingleRelContext) IsSingleRelContext() {}

func NewSingleRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SingleRelContext {
	var p = new(SingleRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_singleRel

	return p
}

func (s *SingleRelContext) GetParser() antlr.Parser { return s.parser }

func (s *SingleRelContext) ID() antlr.TerminalNode {
	return s.GetToken(AzmParserID, 0)
}

func (s *SingleRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SingleRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
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

func (p *AzmParser) SingleRel() (localctx ISingleRelContext) {
	localctx = NewSingleRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, AzmParserRULE_singleRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(76)
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

// ISubjectRelContext is an interface to support dynamic dispatch.
type ISubjectRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	HASH() antlr.TerminalNode

	// IsSubjectRelContext differentiates from other interfaces.
	IsSubjectRelContext()
}

type SubjectRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubjectRelContext() *SubjectRelContext {
	var p = new(SubjectRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_subjectRel
	return p
}

func InitEmptySubjectRelContext(p *SubjectRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_subjectRel
}

func (*SubjectRelContext) IsSubjectRelContext() {}

func NewSubjectRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubjectRelContext {
	var p = new(SubjectRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_subjectRel

	return p
}

func (s *SubjectRelContext) GetParser() antlr.Parser { return s.parser }

func (s *SubjectRelContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(AzmParserID)
}

func (s *SubjectRelContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(AzmParserID, i)
}

func (s *SubjectRelContext) HASH() antlr.TerminalNode {
	return s.GetToken(AzmParserHASH, 0)
}

func (s *SubjectRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubjectRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
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

func (p *AzmParser) SubjectRel() (localctx ISubjectRelContext) {
	localctx = NewSubjectRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, AzmParserRULE_subjectRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(78)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(79)
		p.Match(AzmParserHASH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(80)
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

// IWildcardRelContext is an interface to support dynamic dispatch.
type IWildcardRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	ASTERISK() antlr.TerminalNode

	// IsWildcardRelContext differentiates from other interfaces.
	IsWildcardRelContext()
}

type WildcardRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWildcardRelContext() *WildcardRelContext {
	var p = new(WildcardRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_wildcardRel
	return p
}

func InitEmptyWildcardRelContext(p *WildcardRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_wildcardRel
}

func (*WildcardRelContext) IsWildcardRelContext() {}

func NewWildcardRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WildcardRelContext {
	var p = new(WildcardRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_wildcardRel

	return p
}

func (s *WildcardRelContext) GetParser() antlr.Parser { return s.parser }

func (s *WildcardRelContext) ID() antlr.TerminalNode {
	return s.GetToken(AzmParserID, 0)
}

func (s *WildcardRelContext) COLON() antlr.TerminalNode {
	return s.GetToken(AzmParserCOLON, 0)
}

func (s *WildcardRelContext) ASTERISK() antlr.TerminalNode {
	return s.GetToken(AzmParserASTERISK, 0)
}

func (s *WildcardRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WildcardRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
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

func (p *AzmParser) WildcardRel() (localctx IWildcardRelContext) {
	localctx = NewWildcardRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, AzmParserRULE_wildcardRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(82)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(83)
		p.Match(AzmParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(84)
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

// IArrowRelContext is an interface to support dynamic dispatch.
type IArrowRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	ARROW() antlr.TerminalNode

	// IsArrowRelContext differentiates from other interfaces.
	IsArrowRelContext()
}

type ArrowRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrowRelContext() *ArrowRelContext {
	var p = new(ArrowRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_arrowRel
	return p
}

func InitEmptyArrowRelContext(p *ArrowRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_arrowRel
}

func (*ArrowRelContext) IsArrowRelContext() {}

func NewArrowRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrowRelContext {
	var p = new(ArrowRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_arrowRel

	return p
}

func (s *ArrowRelContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrowRelContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(AzmParserID)
}

func (s *ArrowRelContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(AzmParserID, i)
}

func (s *ArrowRelContext) ARROW() antlr.TerminalNode {
	return s.GetToken(AzmParserARROW, 0)
}

func (s *ArrowRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrowRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
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

func (p *AzmParser) ArrowRel() (localctx IArrowRelContext) {
	localctx = NewArrowRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, AzmParserRULE_arrowRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(86)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(87)
		p.Match(AzmParserARROW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(88)
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

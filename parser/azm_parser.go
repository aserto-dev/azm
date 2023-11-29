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
		"", "", "", "", "ARROW", "HASH", "COLON", "ASTERISK", "ID", "NEWLINE",
		"WS",
	}
	staticData.RuleNames = []string{
		"prog", "stat", "unionRel", "intersectRel", "exclusionRel", "rel", "singleRel",
		"subjectRel", "wildcardRel", "arrowRel",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 10, 78, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1, 0, 4,
		0, 22, 8, 0, 11, 0, 12, 0, 23, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 3, 1, 36, 8, 1, 1, 2, 1, 2, 1, 2, 5, 2, 41, 8, 2, 10,
		2, 12, 2, 44, 9, 2, 1, 3, 1, 3, 1, 3, 5, 3, 49, 8, 3, 10, 3, 12, 3, 52,
		9, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 62, 8, 5, 1,
		6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 0, 0, 10, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 0, 0, 76, 0,
		21, 1, 0, 0, 0, 2, 35, 1, 0, 0, 0, 4, 37, 1, 0, 0, 0, 6, 45, 1, 0, 0, 0,
		8, 53, 1, 0, 0, 0, 10, 61, 1, 0, 0, 0, 12, 63, 1, 0, 0, 0, 14, 65, 1, 0,
		0, 0, 16, 69, 1, 0, 0, 0, 18, 73, 1, 0, 0, 0, 20, 22, 3, 2, 1, 0, 21, 20,
		1, 0, 0, 0, 22, 23, 1, 0, 0, 0, 23, 21, 1, 0, 0, 0, 23, 24, 1, 0, 0, 0,
		24, 1, 1, 0, 0, 0, 25, 26, 3, 4, 2, 0, 26, 27, 5, 9, 0, 0, 27, 36, 1, 0,
		0, 0, 28, 29, 3, 6, 3, 0, 29, 30, 5, 9, 0, 0, 30, 36, 1, 0, 0, 0, 31, 32,
		3, 8, 4, 0, 32, 33, 5, 9, 0, 0, 33, 36, 1, 0, 0, 0, 34, 36, 5, 9, 0, 0,
		35, 25, 1, 0, 0, 0, 35, 28, 1, 0, 0, 0, 35, 31, 1, 0, 0, 0, 35, 34, 1,
		0, 0, 0, 36, 3, 1, 0, 0, 0, 37, 42, 3, 10, 5, 0, 38, 39, 5, 1, 0, 0, 39,
		41, 3, 10, 5, 0, 40, 38, 1, 0, 0, 0, 41, 44, 1, 0, 0, 0, 42, 40, 1, 0,
		0, 0, 42, 43, 1, 0, 0, 0, 43, 5, 1, 0, 0, 0, 44, 42, 1, 0, 0, 0, 45, 50,
		3, 10, 5, 0, 46, 47, 5, 2, 0, 0, 47, 49, 3, 10, 5, 0, 48, 46, 1, 0, 0,
		0, 49, 52, 1, 0, 0, 0, 50, 48, 1, 0, 0, 0, 50, 51, 1, 0, 0, 0, 51, 7, 1,
		0, 0, 0, 52, 50, 1, 0, 0, 0, 53, 54, 3, 10, 5, 0, 54, 55, 5, 3, 0, 0, 55,
		56, 3, 10, 5, 0, 56, 9, 1, 0, 0, 0, 57, 62, 3, 12, 6, 0, 58, 62, 3, 16,
		8, 0, 59, 62, 3, 14, 7, 0, 60, 62, 3, 18, 9, 0, 61, 57, 1, 0, 0, 0, 61,
		58, 1, 0, 0, 0, 61, 59, 1, 0, 0, 0, 61, 60, 1, 0, 0, 0, 62, 11, 1, 0, 0,
		0, 63, 64, 5, 8, 0, 0, 64, 13, 1, 0, 0, 0, 65, 66, 5, 8, 0, 0, 66, 67,
		5, 5, 0, 0, 67, 68, 5, 8, 0, 0, 68, 15, 1, 0, 0, 0, 69, 70, 5, 8, 0, 0,
		70, 71, 5, 6, 0, 0, 71, 72, 5, 7, 0, 0, 72, 17, 1, 0, 0, 0, 73, 74, 5,
		8, 0, 0, 74, 75, 5, 4, 0, 0, 75, 76, 5, 8, 0, 0, 76, 19, 1, 0, 0, 0, 5,
		23, 35, 42, 50, 61,
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
	AzmParserNEWLINE  = 9
	AzmParserWS       = 10
)

// AzmParser rules.
const (
	AzmParserRULE_prog         = 0
	AzmParserRULE_stat         = 1
	AzmParserRULE_unionRel     = 2
	AzmParserRULE_intersectRel = 3
	AzmParserRULE_exclusionRel = 4
	AzmParserRULE_rel          = 5
	AzmParserRULE_singleRel    = 6
	AzmParserRULE_subjectRel   = 7
	AzmParserRULE_wildcardRel  = 8
	AzmParserRULE_arrowRel     = 9
)

// IProgContext is an interface to support dynamic dispatch.
type IProgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllStat() []IStatContext
	Stat(i int) IStatContext

	// IsProgContext differentiates from other interfaces.
	IsProgContext()
}

type ProgContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgContext() *ProgContext {
	var p = new(ProgContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_prog
	return p
}

func InitEmptyProgContext(p *ProgContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_prog
}

func (*ProgContext) IsProgContext() {}

func NewProgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgContext {
	var p = new(ProgContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_prog

	return p
}

func (s *ProgContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgContext) AllStat() []IStatContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatContext); ok {
			len++
		}
	}

	tst := make([]IStatContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatContext); ok {
			tst[i] = t.(IStatContext)
			i++
		}
	}

	return tst
}

func (s *ProgContext) Stat(i int) IStatContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatContext); ok {
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

	return t.(IStatContext)
}

func (s *ProgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterProg(s)
	}
}

func (s *ProgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitProg(s)
	}
}

func (p *AzmParser) Prog() (localctx IProgContext) {
	localctx = NewProgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, AzmParserRULE_prog)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(21)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == AzmParserID || _la == AzmParserNEWLINE {
		{
			p.SetState(20)
			p.Stat()
		}

		p.SetState(23)
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

// IStatContext is an interface to support dynamic dispatch.
type IStatContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UnionRel() IUnionRelContext
	NEWLINE() antlr.TerminalNode
	IntersectRel() IIntersectRelContext
	ExclusionRel() IExclusionRelContext

	// IsStatContext differentiates from other interfaces.
	IsStatContext()
}

type StatContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatContext() *StatContext {
	var p = new(StatContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_stat
	return p
}

func InitEmptyStatContext(p *StatContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_stat
}

func (*StatContext) IsStatContext() {}

func NewStatContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatContext {
	var p = new(StatContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_stat

	return p
}

func (s *StatContext) GetParser() antlr.Parser { return s.parser }

func (s *StatContext) UnionRel() IUnionRelContext {
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

func (s *StatContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(AzmParserNEWLINE, 0)
}

func (s *StatContext) IntersectRel() IIntersectRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntersectRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntersectRelContext)
}

func (s *StatContext) ExclusionRel() IExclusionRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExclusionRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExclusionRelContext)
}

func (s *StatContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterStat(s)
	}
}

func (s *StatContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitStat(s)
	}
}

func (p *AzmParser) Stat() (localctx IStatContext) {
	localctx = NewStatContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, AzmParserRULE_stat)
	p.SetState(35)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(25)
			p.UnionRel()
		}
		{
			p.SetState(26)
			p.Match(AzmParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(28)
			p.IntersectRel()
		}
		{
			p.SetState(29)
			p.Match(AzmParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(31)
			p.ExclusionRel()
		}
		{
			p.SetState(32)
			p.Match(AzmParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(34)
			p.Match(AzmParserNEWLINE)
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

func (p *AzmParser) UnionRel() (localctx IUnionRelContext) {
	localctx = NewUnionRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, AzmParserRULE_unionRel)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(37)
		p.Rel()
	}
	p.SetState(42)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__0 {
		{
			p.SetState(38)
			p.Match(AzmParserT__0)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(39)
			p.Rel()
		}

		p.SetState(44)
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

// IIntersectRelContext is an interface to support dynamic dispatch.
type IIntersectRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRel() []IRelContext
	Rel(i int) IRelContext

	// IsIntersectRelContext differentiates from other interfaces.
	IsIntersectRelContext()
}

type IntersectRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntersectRelContext() *IntersectRelContext {
	var p = new(IntersectRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_intersectRel
	return p
}

func InitEmptyIntersectRelContext(p *IntersectRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_intersectRel
}

func (*IntersectRelContext) IsIntersectRelContext() {}

func NewIntersectRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntersectRelContext {
	var p = new(IntersectRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_intersectRel

	return p
}

func (s *IntersectRelContext) GetParser() antlr.Parser { return s.parser }

func (s *IntersectRelContext) AllRel() []IRelContext {
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

func (s *IntersectRelContext) Rel(i int) IRelContext {
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

func (s *IntersectRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntersectRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntersectRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterIntersectRel(s)
	}
}

func (s *IntersectRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitIntersectRel(s)
	}
}

func (p *AzmParser) IntersectRel() (localctx IIntersectRelContext) {
	localctx = NewIntersectRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, AzmParserRULE_intersectRel)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(45)
		p.Rel()
	}
	p.SetState(50)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__1 {
		{
			p.SetState(46)
			p.Match(AzmParserT__1)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(47)
			p.Rel()
		}

		p.SetState(52)
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

// IExclusionRelContext is an interface to support dynamic dispatch.
type IExclusionRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRel() []IRelContext
	Rel(i int) IRelContext

	// IsExclusionRelContext differentiates from other interfaces.
	IsExclusionRelContext()
}

type ExclusionRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExclusionRelContext() *ExclusionRelContext {
	var p = new(ExclusionRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_exclusionRel
	return p
}

func InitEmptyExclusionRelContext(p *ExclusionRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_exclusionRel
}

func (*ExclusionRelContext) IsExclusionRelContext() {}

func NewExclusionRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExclusionRelContext {
	var p = new(ExclusionRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_exclusionRel

	return p
}

func (s *ExclusionRelContext) GetParser() antlr.Parser { return s.parser }

func (s *ExclusionRelContext) AllRel() []IRelContext {
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

func (s *ExclusionRelContext) Rel(i int) IRelContext {
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

func (s *ExclusionRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExclusionRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExclusionRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterExclusionRel(s)
	}
}

func (s *ExclusionRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitExclusionRel(s)
	}
}

func (p *AzmParser) ExclusionRel() (localctx IExclusionRelContext) {
	localctx = NewExclusionRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, AzmParserRULE_exclusionRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(53)
		p.Rel()
	}
	{
		p.SetState(54)
		p.Match(AzmParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(55)
		p.Rel()
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

	// Getter signatures
	SingleRel() ISingleRelContext
	WildcardRel() IWildcardRelContext
	SubjectRel() ISubjectRelContext
	ArrowRel() IArrowRelContext

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

func (s *RelContext) SingleRel() ISingleRelContext {
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

func (s *RelContext) WildcardRel() IWildcardRelContext {
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

func (s *RelContext) SubjectRel() ISubjectRelContext {
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

func (s *RelContext) ArrowRel() IArrowRelContext {
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

func (s *RelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterRel(s)
	}
}

func (s *RelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitRel(s)
	}
}

func (p *AzmParser) Rel() (localctx IRelContext) {
	localctx = NewRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, AzmParserRULE_rel)
	p.SetState(61)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(57)
			p.SingleRel()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(58)
			p.WildcardRel()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(59)
			p.SubjectRel()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(60)
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

func (p *AzmParser) SingleRel() (localctx ISingleRelContext) {
	localctx = NewSingleRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, AzmParserRULE_singleRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(63)
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

func (p *AzmParser) SubjectRel() (localctx ISubjectRelContext) {
	localctx = NewSubjectRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, AzmParserRULE_subjectRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(65)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(66)
		p.Match(AzmParserHASH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(67)
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

func (p *AzmParser) WildcardRel() (localctx IWildcardRelContext) {
	localctx = NewWildcardRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, AzmParserRULE_wildcardRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(70)
		p.Match(AzmParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(71)
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

func (p *AzmParser) ArrowRel() (localctx IArrowRelContext) {
	localctx = NewArrowRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, AzmParserRULE_arrowRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(73)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(74)
		p.Match(AzmParserARROW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(75)
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

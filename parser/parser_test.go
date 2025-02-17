package parser

import (
	"kaze/ast"
	"kaze/lexer"
	"strconv"
	"testing"
)

func testInfixExpression(t *testing.T, exp *ast.InfixExpression, left interface{}, operator string, right interface{}) bool {
	if !testLiteralExpression(t, exp.Left, left) {
		return false
	}

	if exp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, exp.Operator)
		return false
	}

	if !testLiteralExpression(t, exp.Right, right) {
		return false
	}

	return true
}

func testVarStatement(t *testing.T, stmt ast.Statement, s string) bool {
	if stmt.TokenLiteral() != "var" {
		t.Errorf("stmt.TokenLiteral not 'var'. got=%q", stmt.TokenLiteral())
		return false
	}

	varStmt, ok := stmt.(*ast.VarStatement)
	if !ok {
		t.Errorf("stmt not *ast.VarStatement. got=%T", stmt)
		return false
	}

	if varStmt.Name.Value != s {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", s, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != s {
		t.Errorf("varStmt.Name.TokenLiteral not '%s'. got=%s", s, varStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		switch exp.(type) {
		case *ast.StringLiteral:
			return testStringLiteral(t, exp, v)
		case *ast.Identifier:
			return testIdentifier(t, exp, v)
		}
		t.Fatalf("type of exp not handled. got=%T", exp)
	case bool:
		return testBoolean(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testStringLiteral(t *testing.T, exp ast.Expression, v string) bool {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp not *ast.StringLiteral. got=%T", exp)
		return false
	}

	if str.Value != v {
		t.Errorf("str.Value not %s. got=%s", v, str.Value)
		return false
	}

	if str.TokenLiteral() != v {
		t.Errorf("str.TokenLiteral not %s. got=%s", v, str.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != strconv.FormatInt(value, 10) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if b.Value != value {
		t.Errorf("bool.Value not %t. got=%t", value, b.Value)
		return false
	}

	if b.TokenLiteral() != strconv.FormatBool(value) {
		t.Errorf("bool.TokenLiteral not %t. got=%s", value, b.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"var x = 5;", "x", 5},
		{"var y = true;", "y", true},
		{"var foobar = y;", "foobar", "y"},
		{`var hoge = "hoge";`, "hoge", "hoge"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.VarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
		{"false != true;", false, "!=", true},
		{`"hoge" + "fuga";`, "hoge", "+", "fuga"},
		{`"hoge" == "hoge";`, "hoge", "==", "hoge"},
		{`"hoge" != "fuga";`, "hoge", "!=", "fuga"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestAssignExpression(t *testing.T) {
	tests := []struct {
		input string
		name  interface{}
		value interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y;", "foobar", "y"},
		{`hoge = "hoge";`, "hoge", "hoge"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.AssignExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.AssignExpression. got=%T", stmt.Expression)
		}

		if exp.Name.Value != tt.name {
			t.Fatalf("exp.Name.Value not '%s'. got=%s", tt.name, exp.Name.Value)
		}

		if !testLiteralExpression(t, exp.Value, tt.value) {
			return
		}
	}
}

func TestBlockExpression(t *testing.T) {
	input := `
{
	var x = 5;
	var y = 10;
	var foobar = 838383;
}
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.BlockExpression. got=%T", program.Statements[0])
	}

	block, ok := exprStmt.Expression.(*ast.BlockExpression)
	if !ok {
		t.Fatalf("block is not *ast.BlockExpression. got=%T", exprStmt.Expression)
	}

	tests := []string{"x", "y", "foobar"}

	for i, tt := range tests {
		stmt := block.Statements[i]
		if !testVarStatement(t, stmt, tt) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `return 5; return 10;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}

}

func TestFunctionDefinitionStatement(t *testing.T) {
	input := `
	fun add(x, y) {
		return x + y;
	}
	fun empty() {}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedName       string
		expectedParameters []string
		expectedBody       string
	}{
		{"add", []string{"x", "y"}, "return (x + y);"},
		{"empty", []string{}, ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		fds, ok := stmt.(*ast.FunctionDefinitionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.FunctionDefinitionStatement. got=%T", stmt)
		}
		if fds.Name.Value != tt.expectedName {
			t.Fatalf("fds.Name.Value not %s. got=%s", tt.expectedName, fds.Name.Value)
		}
		if len(fds.Parameters) != len(tt.expectedParameters) {
			t.Fatalf("len(fds.Parameters) not %d. got=%d", len(tt.expectedParameters), len(fds.Parameters))
		}
		for j, ident := range fds.Parameters {
			if ident.Value != tt.expectedParameters[j] {
				t.Fatalf("ident.Value not %s. got=%s", tt.expectedParameters[j], ident.Value)
			}
		}
		if fds.Body.String() != tt.expectedBody {
			t.Fatalf("fds.Body.String() not %s. got=%s", tt.expectedBody, fds.Body.String())
		}
	}
}

func TestCallExpression(t *testing.T) {
	tests := []struct {
		input    string
		function string
		args     []interface{}
	}{
		{"add(1, 2);", "add", []interface{}{"1", "2"}},
		{"add(1 + 2, 3 * 4);", "add", []interface{}{"(1 + 2)", "(3 * 4)"}},
		{"zero()", "zero", []interface{}{}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		call, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T", stmt.Expression)
		}
		if call.Function.String() != tt.function {
			t.Fatalf("call.Function.Value not %s. got=%s", tt.function, call.Function.String())
		}
		if len(call.Arguments) != len(tt.args) {
			t.Fatalf("len(call.Arguments) not %d. got=%d", len(tt.args), len(call.Arguments))
		}
		for i, arg := range tt.args {
			if call.Arguments[i].String() != arg {
				t.Fatalf("call.Arguments[%d] not %s. got=%s", i, arg, call.Arguments[i].String())
			}
		}
	}
}

func TestIfExpressions(t *testing.T) {
	tests := []struct {
		input       string
		condition   string
		consequence string
		alternative string
	}{
		{"if x < y { x }", "(x < y)", "x", ""},
		{"if x < y { x } else { y }", "(x < y)", "x", "y"},
		{"if (x < y) { x }", "(x < y)", "x", ""},
		{"if (x < y) { x } else { y }", "(x < y)", "x", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		ifExp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
		}
		if ifExp.Condition.String() != tt.condition {
			t.Fatalf("ifExp.Condition.String() not %s. got=%s", tt.condition, ifExp.Condition.String())
		}
		if ifExp.Consequence.String() != tt.consequence {
			t.Fatalf("ifExp.Consequence.String() not %s. got=%s", tt.consequence, ifExp.Consequence.String())
		}
		if ifExp.Alternative == nil {
			if tt.alternative != "" {
				t.Fatalf("ifExp.Alternative is nil. got=%s", tt.alternative)
			}
		} else {
			if ifExp.Alternative.String() != tt.alternative {
				t.Fatalf("ifExp.Alternative.String() not %s. got=%s", tt.alternative, ifExp.Alternative.String())
			}
		}
	}
}

func TestWhileStatement(t *testing.T) {
	tests := []struct {
		input     string
		condition string
		body      string
	}{
		{"while x < y { x }", "(x < y)", "x"},
		{"while x < y { x + y }", "(x < y)", "(x + y)"},
		{"while true { break; }", "true", "break"},
		{"while true { continue; }", "true", "continue"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.WhileStatement)
		if !ok {
			t.Fatalf("stmt not *ast.WhileStatement. got=%T", program.Statements[0])
		}
		if stmt.Condition.String() != tt.condition {
			t.Fatalf("stmt.Condition.String() not %s. got=%s", tt.condition, stmt.Condition.String())
		}
		if stmt.Body.String() != tt.body {
			t.Fatalf("stmt.Body.String() not %s. got=%s", tt.body, stmt.Body.String())
		}
	}
}

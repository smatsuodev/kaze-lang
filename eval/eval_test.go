package eval

import (
	"kaze/lexer"
	"kaze/object"
	"kaze/parser"
	"testing"
)

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) {
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", evaluated, evaluated)
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) {
	result, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Fatalf("object is not Boolean. got=%T (%+v)", evaluated, evaluated)
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}
}

func testStringObject(t *testing.T, evaluated object.Object, expected string) {
	result, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if result.Value != expected {
		t.Fatalf("object has wrong value. got=%s, want=%s", result.Value, expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// Boolean
func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"a" < "b"`, true},
		{`"b" < "a"`, false},
		{`"b" > "a"`, true},
		{`"a" > "b"`, false},
		{`"hoge" == "hoge"`, true},
		{`"hoge" == "fuga"`, false},
		{`"hoge" != "fuga"`, true},
		{`"hoge" != "hoge"`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hoge";`, "hoge"},
		{`"hoge" + "hoge";`, "hogehoge"},
		{`"hoge"[0]`, "h"},
		{`"hoge"[0] + "fuga"[0]`, "hf"},
		{`var a = "hoge"; var b = 0; a[b]`, "h"},
		{`fun greet(name) { "Hello, " + name + "!"; }; greet("Alice");`, "Hello, Alice!"},
		{`fun greet(name) { "Hello, " + name + "!"; }; fun add(x,y){x+y}; greet("Alice")[add(3,4)];`, "A"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

// BANG operator
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"var a = 5; a = 10; a;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBlockExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"{ var a = 5; a; }", 5},
		{"var a = 1; { a = 5; }; a", 5},
		{"var a = 1; { var a = 5; }; a", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"var a = 5; return a;", 5},
		{"var a = 5; return a; return 10;", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"fun identity(x) { x; } identity(5);", 5},
		{"fun identity(x) { return x; } identity(5);", 5},
		{"fun sum(x, y) { x + y; } sum(5, 10);", 15},
		{"fun sum(x, y) { x + y; } sum(5 + 5, 10 + 10);", 30},
		{"fun sum(x, y) { var a = x + y; a; } sum(5, 10);", 15},
		{"fun fact(x) { if (x == 0) { return 1; } else { return x * fact(x - 1); } } fact(5);", 120},
		{"var x = 10; fun f(x) { return x; } f(5);", 5},
		{"var x = 10; fun f(x) { return x; } f(x);", 10},
		{"var x = 10; fun f(x) { x = 5; } x;", 10},
		{"var x = 10; fun f() { return x; } f();", 10},
		{"var x = 10; fun f() { x = 5; } x;", 10},
		{"var x = 10; fun f() { x = 5; return x; } f(); x;", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestIfExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10; }", 10},
		{"if false { 10; }", NULL},
		{"if 1 { 10; }", 10},
		{"if 1 < 2 { 10; }", 10},
		{"if 1 > 2 { 10; }", NULL},
		{"if 1 > 2 { 10; } else { 20; }", 20},
		{"if 1 < 2 { 10; } else { 20; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case int64:
			testIntegerObject(t, evaluated, v)
		case bool:
			testBooleanObject(t, evaluated, v)
		case object.Object:
			if evaluated != v {
				t.Fatalf("object has wrong value. got=%+v, want=%+v", evaluated, v)
			}
		}
	}
}

func TestWhileStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var x = 0; while x < 10 { x = x + 1; } x;", 10},
		{"var x = 0; while x < 10 { x = x + 1; if x == 5 { break; } } x;", 5},
		{"var x = 0; var y = 0; while x < 10 { x = x + 1; if x > 5 { continue; } y = y + 1; } y;", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

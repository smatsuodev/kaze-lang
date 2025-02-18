package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Fatalf("strings with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Fatalf("strings with same content have different hash keys")
	}
	if hello1.HashKey() == diff1.HashKey() {
		t.Fatalf("strings with different content have same hash keys")
	}
}

func TestLValueGet(t *testing.T) {
	env := NewEnvironment()

	env.Create("int", &Integer{Value: 1})
	intVar := &Variable{Name: "int", Env: env}
	obj, ok := intVar.Get()
	if !ok {
		t.Fatalf("intVar.Get() returned false")
	}
	if obj.Type() != INTEGER_OBJ {
		t.Fatalf("intVar.Get() returned wrong type")
	}
	if obj.(*Integer).Value != 1 {
		t.Fatalf("intVar.Get() returned wrong value")
	}

	env.Create("array", &Array{Elements: []Object{&Integer{Value: 1}}})
	arrayVar := &Variable{Name: "array", Env: env}
	arrayRef := &IndexRef{Left: arrayVar, Index: &Integer{Value: 0}}
	obj, ok = arrayRef.Get()
	if !ok {
		t.Fatalf("arrayRef.Get() returned false")
	}
	if obj.Type() != INTEGER_OBJ {
		t.Fatalf("arrayRef.Get() returned wrong type")
	}
	if obj.(*Integer).Value != 1 {
		t.Fatalf("arrayRef.Get() returned wrong value")
	}

	env.Create("hash", &Hash{Pairs: map[HashKey]HashPair{(&String{Value: "key"}).HashKey(): HashPair{Key: &String{Value: "key"}, Value: &Integer{Value: 1}}}})
	hashVar := &Variable{Name: "hash", Env: env}
	hashRef := &IndexRef{Left: hashVar, Index: &String{Value: "key"}}
	obj, ok = hashRef.Get()
	if !ok {
		t.Fatalf("hashRef.Get() returned false")
	}
	if obj.Type() != INTEGER_OBJ {
		t.Fatalf("hashRef.Get() returned wrong type")
	}
	if obj.(*Integer).Value != 1 {
		t.Fatalf("hashRef.Get() returned wrong value")
	}

	env.Create("string", &String{Value: "Hello World"})
	stringVar := &Variable{Name: "string", Env: env}
	stringRef := &IndexRef{Left: stringVar, Index: &Integer{Value: 0}}
	obj, ok = stringRef.Get()
	if !ok {
		t.Fatalf("stringRef.Get() returned false")
	}
	if obj.Type() != STRING_OBJ {
		t.Fatalf("stringRef.Get() returned wrong type")
	}
	if obj.(*String).Value != "H" {
		t.Fatalf("stringRef.Get() returned wrong value")
	}
}

func TestMultiDimensionIndexRefGet(t *testing.T) {
	env := NewEnvironment()

	env.Create("array", &Array{Elements: []Object{&Array{Elements: []Object{&Integer{Value: 1}}}}})
	arrayVar := &Variable{Name: "array", Env: env}
	arrayRef := &IndexRef{Left: arrayVar, Index: &Integer{Value: 0}}
	arrayRef2 := &IndexRef{Left: arrayRef, Index: &Integer{Value: 0}}
	obj, ok := arrayRef2.Get()
	if !ok {
		t.Fatalf("arrayRef2.Get() returned false")
	}
	if obj.Type() != INTEGER_OBJ {
		t.Fatalf("arrayRef2.Get() returned wrong type")
	}
	if obj.(*Integer).Value != 1 {
		t.Fatalf("arrayRef2.Get() returned wrong value")
	}
}

func TestUpdateArrayIndexRef(t *testing.T) {
	env := NewEnvironment()

	env.Create("array", &Array{Elements: []Object{&Integer{Value: 1}}})
	arrayVar := &Variable{Name: "array", Env: env}
	arrayRef := &IndexRef{Left: arrayVar, Index: &Integer{Value: 0}}
	newValue := &Integer{Value: 2}
	obj, ok := arrayRef.Update(newValue)
	if !ok {
		t.Fatalf("arrayRef.Update() returned false")
	}
	if obj.Type() != INTEGER_OBJ {
		t.Fatalf("arrayRef.Update() returned wrong type")
	}
	if obj.(*Integer).Value != 2 {
		t.Fatalf("arrayRef.Update() returned wrong value")
	}
	if e, _ := env.Get("array"); e.(*Array).Elements[0].(*Integer).Value != 2 {
		t.Fatalf("arrayRef.Update() did not update the array")
	}
}

func TestUpdateNestedArrayIndexRef(t *testing.T) {
	env := NewEnvironment()

	env.Create("array", &Array{Elements: []Object{&Array{Elements: []Object{&Integer{Value: 1}}}}})
	arrayVar := &Variable{Name: "array", Env: env}
	arrayRef := &IndexRef{Left: arrayVar, Index: &Integer{Value: 0}}
	arrayRef2 := &IndexRef{Left: arrayRef, Index: &Integer{Value: 0}}
	newValue := &Integer{Value: 2}
	obj, ok := arrayRef2.Update(newValue)
	if !ok {
		t.Fatalf("arrayRef2.Update() returned false")
	}
	if obj.Type() != INTEGER_OBJ {
		t.Fatalf("arrayRef2.Update() returned wrong type")
	}
	if obj.(*Integer).Value != 2 {
		t.Fatalf("arrayRef2.Update() returned wrong value")
	}
	if e, _ := env.Get("array"); e.(*Array).Elements[0].(*Array).Elements[0].(*Integer).Value != 2 {
		t.Fatalf("arrayRef2.Update() did not update the array")
	}
}

func TestUpdateHashIndexRef(t *testing.T) {
	env := NewEnvironment()

	env.Create("hash", &Hash{Pairs: map[HashKey]HashPair{(&String{Value: "key"}).HashKey(): HashPair{Key: &String{Value: "key"}, Value: &Integer{Value: 1}}}})
	hashVar := &Variable{Name: "hash", Env: env}
	hashRef := &IndexRef{Left: hashVar, Index: &String{Value: "key"}}
	newValue := &Integer{Value: 2}
	obj, ok := hashRef.Update(newValue)
	if !ok {
		t.Fatalf("hashRef.Update() returned false")
	}
	if obj.Type() != INTEGER_OBJ {
		t.Fatalf("hashRef.Update() returned wrong type")
	}
	if obj.(*Integer).Value != 2 {
		t.Fatalf("hashRef.Update() returned wrong value")
	}
	if e, _ := env.Get("hash"); e.(*Hash).Pairs[(&String{Value: "key"}).HashKey()].Value.(*Integer).Value != 2 {
		t.Fatalf("hashRef.Update() did not update the hash")
	}
}

func TestUpdateNestedHashIndexRef(t *testing.T) {
	env := NewEnvironment()

	env.Create("hash", &Hash{Pairs: map[HashKey]HashPair{(&String{Value: "key"}).HashKey(): HashPair{Key: &String{Value: "key"}, Value: &Hash{Pairs: map[HashKey]HashPair{(&String{Value: "key"}).HashKey(): HashPair{Key: &String{Value: "key"}, Value: &Integer{Value: 1}}}}}}})
	hashVar := &Variable{Name: "hash", Env: env}
	hashRef := &IndexRef{Left: hashVar, Index: &String{Value: "key"}}
	hashRef2 := &IndexRef{Left: hashRef, Index: &String{Value: "key"}}
	newValue := &Integer{Value: 2}
	obj, ok := hashRef2.Update(newValue)
	if !ok {
		t.Fatalf("hashRef2.Update() returned false")
	}
	if obj.Type() != INTEGER_OBJ {
		t.Fatalf("hashRef2.Update() returned wrong type")
	}
	if obj.(*Integer).Value != 2 {
		t.Fatalf("hashRef2.Update() returned wrong value")
	}
	if e, _ := env.Get("hash"); e.(*Hash).Pairs[(&String{Value: "key"}).HashKey()].Value.(*Hash).Pairs[(&String{Value: "key"}).HashKey()].Value.(*Integer).Value != 2 {
		t.Fatalf("hashRef2.Update() did not update the hash")
	}
}

func TestUpdateStringIndexRef(t *testing.T) {
	env := NewEnvironment()

	env.Create("string", &String{Value: "Hello World"})
	stringVar := &Variable{Name: "string", Env: env}
	stringRef := &IndexRef{Left: stringVar, Index: &Integer{Value: 0}}
	newValue := &String{Value: "h"}
	obj, ok := stringRef.Update(newValue)
	if !ok {
		t.Fatalf("stringRef.Update() returned false")
	}
	if obj.Type() != STRING_OBJ {
		t.Fatalf("stringRef.Update() returned wrong type")
	}
	if obj.(*String).Value != "h" {
		t.Fatalf("stringRef.Update() returned wrong value")
	}
	if e, _ := env.Get("string"); e.(*String).Value != "hello World" {
		t.Fatalf("stringRef.Update() did not update the string")
	}
}

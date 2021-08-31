package jsontype_test

import (
	"context"
	"fmt"

	"github.com/sinmetal/til/spanner/jsontype"
)

// 何も考えずにinterface{}を突っ込んでみて、無理ですと言われている図
func ExampleInsertJsonSample1() {
	ctx := context.Background()

	if err := jsontype.InsertJsonSample1(ctx); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// spanner: code = "InvalidArgument", desc = "client doesn't support type *jsontype.JsonBody"
}

// []byteにjsonを突っ込んでみて、やっぱりダメですと言われる図
func ExampleInsertJsonSample2() {
	ctx := context.Background()

	if err := jsontype.InsertJsonSample2(ctx); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// spanner: code = "FailedPrecondition", desc = "Invalid value for column Json in table JsonSample: Expected JSON."
}

// spanner.NullJSONを使った図
func ExampleInsertJsonSample3() {
	ctx := context.Background()

	if err := jsontype.InsertJsonSample3(ctx); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// jsontype.Entity{ID:"sample", Json:spanner.NullJSON{Value:map[string]interface {}{"Count":100, "Date":"2021-08-30T13:00:00Z", "ID":"helloJson"}, Valid:true}}
}

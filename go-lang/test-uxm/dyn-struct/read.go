
// GO Lang
// dynamic-structure / read dynamic struct and detect if some fields exist
// (c) 2020 unix-world.org

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/unix-world/smartgo/dynamic-struct"
)

type DataOne struct {
	Integer int     `json:"int"`
	Text    string  `json:"someText"`
	Float   float64 `json:"double"`
}

type DataTwo struct {
	Boolean bool
	Slice []int
	Anonymous string `json:"-"`
}

func main() {
	instance := dynamicstruct.MergeStructs(DataOne{}, DataTwo{}).
		Build().
		New()

	data := []byte(`
{
"int": 123,
"someText": "example",
"double": 123.45,
"Boolean": true,
"Slice": [1, 2, 3],
"Anonymous": "avoid to read"
}
`)

	err := json.Unmarshal(data, &instance)
	if err != nil {
		log.Fatal(err)
	}

	reader := dynamicstruct.NewReader(instance)

	fmt.Println("Integer", reader.GetField("Integer").Int())
	fmt.Println("Text", reader.GetField("Text").String())
	fmt.Println("Float", reader.GetField("Float").Float64())
	fmt.Println("Boolean", reader.GetField("Boolean").Bool())
	fmt.Println("Slice", reader.GetField("Slice").Interface().([]int))
	if(reader.HasField("Anonymous")) {
		fmt.Println("Anonymous", reader.GetField("Anonymous").String())
	} else {
		fmt.Println("ERROR: The field `Anonymous` ... does exists but not detected ...")
	}

	if(reader.HasField("NotExisting")) {
		fmt.Println("ERROR: The field `NotExisting` ... does not exists but was wrong detected ...", reader.GetField("NotExisting").String())
	} else {
		fmt.Println("OK: The field `NotExisting` ... does not exists ...")
	}

	var dataOne DataOne
	err = reader.ToStruct(&dataOne)
	fmt.Println(err, dataOne)

	var dataTwo DataTwo
	err = reader.ToStruct(&dataTwo)
	fmt.Println(err, dataTwo)
	// Out:
	// Integer 123
	// Text example
	// Float 123.45
	// Boolean true
	// Slice [1 2 3]
	// Anonymous
	// <nil> {123 example 123.45}
	// <nil> {true [1 2 3] }

}

// #END

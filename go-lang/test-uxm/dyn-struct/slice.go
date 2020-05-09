
// GO Lang
// dynamic-structure / make a slice of dynamic struct
// (c) 2020 unix-world.org

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/unix-world/smartgo/dynamic-struct"
)

type Data struct {
	Integer   int     `json:"int"`
	Text      string  `json:"someText"`
	Float     float64 `json:"double"`
	Boolean   bool
	Slice     []int
	Anonymous string `json:"-"`
}

func main() {

	definition := dynamicstruct.ExtendStruct(Data{}).Build()

	slice := definition.NewSliceOfStructs()

	data := []byte(`
[
	{
		"int": 123,
		"someText": "example",
		"double": 123.45,
		"Boolean": true,
		"Slice": [1, 2, 3],
		"Anonymous": "avoid to read"
	}
]
`)

	err := json.Unmarshal(data, &slice)
	if err != nil {
		log.Fatal(err)
	}

	data, err = json.Marshal(slice)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
	// Out:
	// [{"Boolean":true,"Slice":[1,2,3],"int":123,"someText":"example","double":123.45}]

	reader := dynamicstruct.NewReader(slice)
	readersSlice := reader.ToSliceOfReaders()
	for k, v := range readersSlice {
		var value Data
		err := v.ToStruct(&value)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(k, value)
	}

	// Out:
	// 0 {123 example 123.45 true [1 2 3] }

}

// #END

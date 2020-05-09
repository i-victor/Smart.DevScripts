
// GO Lang
// dynamic-structure / make a map of dynamic struct
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

	mapWithStringKey := definition.NewMapOfStructs("")

	data := []byte(`
{
	"element": {
		"int": 123,
		"someText": "example",
		"double": 123.45,
		"Boolean": true,
		"Slice": [1, 2, 3],
		"Anonymous": "avoid to read"
	}
}
`)

	err := json.Unmarshal(data, &mapWithStringKey)
	if err != nil {
		log.Fatal(err)
	}

	data, err = json.Marshal(mapWithStringKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
	// Out:
	// {"element":{"int":123,"someText":"example","double":123.45,"Boolean":true,"Slice":[1,2,3]}}

	reader := dynamicstruct.NewReader(mapWithStringKey)
	readersMap := reader.ToMapReaderOfReaders()
	for k, v := range readersMap {
		var value Data
		err := v.ToStruct(&value)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(k, value)
	}

	// Out:
	// element {123 example 123.45 true [1 2 3] }

}

// #END

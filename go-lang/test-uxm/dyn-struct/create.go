
// GO Lang
// dynamic-structure / create a new structure
// (c) 2020 unix-world.org

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/unix-world/smartgo/dynamic-struct"
)

func main() {
	instance := dynamicstruct.NewStruct().
		AddField("Integer", 0, `json:"int"`).
		AddField("Text", "", `json:"someText"`).
		AddField("Float", 0.0, `json:"double"`).
		AddField("Boolean", false, "").
		AddField("Slice", []int{}, "").
		AddField("Anonymous", "", `json:"-"`).
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

	data, err = json.Marshal(instance)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data)) // Out: {"int":123,"someText":"example","double":123.45,"Boolean":true,"Slice":[1,2,3]}

}

// #END

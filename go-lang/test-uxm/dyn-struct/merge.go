
// GO Lang
// dynamic-structure / merge existing structs
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

	data, err = json.Marshal(instance)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data)) // Out: {"int":123,"someText":"example","double":123.45,"Boolean":true,"Slice":[1,2,3]}

}

// #END

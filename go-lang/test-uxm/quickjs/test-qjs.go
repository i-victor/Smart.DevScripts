
// GO Lang
// (c) 2020 unix-world.org
// version: 2020.02.18

package main

import (
    "github.com/alexsunday/quickjs"
)

func main()  {
    rt := quickjs.NewJsRuntime()
    defer rt.Close()
    ctx := rt.NewContext()
    defer ctx.Close()
    res := ctx.EvalCode("console.log('1+1=', 1+1)")
    //res := ctx.EvalCode("var a = { b:2, c:3 }; return JSON.stringify(a);")
    println(res)
}

// #END

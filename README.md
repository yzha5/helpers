# helpers
初学阶段，自己写的常用的函数（部分网上拷的）

## 安装
```
go get github.com/yzha5/helpers
```

## 使用

### 判断

#### IsNumeric(val interface{}) bool
> 判断传入的数值或字符串是否为数字

```go
package main

import (

"fmt"
"github.com/yzha5/helpers/check"
)

func main() {
    var valueF = "56.25"//浮点
    var valueE = "6E3"//科学计数
    var valueN = "123"//数字
    var valueI = 123//int
    var valueH = "0xFFFFFF"//16进制
    var valueS = "abc"//一般的字符或字符串
    
    var b bool

    b = check.IsNumeric(valueF)
    fmt.Println("valueF:",b) //valueF: true

    b = check.IsNumeric(valueE)
    fmt.Println("valueF:",b) //valueE: true

    b = check.IsNumeric(valueN)
    fmt.Println("valueF:",b) //valueN: true

    b = check.IsNumeric(valueI)
    fmt.Println("valueF:",b) //valueI: true

    b = check.IsNumeric(valueH)
    fmt.Println("valueF:",b) //valueH: true

    b = check.IsNumeric(valueS)
    fmt.Println("valueF:",b) //valueS: false

}
```

#### IsDecimal(val interface{}) bool
> 判断传入的数值或字符串是否为`十进制`数字

```go
package main

import (

"fmt"
"github.com/yzha5/helpers/check"
)

func main() {
    var valueF = "56.25"//浮点
    var valueE = "6E3"//科学计数
    var valueN = "123"//数字
    var valueI = 123//int
    var valueH = "0xFFFFFF"//16进制
    var valueS = "abc"//一般的字符或字符串
    
    var b bool

    b = check.IsDecimal(valueF)
    fmt.Println("valueF:",b) //valueF: true

    b = check.IsDecimal(valueE)
    fmt.Println("valueF:",b) //valueE: false

    b = check.IsDecimal(valueN)
    fmt.Println("valueF:",b) //valueN: true

    b = check.IsDecimal(valueI)
    fmt.Println("valueF:",b) //valueI: true

    b = check.IsDecimal(valueH)
    fmt.Println("valueF:",b) //valueH: false

    b = check.IsDecimal(valueS)
    fmt.Println("valueF:",b) //valueS: false

}
```

### 数据转换

#### ArrayToString(strA []string, spl string, bra [2]string) (str string)
> 将字符串数组转换为一个字符串 可带分隔符
> 
> 假如：参数 `strA:["abc","def"] spl:"|" bra:["(",")"]` ==> `"(abc)|(def)"`
> 
> `str`字符串数组
> 
> `spl`分隔符
> 
> `bra`使用括号将每个数组元素括起来


#### StructToMap(obj interface{}) (map_ map[string]interface{}, err error)
> 将`struct`类型的数据转换为`map`，支持嵌套的数据
>
> obj可以是`指针`也可以是`值`
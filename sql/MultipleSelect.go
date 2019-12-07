package sql

import (
	"errors"
	"fmt"
	"github.com/yzha5/helpers/check"
	"github.com/yzha5/helpers/conv"
	"strings"
)

/**
 * 相对复杂的条件查询
 * 传入查询条件，返回MySql查询语句
 */

type Type string

type Oper string

type Column struct {
	Type   Type
	Col    string
	Oper   Oper
	Values []string
}

type Cond struct {
	And []*Column
	Or  []*Column
}

const (
	EQ       Oper = "EQ"       // =
	NEQ      Oper = "NEQ"      // <>
	GT       Oper = "GT"       // >
	GTE      Oper = "GET"      // >=
	LT       Oper = "LT"       // <
	LTE      Oper = "LTE"      // <=
	BETWEEN  Oper = "BETWEEN"  // BETWEEN xx AND xx
	NBETWEEN Oper = "NBETWEEN" // NOT BETWEEN xx AND xx
	LIKE     Oper = "LIKE"     // LIKE 'xx' 传入的值不需要带%，默认匹配 %value%
	NLIKE    Oper = "NLIKE"    // NOT LIKE 'xx' 传入的值不需要带%，默认匹配 %value%
	NULL     Oper = "NULL"     // IS NULL
	NNULL    Oper = "NNULL"    // IS NOT NULL
	IN       Oper = "IN"       // IN (xx1,xx2,xx3...)
	NIN      Oper = "NIN"      // NOT IN (xx1,xx2,xx3...)
)

const (
	Numeric Type = "numeric"
	String  Type = "string"
	Bool    Type = "bool"
)

/**
 * 相对复杂的条件查询 拼字符串
 * 传入查询条件，返回MySql查询语句
 *
 * 返回的值是下面XXX的内容
 * select * from table where (andXXX) or (orXXX) ...
 */
func MakeSqlCmd(cond *Cond) (and, or string, err error) {
	var (
		ands []string
		ors  []string
	)

	for _, column := range cond.And {
		str, err := makeEachSqlCmd(column)
		if err != nil {
			return "", "", err
		}
		ands = append(ands, str)
	}

	//and 条件必填，否则生成的语句就会变为 SELECT * FORM `table` WHERE OR orxxx;
	if len(ands) < 1 {
		return "", "", errors.New("缺少必填项")
	}

	for _, column := range cond.Or {
		str, err := makeEachSqlCmd(column)
		if err != nil {
			return "", "", err
		}
		ors = append(ors, str)
	}

	and = conv.ArrayToString(ands, " AND ", [2]string{})

	if len(ors) >= 1 {
		or = conv.ArrayToString(ors, " AND ", [2]string{})
	}

	return
}

func makeEachSqlCmd(c *Column) (str string, err error) {

	//---------------------------------
	//数据类型处理
	//---------------------------------
	switch c.Type {
	case Numeric:
	case String:
	case Bool:
		break
	default:
		return "", errors.New("数据类型不正确")
	}

	//---------------------------------
	//默认值处理
	//---------------------------------

	//因为 所以值都是字符串
	//所以 为了防止用户输入空字符串，需要把空字符串设置默认值
	for i, value := range c.Values {
		//如果类型为 numeric 并且 （为空 或者 不是有效的十进制数字）
		if c.Type == Numeric && (value == "" || !check.IsDecimal(value)) {
			//数字 设置默认值为 0
			c.Values[i] = "0"

			//又如果 类型为 bool 并且 （值为空 或者 值不为 true 或 false）
		} else if c.Type == Bool {
			//布尔 设置默认值为 false
			if value != "true" {
				c.Values[i] = "false"
			}
		}
	}

	//---------------------------------
	//筛选条件处理
	//---------------------------------

	//开始筛选条件

	switch c.Oper {
	case EQ:

		if c.Type == String {
			//如果是字符串，值需要用 单引号 括起来
			return fmt.Sprintf("`%s` = '%s'", c.Col, c.Values[0]), nil
		}

		return fmt.Sprintf("`%s` = %s", c.Col, c.Values[0]), nil

	case NEQ:

		if c.Type == String {
			//如果是字符串，值需要用 单引号 括起来
			return fmt.Sprintf("`%s` <> '%s'", c.Col, c.Values[0]), nil
		}

		return fmt.Sprintf("`%s` <> %s", c.Col, c.Values[0]), nil

	case GT:

		if c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [大于] 运算不能用于 布尔类型")
		}

		if c.Values[0] == "" {
			return "", errors.New("[" + c.Col + "] [大于等于] 运算 值不能为空")
		}

		if c.Type == String {
			//如果是字符串，值需要用 单引号 括起来
			return fmt.Sprintf("`%s` > '%s'", c.Col, c.Values[0]), nil
		}

		return fmt.Sprintf("`%s` > %s", c.Col, c.Values[0]), nil

	case GTE:

		if c.Type == String || c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [大于等于] 运算不能用于 字符串类型或布尔类型")
		}

		if c.Values[0] == "" {
			return "", errors.New("[" + c.Col + "] [大于等于] 运算 值不能为空")
		}

		if c.Type == String {
			//如果是字符串，值需要用 单引号 括起来
			return fmt.Sprintf("`%s` >= '%s'", c.Col, c.Values[0]), nil
		}

		return fmt.Sprintf("`%s` >= %s", c.Col, c.Values[0]), nil

	case LT:

		if c.Type == String || c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [小于] 运算不能用于 字符串类型或布尔类型")
		}

		if c.Values[0] == "" {
			return "", errors.New("[" + c.Col + "] [小于] 运算 值不能为空")
		}

		if c.Type == String {
			//如果是字符串，值需要用 单引号 括起来
			return fmt.Sprintf("`%s` < '%s'", c.Col, c.Values[0]), nil
		}

		return fmt.Sprintf("`%s` < %s", c.Col, c.Values[0]), nil

	case LTE:

		if c.Type == String || c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [小于等于] 运算不能用于 字符串类型或布尔类型")
		}

		if c.Values[0] == "" {
			return "", errors.New("[" + c.Col + "] [小于等于] 运算 值不能为空")
		}

		if c.Type == String {
			//如果是字符串，值需要用 单引号 括起来
			return fmt.Sprintf("`%s` <= '%s'", c.Col, c.Values[0]), nil
		}

		return fmt.Sprintf("`%s` <= %s", c.Col, c.Values[0]), nil

	case BETWEEN:

		if c.Type == String || c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [介于两者之间] 运算不能用于 布尔类型")
		}

		if c.Values[0] == "" || c.Values[1] == "" {
			return "", errors.New("[" + c.Col + "] [介于两者之间] 两个值都不能为空")
		}

		if c.Type == String {
			//如果是字符串，值需要用 单引号 括起来
			return fmt.Sprintf("`%s` BETWEEN '%s' AND '%s'", c.Col, c.Values[0], c.Values[1]), nil
		}

		return fmt.Sprintf("`%s` BETWEEN %s AND %s", c.Col, c.Values[0], c.Values[1]), nil

	case NBETWEEN:

		if c.Type == String || c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [介于两者之间] 运算不能用于 布尔类型")
		}

		if c.Values[0] == "" || c.Values[1] == "" {
			return "", errors.New("[" + c.Col + "] [介于两者之间] 两个值都不能为空")
		}

		if c.Type == String {
			//如果是字符串，值需要用 单引号 括起来
			return fmt.Sprintf("`%s` NOT BETWEEN '%s' AND '%s'", c.Col, c.Values[0], c.Values[1]), nil
		}

		return fmt.Sprintf("`%s` NOT BETWEEN %s AND %s", c.Col, c.Values[0], c.Values[1]), nil

	case LIKE:

		if c.Type == Numeric || c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [相似] 运算不能用于 数字类型或布尔类型")
		}

		if c.Values[0] == "" {
			return "", errors.New("[" + c.Col + "] [相似] 运算 值不能为空")
		}

		return fmt.Sprintf("`%s` LIKE '%%%s%%'", c.Col, c.Values[0]), nil

	case NLIKE:

		if c.Type == Numeric || c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [不相似] 运算不能用于 数字类型或布尔类型")
		}

		if c.Values[0] == "" {
			return "", errors.New("[" + c.Col + "] [不相似] 运算 值不能为空")
		}

		return fmt.Sprintf("`%s` NOT LIKE '%%%s%%'", c.Col, c.Values[0]), nil

	case IN:

		var str string
		if c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [包含] 运算不能用于 布尔类型")
		} else if c.Type == Numeric {
			str = strings.Replace(strings.Trim(fmt.Sprint(c.Values), "[]"), " ", ",", -1)
		} else {
			str = conv.ArrayToString(c.Values, ",", [2]string{"'", "'"})
		}

		return fmt.Sprintf("`%s` IN (%s)", c.Col, str), nil

	case NIN:

		var str string
		if c.Type == Bool {
			return "", errors.New("[" + c.Col + "] [不包含] 运算不能用于 布尔类型")
		} else if c.Type == Numeric {
			str = strings.Replace(strings.Trim(fmt.Sprint(c.Values), "[]"), " ", ",", -1)
		} else {
			str = conv.ArrayToString(c.Values, ",", [2]string{"'", "'"})
		}

		return fmt.Sprintf("`%s` NIN (%s)", c.Col, str), nil

	case NULL:
		return fmt.Sprintf("`%s` IS NULL", c.Col), nil

	case NNULL:
		return fmt.Sprintf("`%s` IS NOT NULL", c.Col), nil

	default:
		return "", errors.New("[" + c.Col + "] 无效的运算符")
	}
}

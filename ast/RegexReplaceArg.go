package ast

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

const Adapter_FormateDate = `2006-01-02 15:04:05`

//执行替换操作
func Replace(findStrs []string, data string, arg map[string]interface{}, engine ExpressionEngine, arg_array *[]interface{}) (string, error) {
	for _, findStr := range findStrs {
		var choice interface{}
		var argValue = arg[findStr]
		if argValue != nil {
			choice = argValue
			*arg_array = append(*arg_array, argValue)
		} else {
			//exec lexer
			var err error
			evalData, err := engine.LexerAndEval(findStr, arg)
			if err != nil {
				return "", errors.New(engine.Name() + ":" + err.Error())
			}
			choice = evalData
			*arg_array = append(*arg_array, evalData)
		}
		data = strings.Replace(data, "#{"+findStr+"}", Convert(choice), -1)
		//data = strings.Replace(data, "#{"+findStr+"}", indexConvert.Convert(), -1)
	}
	return data, nil
}

//执行替换操作
func ReplaceRaw(findStrs []string, data string, typeConvert SqlArgTypeConvert, arg map[string]interface{}, engine ExpressionEngine) (string, error) {
	for _, findStr := range findStrs {
		var evalData interface{}
		//find param arg
		var argValue = arg[findStr]
		if argValue != nil {
			evalData = argValue
		} else {
			//exec lexer
			var err error
			evalData, err = engine.LexerAndEval(findStr, arg)
			if err != nil {
				return "", errors.New(engine.Name() + ":" + err.Error())
			}
		}
		var resultStr string
		if typeConvert != nil {
			resultStr = typeConvert.Convert(evalData)
		} else {
			resultStr = fmt.Sprint(evalData)
		}
		data = strings.Replace(data, "${"+findStr+"}", resultStr, -1)
	}
	arg = nil
	typeConvert = nil
	return data, nil
}

//find like #{*} value *
func FindExpress(str string) []string {
	var finds = []string{}
	var item []byte
	var lastIndex = -1
	var startIndex = -1
	var strBytes = []byte(str)
	for index, v := range strBytes {
		if v == 35 {
			lastIndex = index
		}
		if v == 123 && lastIndex == (index-1) {
			startIndex = index + 1
		}
		if v == 125 && startIndex != -1 {
			item = strBytes[startIndex:index]

			//去掉逗号之后的部分
			if bytes.Contains(item, []byte(",")) {
				item = bytes.Split(item, []byte(","))[0]
			}
			finds = append(finds, string(item))
			item = nil
			startIndex = -1
			lastIndex = -1
		}
	}
	item = nil
	strBytes = nil

	var strs = []string{}
	for _, k := range finds {
		strs = append(strs, k)
	}
	return strs
}

//find like ${*} value *
func FindRawExpressString(str string) []string {
	var finds = []string{}
	var item []byte
	var lastIndex = -1
	var startIndex = -1
	var strBytes = []byte(str)
	for index, v := range str {
		if v == 36 {
			lastIndex = index
		}
		if v == 123 && lastIndex == (index-1) {
			startIndex = index + 1
		}
		if v == 125 && startIndex != -1 {
			item = strBytes[startIndex:index]
			//去掉逗号之后的部分
			if bytes.Contains(item, []byte(",")) {
				item = bytes.Split(item, []byte(","))[0]
			}
			finds = append(finds, string(item))
			item = nil
			startIndex = -1
			lastIndex = -1
		}
	}
	item = nil
	strBytes = nil

	var strs = []string{}
	for _, k := range finds {
		strs = append(strs, k)
	}
	return strs
}

func Convert(argValue interface{}) string {
	if argValue == nil {
		return "''"
	}
	switch argValue.(type) {
	case string:
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(argValue.(string))
		argStr.WriteString(`'`)
		return argStr.String()
	case *string:
		var v = argValue.(*string)
		if v == nil {
			return "''"
		}
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(*v)
		argStr.WriteString(`'`)
		return argStr.String()
	case bool:
		if argValue.(bool) {
			return "true"
		} else {
			return "false"
		}
	case *bool:
		var v = argValue.(*bool)
		if v == nil {
			return "''"
		}
		if *v {
			return "true"
		} else {
			return "false"
		}
	case time.Time:
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(argValue.(time.Time).Format(Adapter_FormateDate))
		argStr.WriteString(`'`)
		return argStr.String()
	case *time.Time:
		var timePtr = argValue.(*time.Time)
		if timePtr == nil {
			return "''"
		}
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(timePtr.Format(Adapter_FormateDate))
		argStr.WriteString(`'`)
		return argStr.String()

	case int, int16, int32, int64, float32, float64:
		return fmt.Sprint(argValue)
	case *int:
		var v = argValue.(*int)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *int16:
		var v = argValue.(*int16)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *int32:
		var v = argValue.(*int32)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *int64:
		var v = argValue.(*int64)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *float32:
		var v = argValue.(*float32)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *float64:
		var v = argValue.(*float64)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	}

	if argValue == nil {
		return ""
	}
	return fmt.Sprint(argValue)
}

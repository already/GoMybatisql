package ast

import (
	"github.com/already/batisql/v7/utils"
)

//判断节点
type NodeIf struct {
	childs []Node
	t      NodeType
	test   string

	holder *NodeConfigHolder
}

func (it *NodeIf) Type() NodeType {
	return NIf
}

func (it *NodeIf) Eval(env map[string]interface{}, arg_array *[]interface{}) ([]byte, error) {
	if it.holder == nil {
		return nil, nil
	}
	var result, err = it.holder.GetExpressionEngineProxy().LexerAndEval(it.test, env)
	if err != nil {
		err = utils.NewError("GoMybatisSqlBuilder", "[GoMybatis] <test `", it.test, `> fail,`, err.Error())
	}
	if result.(bool) {
		return DoChildNodes(it.childs, env, arg_array)
	}
	return nil, nil
}

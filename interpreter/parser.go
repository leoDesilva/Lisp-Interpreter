package interpreter 

import (
    "fmt"
    "strconv"
)

type Parser struct {
    Tokens []Node
    Token Node 
    Index int
}

func next(p *Parser){
    p.Index++
    if (p.Index >= len(p.Tokens)){
        p.Index = -1
        return
    }
    p.Token = p.Tokens[p.Index]
}

func Parse(p *Parser) ([]interface{}, bool){
    var ast []interface{}

    for p.Index != -1 && p.Index < len(p.Tokens) {
        node, err := parseExpr(p)
        if err {return ast, true}
        fmt.Println(node)
        ast = append(ast, node)
        next(p)
    }

    return ast, false
}

func parseExpr(p *Parser) (interface{},bool) {
    if p.Token.Type == TT_LPAREN{
        node,err := parseList(p)
        if err { return node, true }
        return node, false
    } else{
        node,err := parseAtom(p)
        if err { return node, true }
        return node, false
    }

    return EmptyNode{TT_EOF}, true
}

func parseList(p *Parser) (interface{}, bool){
         //[x] (OP expr expr*) ;OP = +/-/*/>/</>=
         //[x] ( setf identifier expr)
         //[x] ( identifier {expr}) ; function call 
         //[ ] ( defun identifier ({identifier}) expr)
         //[ ] ( if expr expr expr) ; condition then else
         //[ ] ( block expr* )
         //[x] ( expr )
         next(p)

         if p.Token.Class == "BIN_OP"{
             node, err := parseBinOp(p)
             if err {return EmptyNode{TT_EOF}, true}
             return node, false
         } else if matches(p.Token, "setf"){
             node,err := parseAssignment(p)
             if err {return EmptyNode{TT_EOF}, true}
             return node,false
         } else if p.Token.Type == TT_IDENTIFIER {
             node, err := parseFunctionCall(p)
             if err {return EmptyNode{TT_EOF}, true}
             return node,false
         } else if p.Token.Type == TT_LPAREN {
             node,err := parseExpr(p)
             if err {return EmptyNode{TT_EOF}, true}
             next(p)
             return node, false
         } else if p.Token.Class == "ATOM"{
             node,err := parseAtom(p)
             if err {return EmptyNode{TT_EOF}, true}
             next(p)
             return node, false
         }


         return EmptyNode{TT_EOF}, true
}

func parseFunctionCall(p *Parser) (interface{}, bool){
    identifier := p.Token.Value
    next(p)
    
    var parameters []interface{}

    for p.Token.Type != TT_RPAREN {
        param,err := parseExpr(p) 
        if err {return EmptyNode{TT_EOF}, true}
        parameters = append(parameters,param)
        next(p)
    }

    return FunctionCallNode{TT_FUNCTION_CALL, identifier, parameters}, false
    
}

func parseAssignment(p *Parser) (interface{}, bool) {
    next(p)

    if p.Token.Type != TT_IDENTIFIER {return EmptyNode{TT_EOF}, true}
    identifier := p.Token.Value
    next(p)

    value,err := parseExpr(p)
    if err {return EmptyNode{TT_EOF}, true}
    next(p)

    return VarAssignNode{TT_VAR_ASSIGN, identifier, value}, false
}

func parseBinOp(p *Parser) (interface{}, bool){
    op := p.Token.Type
    next(p)
    operandA,err := parseExpr(p)
    next(p)
    operandB,err := parseExpr(p)
    next(p)

    if err { return EmptyNode{TT_EOF}, true}

    return BinOpNode{TT_BIN_OP, operandA, op, operandB}, false
}

func parseAtom(p *Parser)(interface{}, bool){
    var node interface{}

    if p.Token.Type == TT_EOF{
        return EmptyNode{TT_EOF}, false

    } else if p.Token.Type == TT_INT{
        value, _ := strconv.Atoi(p.Token.Value)
        return IntNode{
            Type: TT_INT,
            Value: value,
        }, false

    } else if p.Token.Type == TT_STRING{
        return StringNode{
            Type: TT_STRING,
            Value: p.Token.Value,
        }, false

    } else if p.Token.Type == TT_IDENTIFIER{
        return VarAcessNode {
            Type: TT_VAR_ACCESS,
            Identifier: p.Token.Value,
        }, false
    }

    return node, true
}


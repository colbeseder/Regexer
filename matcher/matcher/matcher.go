package matcher

import (
    "fmt"
    "strings"
    "encoding/json"
)

const (
    START = 0
    END = 1
    LITERAL = 2
    GROUP = 3
    WILD = 4
)

type Node struct {
    Flavour int
    Value string
    Children []*Node
    Max int
    Min int
}

func (n *Node) addChild(child *Node) {
    n.Children = append(n.Children, child)
}

func (n *Node) print(indent string) {
    fmt.Println(indent, n.Value, " [", n.Min, ",", n.Max, "]")
    for i:=0 ; i < len(n.Children) ; i++ {
        n.Children[i].print(indent + "--")
    }
}

func (n *Node) toJSON() string{
    treeAsJSON, err := json.Marshal(n)
    if err != nil {
        panic(err)
    }
    return string(treeAsJSON)
}

func createNode(Flavour int, Value string, Max int, Min int) Node {
    return Node{Flavour: Flavour, Value: Value, Max: Max, Min: Min}
}

func Tokenize(pattern string) []string {
    // So far, just literal characters and groups
    var tokens []string
    for i := 0; i < len(pattern); i++ {
        var token string
        if string(pattern[i]) == "[" {
            //group
            i++
            group := ""
            for  ; string(pattern[i]) != "]" ; i++ {
                group += string(pattern[i])
            }
            token = "[" + group + "]"
        } else {
            //literal
            token = string(pattern[i])
        }
        tokens = append(tokens, token)
    }
    return tokens
}

func Tokens2AST(tokens []string, head *Node){
    current := head
    for i := 0 ; i < len(tokens) ; i++ {
        token := tokens[i]
        var next Node
        if token[0] == "["[0] {
            //group
            next = createNode(GROUP, token[1:len(token)-1], 1, 1)
        } else if token == "." {
            next = createNode(WILD, ".", 1, 1)
        } else {
            next = createNode(LITERAL, token, 1, 1)
        }
        current.addChild(&next)
        current = &next
    }
}

func pattern2AST(pattern string) *Node {
    tokens := Tokenize(pattern)
    //fmt.Println(tokens)
    head := createNode(START, "START", 0, 0)
    Tokens2AST(tokens, &head)
    return &head
}

func doesMatchNode(c rune, n *Node) bool {
    switch Flavour := n.Flavour; Flavour {
    case LITERAL :
        return rune(n.Value[0]) == c
    case GROUP :
        return strings.ContainsRune(n.Value, c)
    case WILD :
        return true
    default:
        return false
    }
}

func MatchTree(head *Node, haystack []rune, offset int) (bool, string) {
    pStartMatching := offset
    pHaystack := pStartMatching
    current := head
    countMatchesOnNode := 0
    capture := ""
    for {
        if countMatchesOnNode == current.Max {
            if len(current.Children) == 0 {
                return true, capture
            } else {
                current = current.Children[0]
                countMatchesOnNode = 0
            }
        }
        if pHaystack >= len(haystack) {
            break
        }
        if doesMatchNode(rune(haystack[pHaystack]), current){
            capture += string(haystack[pHaystack])
            countMatchesOnNode += 1
            pHaystack += 1
        } else {
            current = head
            countMatchesOnNode = 0
            capture = ""
            pStartMatching += 1
            pHaystack = pStartMatching
        }
    }
    return false , ""
}

func Match(pattern string, haystack string) (bool, string) {
    head := pattern2AST(pattern)
    head.toJSON()
    //head.print("")
    return MatchTree(head, []rune(haystack), 0)
}

func test(pattern string, haystack string, expectedR bool, expectedC string){
    r, c := Match(pattern, haystack)
    if r != expectedR || c != expectedC {
        fmt.Println("Test failed: ")
        fmt.Println("expected", expectedR, expectedC, " , But got" , r, c)
    }
}

func RunTests(){
    test("hello", "hello", true, "hello")
    test("hello", "goodbye", false, "")
    test("h[ea]llo", "hello", true, "hello")
    test("h[ea]llo", "hollo", false, "")
    test("h[ea]llo", "aaahellobbb", true, "hello")
    test("h[ea]llo", "aaahhellobbb", true, "hello")
    test("h.llo", "hello", true, "hello")
    test("h.llo", "hllo", false, "")
    test(".....", "aaahhellobbb", true, "aaahh")
    test("..hello", "aaahhellobbb", true, "ahhello")
}

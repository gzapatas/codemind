package ingestion

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	ts_go "github.com/smacker/go-tree-sitter/golang"
	ts_js "github.com/smacker/go-tree-sitter/javascript"
	ts_rust "github.com/smacker/go-tree-sitter/rust"
	ts_ts "github.com/smacker/go-tree-sitter/typescript"
)

func init() {
	// register tree-sitter chunker with higher priority so it's preferred
	RegisterChunkerWithPriority("go", treeSitterChunker, 20)
	RegisterChunkerWithPriority("rust", treeSitterChunker, 20)
	RegisterChunkerWithPriority("javascript", treeSitterChunker, 20)
	RegisterChunkerWithPriority("js", treeSitterChunker, 20)
	RegisterChunkerWithPriority("typescript", treeSitterChunker, 20)
	RegisterChunkerWithPriority("ts", treeSitterChunker, 20)
}

func treeSitterChunker(content []byte, lang string) ([]Chunk, error) {
	p := sitter.NewParser()
	switch lang {
	case "go":
		p.SetLanguage(ts_go.GetLanguage())
	case "rust":
		p.SetLanguage(ts_rust.GetLanguage())
	case "javascript", "js":
		p.SetLanguage(ts_js.GetLanguage())
	case "typescript", "ts":
		p.SetLanguage(ts_ts.GetLanguage())
	default:
		return nil, fmt.Errorf("treesitter: unsupported language: %s", lang)
	}

	tree, err := p.ParseCtx(context.Background(), nil, content)
	if err != nil || tree == nil {
		return nil, fmt.Errorf("treesitter parse error: %v", err)
	}

	var out []Chunk
	var walk func(node *sitter.Node)
	walk = func(node *sitter.Node) {
		t := node.Type()
		switch t {
		// common function-like nodes across grammars
		case "function_declaration", "function_item", "method_declaration", "method_item", "function":
			name := extractNameFromNode(node, content)
			out = append(out, makeTSChunk(node, "function", name, content))
		// type-like nodes
		case "type_declaration", "type_spec", "type_item", "struct_item", "enum_item", "class_declaration", "class":
			name := extractNameFromNode(node, content)
			out = append(out, makeTSChunk(node, "type", name, content))
		}
		for i := 0; i < int(node.NamedChildCount()); i++ {
			walk(node.NamedChild(i))
		}
	}

	walk(tree.RootNode())
	return out, nil
}

func makeTSChunk(n *sitter.Node, kind, name string, content []byte) Chunk {
	start := n.StartPoint().Row + 1
	end := n.EndPoint().Row + 1
	text := string(content[n.StartByte():n.EndByte()])
	return Chunk{Kind: kind, Name: name, StartLine: int(start), EndLine: int(end), Text: text}
}

func extractNameFromNode(n *sitter.Node, content []byte) string {
	if child := n.ChildByFieldName("name"); child != nil {
		return string(content[child.StartByte():child.EndByte()])
	}
	// fall back to scanning for an identifier-like child
	for i := 0; i < int(n.NamedChildCount()); i++ {
		c := n.NamedChild(i)
		typ := c.Type()
		if typ == "identifier" || typ == "type_identifier" || typ == "property_identifier" || typ == "field_identifier" {
			return string(content[c.StartByte():c.EndByte()])
		}
	}
	// last resort: empty name
	return ""
}

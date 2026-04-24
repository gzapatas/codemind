package ingestion

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	ts_go "github.com/smacker/go-tree-sitter/golang"
	ts_js "github.com/smacker/go-tree-sitter/javascript"
	ts_rust "github.com/smacker/go-tree-sitter/rust"
	ts_tsx "github.com/smacker/go-tree-sitter/typescript/tsx"
	ts_ts "github.com/smacker/go-tree-sitter/typescript/typescript"
)

func init() {
	// register tree-sitter chunker with higher priority so it's preferred
	RegisterChunkerWithPriority("go", treeSitterChunker, 20)
	RegisterChunkerWithPriority("rs", treeSitterChunker, 20)
	RegisterChunkerWithPriority("js", treeSitterChunker, 20)
	RegisterChunkerWithPriority("ts", treeSitterChunker, 20)
	RegisterChunkerWithPriority("tsx", treeSitterChunker, 20)
	RegisterChunkerWithPriority("jsx", treeSitterChunker, 20)
}

func treeSitterChunker(content []byte, lang string) ([]Chunk, error) {
	p := sitter.NewParser()
	switch lang {
	case "go":
		p.SetLanguage(ts_go.GetLanguage())
	case "rs":
		p.SetLanguage(ts_rust.GetLanguage())
	case "js":
		p.SetLanguage(ts_js.GetLanguage())
	case "ts":
		p.SetLanguage(ts_ts.GetLanguage())
	case "tsx":
		p.SetLanguage(ts_tsx.GetLanguage())
	case "jsx":
		p.SetLanguage(ts_tsx.GetLanguage())
	default:
		return nil, fmt.Errorf("treesitter: unsupported language: %s", lang)
	}

	tree, err := p.ParseCtx(context.Background(), nil, content)
	if err != nil || tree == nil {
		return nil, fmt.Errorf("treesitter parse error: %v", err)
	}

	// optional debug dump
	//dumpNode(tree.RootNode(), content, 0)

	var out []Chunk
	var walk func(node *sitter.Node)
	walk = func(node *sitter.Node) {
		t := node.Type()
		captured := false
		switch t {
		// common function-like nodes across grammars
		case "function_declaration", "function_item", "method_declaration", "method_item", "function", "method_definition", "method_signature", "constructor":
			name := extractNameFromNode(node, content)
			out = append(out, makeTSChunk(node, "function", name, content))
			captured = true
		// type-like nodes
		case "type_declaration", "type_spec", "type_item", "struct_item", "enum_item", "class_declaration", "class":
			name := extractNameFromNode(node, content)
			out = append(out, makeTSChunk(node, "type", name, content))
			// also extract method-like descendants inside the class/type (limited depth)
			var search func(n *sitter.Node, depth int)
			search = func(n *sitter.Node, depth int) {
				if n == nil || depth <= 0 {
					return
				}
				typ := n.Type()
				switch typ {
				case "method_definition", "method_declaration", "method_item", "function", "function_declaration", "function_item", "constructor":
					mname := extractNameFromNode(n, content)
					out = append(out, makeTSChunk(n, "method", mname, content))
					return
				}
				for i := 0; i < int(n.NamedChildCount()); i++ {
					search(n.NamedChild(i), depth-1)
				}
			}
			for i := 0; i < int(node.NamedChildCount()); i++ {
				search(node.NamedChild(i), 4)
			}
			captured = true
		}
		if captured {
			return
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

func dumpNode(n *sitter.Node, content []byte, depth int) {
	if n == nil {
		return
	}
	if depth > 50 {
		return
	}
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	start := n.StartPoint().Row + 1
	end := n.EndPoint().Row + 1
	snippet := string(content[n.StartByte():n.EndByte()])
	if len(snippet) > 100 {
		snippet = snippet[:100] + "..."
	}
	fmt.Printf("%s- %s (%d-%d): %q\n", indent, n.Type(), start, end, snippet)
	for i := 0; i < int(n.NamedChildCount()); i++ {
		dumpNode(n.NamedChild(i), content, depth+1)
	}
}

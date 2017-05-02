package agent

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	//"regexp"
	//"strconv"
	//"strings"
	//"time"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/xytis/congix/structs/config"
)

type Config struct {
	//Mapping *config.MappingConfig `mapstructure:"mapping"`
	Nginx  *config.NginxConfig  `mapstructure:"nginx"`
	Consul *config.ConsulConfig `mapstructure:"consul"`
}

func Parse(r io.Reader) (*Config, error) {
	fmt.Println("parsing a reader")
	// Copy the reader into an in-memory buffer first since HCL requires it.
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}

	// Parse the buffer
	root, err := hcl.Parse(buf.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing: %s", err)
	}
	buf.Reset()
	fmt.Println("parsed everything")

	ast.Walk(root, func(node ast.Node) (ast.Node, bool) {
		if node == nil {
			return nil, false
		}
		switch n := node.(type) {
		case *ast.File:
			fmt.Printf("file: %v\n", n.Pos())
		case *ast.ObjectList:
			fmt.Printf("ObjectList: %v\n", n.Pos())
		case *ast.ObjectKey:
			fmt.Printf("ObjectKey: %v\n", n.Token)
		case *ast.ObjectItem:
			fmt.Printf("ObjectItem: %v\n", n.Pos())
		case *ast.LiteralType:
			fmt.Printf("LiteralType: %v\n", n.Token)
		case *ast.ListType:
			fmt.Printf("ListType: %v\n", n.Pos())
		case *ast.ObjectType:
			fmt.Printf("ObjectType: %v\n", n.Pos())
		default:
			fmt.Printf("unknown type: %T\n", n)
		}
		return node, true
	})

	return nil, nil
}

func ParseFile(path string) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

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

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
	"github.com/xytis/congix/helper"
	"github.com/xytis/congix/structs/config"
)

type Config struct {
	// Configuration to reach and control nginx plus server
	Nginx *config.NginxConfig `mapstructure:"nginx"`

	// Runtime mapping which control current nginx state
	Mapping *config.Mapping `mapstructure:"mapping"`

	// List of config files that have been loaded (in order)
	Files []string `mapstructure:"-"`
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

	// Top-level item must be a list
	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: root must be an object")
	}

	// Check for invalid keys
	valid := []string{
		"nginx",
		"mapping",
	}
	if err := checkForInvalidKeys(list, valid); err != nil {
		return nil, err
	}

	config := Config{
		Mapping: &config.Mapping{},
		Nginx:   &config.NginxConfig{},
	}

	{
		matches := list.Filter("nginx")
		if len(matches.Items) > 0 {
			if err := parseNginxStanza(&config.Nginx, matches); err != nil {
				return nil, fmt.Errorf("error parsing 'nginx': %s", err)
			}
		}
	}

	{
		matches := list.Filter("mapping")
		if len(matches.Items) < 1 {
			return nil, fmt.Errorf("missing 'mapping' stanza")
		}
		if err := parseMappingStanza(&config.Mapping, matches); err != nil {
			return nil, fmt.Errorf("error parsing 'mapping': %s", err)
		}
	}

	return &config, nil
}

func LoadConfig(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// TODO:
	//if fi.IsDir() {
	//	return LoadConfigDir(path)
	//}

	cleaned := filepath.Clean(path)
	config, err := ParseConfigFile(cleaned)
	if err != nil {
		return nil, fmt.Errorf("Error loading %s: %s", cleaned, err)
	}

	config.Files = append(config.Files, cleaned)
	return config, nil
}

func ParseConfigFile(path string) (*Config, error) {
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

func checkForInvalidKeys(node ast.Node, valid []string) error {
	var list *ast.ObjectList
	switch n := node.(type) {
	case *ast.ObjectList:
		list = n
	case *ast.ObjectType:
		list = n.List
	default:
		return fmt.Errorf("cannot validate keys of type %T", n)
	}

	validMap := make(map[string]struct{}, len(valid))
	for _, v := range valid {
		validMap[v] = struct{}{}
	}

	var result error
	for _, item := range list.Items {
		key := item.Keys[0].Token.Value().(string)
		if _, ok := validMap[key]; !ok {
			result = multierror.Append(result, fmt.Errorf(
				"invalid key: %s", key))
		}
	}

	return result
}

func parseNginxStanza(result **config.NginxConfig, list *ast.ObjectList) error {
	if len(list.Items) != 1 {
		return fmt.Errorf("only one 'nginx' stanza allowed")
	}

	// Get nginx config object
	obj := list.Items[0]

	if _, ok := obj.Val.(*ast.ObjectType); !ok {
		return fmt.Errorf("'nginx' stanza should be an object")
	}

	valid := []string{
		"address",
		"status_endpoint",
		"upstream_endpoint",
	}
	if err := checkForInvalidKeys(obj.Val.(*ast.ObjectType).List, valid); err != nil {
		return multierror.Prefix(err, "nginx:")
	}

	// Decode the full thing into a map[string]interface for ease
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, obj.Val); err != nil {
		return err
	}

	// Decode the rest
	if err := mapstructure.WeakDecode(m, result); err != nil {
		return err
	}

	// At this point result object is populated with information from file (if any)

	return nil
}

func parseMappingStanza(result **config.Mapping, list *ast.ObjectList) error {
	for _, obj := range list.Items {
		var listVal *ast.ObjectList
		if ot, ok := obj.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return fmt.Errorf("'mapping' stanza should be an object")
		}

		if o := listVal.Filter("entry"); len(o.Items) > 0 {
			for _, item := range o.Items {
				entry := &config.Entry{}
				if err := parseMappingEntry(&entry, item); err != nil {
					return multierror.Prefix(err, "mapping:")
				}
				(*result).Entries = append((*result).Entries, entry)
			}
		}
	}

	return nil
}

func parseMappingEntry(result **config.Entry, item *ast.ObjectItem) error {
	name := item.Keys[0].Token.Value().(string)

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return err
	}
	// Decode the rest
	if err := mapstructure.WeakDecode(m, result); err != nil {
		return err
	}
	(*result).Name = helper.StringToPtr(name)

	return nil
}

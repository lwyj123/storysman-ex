// Copyright (c) 2016-2017 Brandon Buck

package tmpl

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gobuffalo/velvet"
	"github.com/lwyj123/storysman-ex/logger"
)

const (
	openTemplateTags  = "{{"
	closeTemplateTags = "}}"
)

var compiledTemplates = make(map[string]Renderer)

// Register will compile and register a template using the string given and
// store the compiled template in the map.
func Register(name, contents string) error {
	r, err := getVelvetRenderer(contents)
	if err != nil {
		return err
	}
	compiledTemplates[name] = r

	return nil
}

// RegisterFile will compile and register a template using the contents of a
// file and storing by name in the compiled map.
func RegisterFile(name, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return Register(name, string(contents))
}

// Unregister removes a previously compiled template, normally only used for
// testing purposes but in cases a template is, for some reason, miscompiled
// (although longer term solutions should be much more preferred than
// unregistering)
func Unregister(name string) {
	delete(compiledTemplates, name)
}

// RenderOnce bypasses storing the saved templates. Tihs should only be used
// if the template is only ever going to be compiled and rendered _one time_.
func RenderOnce(contents string, data interface{}) (string, error) {
	r, err := getVelvetRenderer(contents)
	if err != nil {
		return "", err
	}
	result, err := r.Render(data)

	return result, err
}

// MustRenderOnce performs a RenderOnce and will exit the program when an error
// occurs.
func MustRenderOnce(contents string, data interface{}) string {
	result, err := RenderOnce(contents, data)
	if err != nil {
		logger.NewWithSource("tmpl").WithError(err).Fatal("Failed to render template once.")
	}

	return result
}

// Template returns the Renderer associated with a registered template, if any.
func Template(name string) (Renderer, error) {
	if tmpl, ok := compiledTemplates[name]; ok {
		return tmpl, nil
	}

	return nil, fmt.Errorf("No template has been registered with the name \"%s\"", name)
}

func getVelvetRenderer(contents string) (Renderer, error) {
	template, err := velvet.NewTemplate(contents)
	if err != nil {
		return nil, err
	}

	return &velvetRenderer{template}, nil
}

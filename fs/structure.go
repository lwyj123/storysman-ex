package fs

import (
	"errors"
	"path/filepath"

	"github.com/lwyj123/storysman-ex/logger"
)

// ProjectStructure is a skeleton of what a `DragonMUD` project should look
// like after `dragon init` is called.
var ProjectStructure = Dir{
	"Dragonfile.toml": File{},
	".gitignore":      File{},
	"plugins":         Dir{},
	"commands": Dir{
		"init.lua": File{},
	},
	"server": Dir{
		"init.lua": File{},
	},
	"client": Dir{
		"init.lua": File{},
	},
	"views": Dir{},
}

// PluginStructure represents what a plugin is intended to look like.
var PluginStructure = Dir{
	"DragonInfo.toml": File{},
	"commands": Dir{
		"init.lua": File{},
	},
	"server": Dir{
		"init.lua": File{},
	},
	"client": Dir{
		"init.lua": File{},
	},
	"views": Dir{},
}

// CreateStructureParams makes it easier and more meaningful to call
// CreateFromStructure.
type CreateStructureParams struct {
	Log          logger.Log
	BaseName     string
	Structure    Item
	TemplateData interface{}
}

// CreateFromStructure takes in an fs.Item and generates the filesystem records
// (files and directories) according the structure provided. This is a
// recursive operation, calling itself for all instances of ItemTypeDir.
func CreateFromStructure(params CreateStructureParams) error {
	if params.Log == nil || params.BaseName == "" || params.Structure == nil {
		return errors.New("missing Log, BaseName or Structure from CreateForStructure params")
	}

	if params.Structure.Type() != ItemTypeDir {
		return errors.New("structure given is not a directory, cannot create from file")
	}

	for name, item := range params.Structure.(Dir) {
		fpath := filepath.Join(params.BaseName, name)
		item.Create(params.Log, fpath, params.TemplateData)
	}

	return nil
}

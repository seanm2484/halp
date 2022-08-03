package cheats

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"gopkg.in/yaml.v3"
)

type Client struct {
}

type CheatSearch struct {
	TitleField string   `yaml:"description"`
	Desc       string   `yaml:"command"`
	Vars       []string `yaml:"variables"`
	Filename   string   `yaml:"file"`
}

func (i CheatSearch) Title() string       { return i.TitleField }
func (i CheatSearch) Description() string { return i.Desc }
func (i CheatSearch) FilterValue() string { return i.TitleField }
func (i CheatSearch) Variables() []string { return i.Vars }
func (i CheatSearch) File() string        { return i.Filename }

func GetList(filePath string) []list.Item {
	files := func(root string, ext string) []string {
		var a []string
		filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if filepath.Ext(d.Name()) == ext {
				a = append(a, s)
			}
			return nil
		})
		return a
	}(filePath, ".yaml")

	cheats := make([]list.Item, 0)
	for _, f := range files {
		// get all the cheats that we could ever find
		cheats = append(cheats, loadYAMLItems(f)...)
	}

	return cheats
}

func (mw CheatSearch) GetDescriptions(filePath string) []CheatSearch {
	files := func(root string, ext string) []string {
		var a []string
		filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if filepath.Ext(d.Name()) == ext {
				a = append(a, s)
			}
			return nil
		})
		return a
	}(filePath, ".yaml")

	cheats := make([]CheatSearch, 0)
	for _, f := range files {
		// get all the cheats that we could ever find
		cheats = append(cheats, mw.loadYAML(f)...)
	}

	return cheats
}

func loadYAMLItems(filePath string) []list.Item {
	f, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var cheats []CheatSearch
	if err := yaml.Unmarshal(f, &cheats); err != nil {
		panic(err)
	}

	items := make([]list.Item, 0)
	for _, cheat := range cheats {
		items = append(items, cheat)
	}
	return items
}

func (mw CheatSearch) loadYAML(filePath string) []CheatSearch {
	f, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var cheats []CheatSearch
	if err := yaml.Unmarshal(f, &cheats); err != nil {
		//log.Fatal(err)
		panic(err)
	}

	return cheats
}

func (mw CheatSearch) FindSelectedCheat(cheatDesc string) (cheat CheatSearch, err error) {
	// we have a description, so iterate over the list to find the entry
	// that matches.
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	cheatList := mw.GetDescriptions(dirname + "/.halp/cheatFiles")
	for _, c := range cheatList {
		if c.TitleField == cheatDesc {
			return c, nil
		}
	}
	return cheat, errors.New("not found")
}

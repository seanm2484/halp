package cheats

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"gopkg.in/yaml.v3"
)

type Client struct {
}

type CheatSearch struct {
	TitleField string   `yaml:"description"`
	Desc       string   `yaml:"command"`
	Vars       []string `yaml:"variables"`
}

func (i CheatSearch) Title() string       { return i.TitleField }
func (i CheatSearch) Description() string { return i.Desc }
func (i CheatSearch) FilterValue() string { return i.TitleField }
func (i CheatSearch) Variables() []string { return i.Vars }

func GetList(filePath string) []list.Item {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		panic(err)
	}
	cheats := make([]list.Item, 0)
	for _, f := range files {
		// get all the cheats that we could ever find
		cheats = append(cheats, loadYAMLItems(filePath+"/"+f.Name())...)
	}

	return cheats
}

func (mw CheatSearch) GetDescriptions(filePath string) []CheatSearch {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		panic(err)
	}
	cheats := make([]CheatSearch, 0)
	for _, f := range files {
		// get all the cheats that we could ever find
		cheats = append(cheats, mw.loadYAML(filePath+"/"+f.Name())...)
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
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return cheats
}

func (mw CheatSearch) FindSelectedCheat(cheatDesc string) (cheat CheatSearch, err error) {
	// we have a description, so iterate over the list to find the entry
	// that matches.
	cheatList := mw.GetDescriptions("./cheatFiles")
	for _, c := range cheatList {
		if c.TitleField == cheatDesc {
			return c, nil
		}
	}
	return cheat, errors.New("not found")
}

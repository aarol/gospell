package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eskriett/spell"
	"github.com/urfave/cli/v2"
)

type Suggestion struct {
	word         string
	replacements []string
}

type model struct {
	filename    string
	dict        *spell.Spell
	words       []string
	suggestions []Suggestion
}

func initialModel(path string) model {
	return model{
		filename: path,
	}
}

type dictMsg (*spell.Spell)
type loadError (error)

func loadDictionary() tea.Msg {
	sp, err := spell.Load("dict.spell")
	if err != nil {
		return loadError(err)
	} else {
		return dictMsg(sp)
	}
}

type wordsMsg ([]string)

func loadWords(path string) tea.Cmd {
	return func() tea.Msg {

		words := make([]string, 0)
		file, err := os.Open(path)
		if err != nil {
			return loadError(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			words = append(words, strings.Split(scanner.Text(), " ")...)
		}
		return wordsMsg(words)
	}
}

type suggestionsMsg ([]Suggestion)

func createSuggestions(m model) tea.Cmd {
	res := make([]Suggestion, 0)
	for _, v := range m.words {
		m.dict.
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(loadDictionary, loadWords(m.filename))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case dictMsg:
		m.dict = msg
	case wordsMsg:
		m.words = msg
	case loadError:
		fmt.Println("loading failed:", msg)
	}
	return m, nil
}

func (m model) View() string {
	if m.dict == nil {
		return "Loading dictionary.."
	}
	if m.words == nil {
		return "Loading words.."
	}
	return strings.Join(m.words, ", ")
}

func main() {
	app := &cli.App{
		Name:            "gospell",
		ArgsUsage:       "filename",
		HideHelpCommand: true, // no commands
		Action: func(ctx *cli.Context) error {
			path := ctx.Args().First()
			if path != "" {
				return tea.NewProgram(initialModel(path)).Start()
			}
			return cli.ShowAppHelp(ctx)
		},
	}
	app.Run(os.Args)
}

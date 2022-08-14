package file

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nxadm/tail"
	"log"
	"os"
	"path/filepath"
	"spin/handler"
	"spin/settings"
	"spin/syntax"
)

func Setup(config *settings.Config) {
	m := new(handler.Model)

	temporaryFile, err := os.CreateTemp("", fmt.Sprintf("%s-", filepath.Base(os.Args[0])))
	if err != nil {
		log.Fatal("Could not create temporary file", err)
	}

	m.TempFile = temporaryFile

	if _, err = m.TempFile.WriteString(""); err != nil {
		log.Fatal("Unable to write to temporary file", err)
	}

	////////////////////////////////////////////////////////// Tail
	fileToTail := getPathToFileToBeTailed()

	file, tailErr := tail.TailFile(
		fileToTail, tail.Config{Follow: config.Follow})
	if tailErr != nil {
		panic(err)
	}

	m.TailFile = file

	go func() {
		for line := range m.TailFile.Lines {
			syntaxHighlightedLine := syntax.Highlight(line.Text)
			_, _ = m.TempFile.WriteString(syntaxHighlightedLine + "\n")
		}
	}()
	////////////////////////////////////////////////////////// Tail

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
	}

	tpErr := m.TempFile.Close()
	if tpErr != nil {
		panic(tpErr)
	}

	tErr := m.TailFile.Stop()
	if tErr != nil {
		panic(tErr)
	}

}

func getPathToFileToBeTailed() string {
	return os.Args[1]
}

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	rootCmd, err := cmd.NewRootCmd()
	if err != nil {
		panic(err)
	}

	err = generateDocs(rootCmd, "./website/docs/reference/cli")
	if err != nil {
		panic(err)
	}
}

func generateDocs(cmd *cobra.Command, parentDirPath string) error {
	var currentFilePath string
	sideBarLabel := cmd.Name()

	// leaf commands go in the same folder as their parent
	if len(cmd.Commands()) == 0 {
		currentFilePath = filepath.Join(parentDirPath, cmd.Name()+".md")
	} else {
		// branch commands go in index.md files in their own folders
		var currentDirPath string

		// handle root cmd
		if cmd.CommandPath() == "opctl" {
			sideBarLabel = "CLI"
			currentDirPath = parentDirPath
		} else {
			currentDirPath = filepath.Join(parentDirPath, cmd.Name())
		}

		if err := os.MkdirAll(currentDirPath, 0755); err != nil {
			return err
		}

		currentFilePath = filepath.Join(currentDirPath, "index.md")
	}

	file, err := os.Create(currentFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(
		fmt.Sprintf(
			`---
sidebar_label: %s
hide_title: true
---

`,
			sideBarLabel,
		),
	)
	if err != nil {
		return err
	}

	if err := doc.GenMarkdownCustom(
		cmd,
		file,
		func(defaultFileName string) string {
			fileNameParts := strings.Split(
				strings.TrimSuffix(defaultFileName, ".md"),
				"_",
			)

			cmdPathParts := strings.Split(
				cmd.CommandPath(),
				" ",
			)

			// handle link to parent
			if len(fileNameParts) == (len(cmdPathParts) - 1) {
				// leafs go in same folder as parent
				if len(cmd.Commands()) == 0 {
					return "index.md"
				}

				// branches get their own folder
				return "../index.md"
			}

			// handle link to child
			childCmdName := fileNameParts[len(fileNameParts)-1]
			var childCmd *cobra.Command
			for _, subCmd := range cmd.Commands() {
				if subCmd.Name() == childCmdName {
					childCmd = subCmd
				}
			}

			// leafs go in same folder as parent
			if len(childCmd.Commands()) == 0 {
				return fmt.Sprintf(
					"%s.md",
					childCmdName,
				)
			}

			// branches get their own folder
			return fmt.Sprintf(
				"%s/index.md",
				childCmdName,
			)
		},
	); err != nil {
		return err
	}

	// Recursively generate docs for child commands
	for _, subCmd := range cmd.Commands() {
		// skip help and hidden commands
		if !subCmd.IsAvailableCommand() {
			continue
		}

		if err := generateDocs(subCmd, path.Dir(currentFilePath)); err != nil {
			return err
		}
	}
	return nil
}

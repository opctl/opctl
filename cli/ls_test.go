package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	cliparamsatisfierFakes "github.com/opctl/opctl/cli/internal/cliparamsatisfier/fakes"
	nodeproviderFakes "github.com/opctl/opctl/cli/internal/nodeprovider/fakes"
)

const op1 = `
name: op1
description: A single line description
`

const op2 = `
name: op2
description: |
  A multiline description

  * one
  * two
  * three
`

func TestLS(t *testing.T) {
	g := NewGomegaWithT(t)

	/* arrange */
	cliParamSatisfier := new(cliparamsatisfierFakes.FakeCLIParamSatisfier)
	nodeProvider := new(nodeproviderFakes.FakeNodeProvider)
	testDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	os.Mkdir(filepath.Join(testDir, "op1"), 0700)
	ioutil.WriteFile(filepath.Join(testDir, "op1", "op.yml"), []byte(op1), 0600)
	os.Mkdir(filepath.Join(testDir, "op2"), 0700)
	ioutil.WriteFile(filepath.Join(testDir, "op2", "op.yml"), []byte(op2), 0600)

	outputFileName := filepath.Join(testDir, "output")
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	// mock the environment - we need to be able to read something written to
	// stdout

	oldStdout := *os.Stdout
	os.Stdout = outputFile
	defer func() {
		os.Stdout = &oldStdout
	}()

	/* act */
	ls(context.Background(), cliParamSatisfier, nodeProvider, testDir)

	outputFile.Close()

	/* assert */
	output, err := ioutil.ReadFile(outputFileName)
	g.Expect(err).To(BeNil())

	// because we don't know where this is running and there's some randomness in
	// temporary directories, the spacing can't be anticipated. Rather than
	// basically reimplementing tabwriter here, I'm just checking some small
	// aspects of the output, instead of exactly comparing it to an expectation
	t.Log(string(output))
	lines := strings.Split(string(output), "\n")

	firstLine := lines[0]
	t.Log(firstLine)
	g.Expect(firstLine[:4]).To(Equal("REF\t"))
	g.Expect(firstLine[len(firstLine)-11:]).To(Equal("DESCRIPTION"))

	lastLineWithContent := lines[len(lines)-2]
	t.Log(lastLineWithContent)
	// the last line should end with the last line of op2's description
	g.Expect(lastLineWithContent[len(lastLineWithContent)-7:]).To(Equal("* three"))
	// the last line should start with tabs
	g.Expect(lastLineWithContent[:2]).To(Equal("\t\t"))
}

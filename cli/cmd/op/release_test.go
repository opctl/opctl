package op

import (
	"bytes"
	"strings"
	"testing"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

func TestConvertAndPlanAppearInOpHelp(t *testing.T) {
	cmd := NewOpCmd(cliparamsatisfier.New(nil), &local.NodeConfig{})
	var help bytes.Buffer
	cmd.SetOut(&help)
	if err := cmd.Help(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"convert", "plan"} {
		if !strings.Contains(help.String(), name) {
			t.Fatalf("%s is absent from op help: %s", name, &help)
		}
	}
}

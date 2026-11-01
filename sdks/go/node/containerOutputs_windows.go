package node

import (
	"fmt"
	"os"
)

func openContainerOutputs(path string) (*os.File, error) {
	return nil, fmt.Errorf("container outputs files require a Unix node")
}

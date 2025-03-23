// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Define the interfaces to mock and their destinations
	mocks := []struct {
		src  string
		dest string
	}{
		{
			src:  "github.com/fr0g-66723067/cc/internal/ai.Provider",
			dest: "./internal/ai/mocks/mock_provider.go",
		},
		{
			src:  "github.com/fr0g-66723067/cc/internal/container.Provider",
			dest: "./internal/container/mocks/mock_provider.go",
		},
		{
			src:  "github.com/fr0g-66723067/cc/internal/vcs.Provider",
			dest: "./internal/vcs/mocks/mock_provider.go",
		},
		{
			src:  "github.com/fr0g-66723067/cc/pkg/plugin.Plugin",
			dest: "./pkg/plugin/mocks/mock_plugin.go",
		},
	}

	// Run mockgen for each interface
	for _, m := range mocks {
		fmt.Printf("Generating mock for %s...\n", m.src)
		cmd := exec.Command("mockgen", "-package", "mocks", "-destination", m.dest, m.src)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to generate mock for %s: %v", m.src, err)
		}
	}

	fmt.Println("Mock generation complete!")
}
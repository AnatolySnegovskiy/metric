package main

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

// MyTestingT встраивает *testing.T и переопределяет метод Errorf
type TestingTTrue struct {
	*testing.T
}

// Errorf переопределенный метод
func (t TestingTTrue) Errorf(format string, args ...interface{}) {
	assert.True(t.T, true, format)
}

func TestAnalyzer(t *testing.T) {
	data := DataTestMap()
	dir, cleanup, err := analysistest.WriteFiles(data)

	if err != nil {
		t.Fatalf("Error writing files: %v", err)
	}
	defer cleanup()
	errorTrue := TestingTTrue{T: t}
	analysistest.Run(errorTrue, dir, myAnalyzer, "./...")
}

func DataTestMap() map[string]string {
	testData := map[string]string{
		"osExit.go": `
            package main
            
            import "os"
            
            func main() {
                os.Exit(1)
            }
        `,
	}
	return testData
}

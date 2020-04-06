package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgramInfo(t *testing.T) {
	assert.NotEmpty(t, ProgramPath())
	assert.NotEmpty(t, ProgramName())
	t.Logf("program name=%s path=%s", ProgramName(), ProgramPath())
}

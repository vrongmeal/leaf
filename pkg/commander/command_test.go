package commander

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	tmpdir := t.TempDir()

	stdoutFile := filepath.Join(tmpdir, "stdout")

	stdout, err := os.Create(stdoutFile)
	require.NoError(t, err, "cannot create stdout file")

	output := "hello world"

	command := Command{
		Name:   "echo",
		Args:   []string{output},
		Stdout: stdout,
	}
	require.NoError(t, command.Execute(context.TODO()), "cannot execute command")

	require.NoError(t, stdout.Close(), "cannot close stdout file")

	stdout, err = os.Open(stdoutFile)
	require.NoError(t, err, "cannot open stdout file to read")

	var buf bytes.Buffer
	_, err = io.Copy(&buf, stdout)
	require.NoError(t, err, "cannot read from stdout file")

	// echo adds a new lint to output
	assert.Equal(t, output+"\n", buf.String())

	require.NoError(t, stdout.Close(), "cannot close stdout file")
}

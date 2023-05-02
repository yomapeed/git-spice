package git

import (
	"context"
	"fmt"
	"io"
)

// Type specifies the type of a Git object.
type Type string

const (
	BlobType   Type = "blob"
	CommitType Type = "commit"
	TreeType   Type = "tree"
)

func (t Type) String() string {
	return string(t)
}

// ReadObject reads the object with the given hash from the repository
// into the given writer.
func (r *Repository) ReadObject(ctx context.Context, typ Type, hash Hash, dst io.Writer) error {
	cmd := r.gitCmd(ctx, "cat-file", string(typ), hash.String()).Stdout(dst)
	if err := cmd.Run(r.exec); err != nil {
		return fmt.Errorf("cat-file: %w", err)
	}
	return nil
}

// WriteObject writes an object of the given type to the repository,
// and returns the hash of the written object.
func (r *Repository) WriteObject(ctx context.Context, typ Type, src io.Reader) (Hash, error) {
	cmd := r.gitCmd(ctx, "hash-object", "-w", "--stdin", "-t", string(typ)).Stdin(src)
	out, err := cmd.OutputString(r.exec)
	if err != nil {
		return ZeroHash, fmt.Errorf("hash-object: %w", err)
	}
	return Hash(out), nil
}

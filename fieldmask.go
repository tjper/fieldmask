package fieldmask

import (
	"bytes"
	"reflect"
	"regexp"

	"github.com/pkg/errors"
	"google.golang.org/genproto/protobuf/field_mask"
)

// pathRE checks if a path is in the desired format.
var pathRE = regexp.MustCompile(`^(\w+){1}(\.\w+)?$`)

type fieldMasker interface {
	GetMask() *field_mask.FieldMask
}

// ToStruct uses reflection to map the values of a fieldMasker onto a struct
// with the same field tags. An error is returned if the fieldMasker has an
// invalid set of paths. Each path is expected to consist of fields seperated
// periods.
// E.G:
//    valid:   [user.name, status]
//    invalid: [user/name, userName]
// For more details on field masks read the following godoc:
// https://godoc.org/google.golang.org/genproto/protobuf/field_mask
func ToStruct(fm fieldMasker, st interface{}, opts ...PathOption) error {
	for _, path := range fm.GetMask().GetPaths() {
		for _, o := range opts {
			path = o(path)
		}
		if err := validate(pathRE, path); err != nil {
			return errors.Wrapf(err, "failed to ToStruct\tfm = %v\tst = %v", fm, st)
		}
		var (
			b   = []byte(path)
			src = reflect.ValueOf(fm).Elem()
			dst = reflect.ValueOf(st).Elem()
		)
		for bytes.ContainsRune(b, '.') {
			b = bytes.SplitN(b, []byte("."), 2)[0]
			src = src.FieldByName(string(b)).Elem()
			dst = dst.FieldByName(string(b)).Elem()
		}
		dst.FieldByName(string(b)).Set(dst.FieldByName(string(b)))
	}
	return nil
}

// PathOption is a function that manipulates a path string.
type PathOption func(string) string

// WithAliases maps path strings to an alias path string.
func WithAliases(aliases map[string]string) PathOption {
	return func(path string) string {
		alias, ok := aliases[path]
		if !ok {
			return path
		}
		return alias
	}
}

// validate ensures the path is valid. WARNING: failure to validate a path
// before usage could result in unexpected endpoint behavior.
func validate(re *regexp.Regexp, path string) error {
	if re.MatchString(path) {
		return errors.Errorf("failed to Validate fieldMasker\tinvalid characters found in path\tpath = %s", path)
	}
	if path == "id" {
		return errors.New("failed to Validate fieldMasker\tid field may not be set via field mask")
	}
	return nil
}

package confetti

import (
	"cmp"
	"errors"
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

// envLoader loads config from environment variables.
// If string is not empty it is used as a prefix for the environment variable.
type envLoader struct {
	prefix    string
	separator string
}

const DefaultSeparator = ","

func (e envLoader) Load(config any, ownConfig *confetti) (err error) {
	var errOnUnknown bool

	if ownConfig != nil {
		errOnUnknown = ownConfig.errOnUnknown
	}

	return loadEnv(config, e.prefix, e.separator, errOnUnknown)
}

// loadEnv recursively sets struct fields from env vars for arbitrarily deep nesting.
// If prefix is not empty, it is used as a prefix for the environment variable.
func loadEnv(config any, prefix, separator string, errOnUnknown bool) error {
	var unknowns map[string]struct{}

	if prefix != "" && errOnUnknown {
		unknowns = map[string]struct{}{}

		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 && strings.HasPrefix(parts[0], strings.ToUpper(prefix)+"_") {
				unknowns[parts[0]] = struct{}{}
			}
		}
	}

	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("config must be pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	if prefix != "" {
		prefix = strings.ToUpper(prefix)
	}

	for i := range t.NumField() {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// Check for struct tag override.
		tagEnv := field.Tag.Get("env")
		name := cmp.Or(tagEnv, camelToUpperSnake(field.Name))

		envName := name
		if prefix != "" && tagEnv == "" {
			envName = prefix + "_" + name
		}

		if fieldVal.Kind() == reflect.Struct {
			subPrefix := name
			if prefix != "" && tagEnv == "" {
				subPrefix = prefix + "_" + name
			}

			if err := loadEnv(fieldVal.Addr().Interface(), subPrefix, separator, errOnUnknown); err != nil {
				return err
			}

			continue
		}

		val, ok := os.LookupEnv(envName)
		if !ok {
			continue
		}

		delete(unknowns, envName)

		switch fieldVal.Kind() { //nolint:exhaustive // ok
		case reflect.String:
			fieldVal.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldVal.Type().PkgPath() == "time" && fieldVal.Type().Name() == "Duration" {
				d, err := time.ParseDuration(val)
				if err != nil {
					return fmt.Errorf("env %s: %w", envName, err)
				}

				fieldVal.SetInt(int64(d))

				break
			}

			iv, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf("env %s: %w", envName, err)
			}

			fieldVal.SetInt(iv)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uv, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return fmt.Errorf("env %s: %w", envName, err)
			}

			fieldVal.SetUint(uv)
		case reflect.Bool:
			bv, err := parseBool(val)
			if err != nil {
				return fmt.Errorf("env %s: %w", envName, err)
			}

			fieldVal.SetBool(bv)
		case reflect.Float32, reflect.Float64:
			fv, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return fmt.Errorf("env %s: %w", envName, err)
			}

			fieldVal.SetFloat(fv)
		case reflect.Slice:
			elemKind := fieldVal.Type().Elem().Kind()
			parts := strings.Split(val, separator)
			slice := reflect.MakeSlice(fieldVal.Type(), len(parts), len(parts))

			for j, part := range parts {
				part = strings.TrimSpace(part)

				switch elemKind {
				case reflect.String:
					slice.Index(j).SetString(part)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					iv, err := strconv.ParseInt(part, 10, 64)
					if err != nil {
						return fmt.Errorf("env %s[%d]: %w", envName, j, err)
					}

					slice.Index(j).SetInt(iv)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					uv, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						return fmt.Errorf("env %s[%d]: %w", envName, j, err)
					}

					slice.Index(j).SetUint(uv)
				case reflect.Float32, reflect.Float64:
					fv, err := strconv.ParseFloat(part, 64)
					if err != nil {
						return fmt.Errorf("env %s[%d]: %w", envName, j, err)
					}

					slice.Index(j).SetFloat(fv)
				case reflect.Bool:
					bv, err := parseBool(part)
					if err != nil {
						return fmt.Errorf("env %s[%d]: %w", envName, j, err)
					}

					slice.Index(j).SetBool(bv)
				default:
					return fmt.Errorf("env %s: unsupported slice element type %s", envName, elemKind)
				}
			}

			fieldVal.Set(slice)
		}
	}

	if len(unknowns) > 0 {
		unk := slices.Collect(maps.Keys(unknowns))
		slices.Sort(unk)

		return fmt.Errorf("unknown environment variables: %v", unk)
	}

	return nil
}

// parseBool parses a string into a boolean value, accepting
// the following as true: "1", "t", "true", "y", "yes" (case-insensitive)
// and as false: "0", "f", "false", "n", "no". Returns an error for anything else.
func parseBool(s string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "t", "true", "y", "yes":
		return true, nil
	case "0", "f", "false", "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %q", s)
	}
}

// camelToUpperSnake converts CamelCase to UPPER_SNAKE_CASE,
// with support for acronyms and abbreviations.
func camelToUpperSnake(s string) string {
	runes, out := []rune(s), []rune{}

	for i := range runes {
		if i > 0 && isUpper(runes[i]) && (isLower(runes[i-1]) || (i+1 < len(runes) && isLower(runes[i+1]))) {
			out = append(out, '_')
		}

		out = append(out, runes[i])
	}

	return strings.ToUpper(string(out))
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

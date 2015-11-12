package cfg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// c represents the configuration store, with a map to store the loaded keys
// from the environment.
var c struct {
	m map[string]string
}

// Init is to be called only once, to load up the giving namespace if found,
// in the environment variables. All keys will be made lowercase.
func Init(namespace string) {
	c.m = make(map[string]string)

	// Get the lists of available environment variables.
	envs := os.Environ()
	if len(envs) == 0 {
		panic("No environment variables found")
	}

	// Create the uppercase version to meet the standard {NAMESPACE_} format.
	uspace := fmt.Sprintf("%s_", strings.ToUpper(namespace))

	// Loop and match each variable using the uppercase namespace.
	for _, val := range envs {
		if !strings.HasPrefix(val, uspace) {
			continue
		}

		part := strings.Split(val, "=")
		c.m[strings.ToLower(strings.TrimPrefix(part[0], uspace))] = part[1]
	}

	// Did we find any keys for this namespace?
	if len(c.m) == 0 {
		panic(fmt.Sprintf("Namespace %q was not found", namespace))
	}
}

// String returns the value of the giving key as a string, else it will return
// a non-nil error if key was not found
func String(key string) (string, error) {
	value, found := c.m[key]
	if !found {
		return "", fmt.Errorf("Unknown key %s !", key)
	}

	return value, nil
}

// MustString returns the value of the giving key as a string, else it will panic
// if the key was not found.
func MustString(key string) string {
	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown key %s !", key))
	}

	return value
}

// Int returns the value of the giving key as an int, else it will return a
// non-nil error, if the key was not found or the value can't be convered to an int.
func Int(key string) (int, error) {
	value, found := c.m[key]
	if !found {
		return 0, fmt.Errorf("Unknown Key %s !", key)
	}

	iv, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return iv, nil
}

// MustInt returns the value of the giving key as an int, else it will panic
// if the key was not found or the value can't be convered to an int.
func MustInt(key string) int {
	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}

	iv, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Key %q value is not an int", key))
	}

	return iv
}

// Time returns the value of the giving key as a Time, else it will return a non-nil
// error, if the key was not found or the value can't be convered to a Time.
func Time(key string) (time.Time, error) {
	value, found := c.m[key]
	if !found {
		return time.Time{}, fmt.Errorf("Unknown Key %s !", key)
	}

	tv, err := time.Parse(time.UnixDate, value)
	if err != nil {
		return tv, err
	}

	return tv, nil
}

// MustTime returns the value of the giving key as a Time, else it will panic
// if the key was not found or the value can't be convered to a Time.
func MustTime(key string) time.Time {
	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}

	tv, err := time.Parse(time.UnixDate, value)
	if err != nil {
		panic(fmt.Sprintf("Key %q value is not a Time", key))
	}

	return tv
}

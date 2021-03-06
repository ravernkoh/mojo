package mojo

import (
	"strings"
)

// Assemble assembles the given objects back into arguments.
func (objs Objects) Assemble() ([]string, error) {
	return assemble(objs)
}

// assemble assembles arguments from the given objects.
func assemble(objs []Object) ([]string, error) {
	var args []string

	for len(objs) > 0 {
		switch objs[0].(type) {
		case CommandObject:
			args = append(args, objs[0].(CommandObject).Name)
			objs = objs[1:]
		case ArgumentObject:
			args = append(args, objs[0].(ArgumentObject).Value)
			objs = objs[1:]
		case FlagObject:
			newArgs, n, err := assembleFlag(objs)
			if err != nil {
				return nil, err
			}

			// Append the new arguments and remove n objects.
			args = append(args, newArgs...)
			objs = objs[n:]
		}
	}

	return args, nil
}

// assembleFlag assembles a flag from the given objects and returns the
// arguments, along with how many objects were used.
//
// Panics if the first object provided is not a flag.
func assembleFlag(objs []Object) ([]string, int, error) {
	var (
		args []string
		n    = 1
	)

	// Panic if the first object is not a flag.
	obj := objs[0].(FlagObject)
	objs = objs[1:]

	// Extract the possible name.
	var name strings.Builder
	name.WriteString(obj.Name)

	// If the flag is a multiple flag, then add all the names together.
	if obj.MultipleFlagsStart {
		for {
			// If there are no more objects, this means that there
			// was no end flag. Return error.
			if len(objs) == 0 {
				return nil, 0, FlagError{
					Err: ErrIncompleteMultipleFlag,
				}
			}

			// If the next object is not a flag, this means that
			// there was no end flag. Return error.
			flagObj, ok := objs[0].(FlagObject)
			if !ok {
				return nil, 0, FlagError{
					Err: ErrIncompleteMultipleFlag,
				}
			}

			// Append the name to the name builder after removing
			// the dash.
			name.WriteString(flagObj.Name[1:])
			objs = objs[1:]
			n++

			// Check whether it is the end flag.
			if flagObj.MultipleFlagsEnd {
				obj = flagObj
				break
			}
		}
	}

	// If the flag isn't a bool flag and it is a combined flag, append the
	// value to the name.
	if !obj.Bool && obj.CombinedFlagValues {
		name.WriteString("=" + obj.Value)
	}

	// Append the name to the arguments.
	args = append(args, name.String())

	// If the flag isn't a bool flag and also isn't a combined flag, append
	// the value to the arguments.
	if !obj.Bool && !obj.CombinedFlagValues {
		args = append(args, obj.Value)
	}

	return args, n, nil
}

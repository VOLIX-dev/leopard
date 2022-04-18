package helpers

import "strings"

func SerializeStack(stack []byte) map[string]interface{} {
	stackTrace := string(stack)

	stackTrace = strings.Replace(stackTrace, "\t", "    ", -1)
	lines := strings.Split(stackTrace, "\n")

	goRoutine := lines[0]
	lines = lines[1 : len(lines)-1]

	return map[string]interface{}{
		"thread": goRoutine,
		"stack":  lines,
	}
}

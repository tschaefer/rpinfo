/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package vcgencmd

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

type Exec interface {
	Run(args ...string) (map[string]string, error)
}

type Cmd struct{}

func (r Cmd) Run(args ...string) (map[string]string, error) {
	execCommand := exec.Command("vcgencmd", args...)
	out, err := execCommand.Output()
	if err != nil {
		if out != nil {
			err = fmt.Errorf("vcgencmd error: %s - %v", strings.TrimSpace(string(out)), err)
		} else {
			err = fmt.Errorf("vcgencmd error: %v", err)
		}
		return nil, err
	}

	output := strings.TrimSpace(string(out))

	outputMap := make(map[string]string)
	for line := range strings.SplitSeq(output, "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			slog.Warn(fmt.Sprintf("vcgencmd warn: skipping data: %s", line))
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		outputMap[key] = value
	}

	return outputMap, nil
}

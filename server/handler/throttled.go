/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package handler

import (
	"fmt"
	"strconv"
	"strings"
)

// StatusFlag represents a single status flag
type throttledStatusFlag struct {
	Bit  uint
	Desc string
}

var throttledStatusFlags = []throttledStatusFlag{
	{0, "Undervoltage detected"},
	{1, "Arm frequency capped"},
	{2, "Currently throttled"},
	{3, "Soft temperature limit active"},
	{16, "Undervoltage has occurred"},
	{17, "Arm frequency capping has occurred"},
	{18, "Throttling has occurred"},
	{19, "Soft temperature limit has occurred"},
}

func parseThrottledHex(hexStr string) ([]string, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	val, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %v", err)
	}

	var results []string
	for _, flag := range throttledStatusFlags {
		if val&(1<<flag.Bit) != 0 {
			results = append(results, flag.Desc)
		}
	}
	return results, nil
}

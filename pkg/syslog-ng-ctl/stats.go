// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Stats(cc ControlChannel) (stats []Stat, errs error) {
	rsp, err := cc.SendCommand("STATS")
	if err != nil {
		return stats, err
	}
	rsp = strings.TrimRight(rsp, "\n") // remove trailing new line
	lines := strings.Split(rsp, "\n")
	// TODO: sanity check: match header line
	lines = lines[1:] // drop header line: SourceName;SourceId;SourceInstance;State;Type;Number
	for _, line := range lines {
		fields := strings.Split(line, ";")
		if len(fields) != 6 {
			errs = errors.Join(errs, InvalidStatLine(line))
			continue
		}
		if len(fields[3]) != 1 {
			errs = errors.Join(errs, InvalidStatLine(line))
			continue
		}
		num, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		stats = append(stats, Stat{
			SourceName:     fields[0],
			SourceID:       fields[1],
			SourceInstance: fields[2],
			SourceState:    SourceState(fields[3][0]),
			Type:           fields[4],
			Number:         num,
		})
	}
	return
}

type Stat struct {
	SourceName     string
	SourceID       string
	SourceInstance string
	SourceState    SourceState
	Type           string
	Number         uint64
}

type SourceState byte

const (
	SourceStateActive   SourceState = 'a'
	SourceStateDynamic  SourceState = 'd'
	SourceStateOrphaned SourceState = 'o'
)

type InvalidStatLine string

func (err InvalidStatLine) Error() string {
	return fmt.Sprintf("invalid stat line: %q", string(err))
}
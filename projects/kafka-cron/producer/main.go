package main

import (
	"bufio"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"strings"
)

type CronSpec struct {
	Minute uint64 `json:"minute,omitempty"`
	Hour   uint64 `json:"hour,omitempty"`
	Dom    uint64 `json:"dom,omitempty"`
	Month  uint64 `json:"month,omitempty"`
	Dow    uint64 `json:"dow,omitempty"`
}

type Job struct {
	Command string `json:"command,omitempty"`
	CronSpec
}

func parseCronFile(reader io.Reader) ([]Job, error) {
	var jobs []Job
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lastInd := strings.LastIndex(line, " ")
		if lastInd == -1 {
			return nil, fmt.Errorf("failed to parse \"%s\"", line)
		}
		cronExpr := line[:lastInd]
		command := line[lastInd+1:]
		sched, err := parser.Parse(cronExpr)
		if err != nil {
			return nil, err
		}
		specShed, ok := sched.(*cron.SpecSchedule)
		if ok == false {
			return nil, fmt.Errorf("type assertion failed")
		}
		job := Job{
			Command: command,
			CronSpec: CronSpec{
				specShed.Minute,
				specShed.Hour,
				specShed.Dom,
				specShed.Month,
				specShed.Dow,
			},
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func main() {

}

package mgmt

import (
	"encoding/json"
	"fmt"
)

const (
	defaultTimeout = 1800
)

func (c *NutanixMGMTClient) Wait(taskUUID string) (*TaskPollResultDTO, error) {
	return c.WaitTimeout(taskUUID, defaultTimeout)
}

func (c *NutanixMGMTClient) WaitTimeout(taskUUID string, timeout int) (*TaskPollResultDTO, error) {
	query := map[string][]string{}
	query["timeoutseconds"] = []string{fmt.Sprintf("%d", timeout)}
	respBytes, err := c.DoRequest("GET", fmt.Sprintf("%s/%s/%s", "tasks", taskUUID, "poll"), query, nil)
	if err != nil {
		return nil, err
	}
	taskPollResult := &TaskPollResultDTO{}
	err = json.Unmarshal(respBytes, taskPollResult)
	if err != nil {
		return nil, err
	}
	return taskPollResult, nil
}

package rest

import (
	"encoding/json"
	"fmt"

	"nutanix/client/api/mgmt"
)

const (
	defaultTimeout = 1800
)

func (c *NutanixRESTClient) Wait(taskUUID string) (*mgmt.TaskPollResultDTO, error) {
	return c.WaitTimeout(taskUUID, defaultTimeout)
}

func (c *NutanixRESTClient) WaitTimeout(taskUUID string, timeout int) (*mgmt.TaskPollResultDTO, error) {
	query := map[string][]string{}
	query["timeoutseconds"] = []string{fmt.Sprintf("%d", timeout)}
	respBytes, err := c.DoRequest("GET", fmt.Sprintf("%s/%s/%s", "tasks", taskUUID, "poll"), query, nil)
	if err != nil {
		return nil, err
	}
	taskPollResult := &mgmt.TaskPollResultDTO{}
	err = json.Unmarshal(respBytes, taskPollResult)
	if err != nil {
		return nil, err
	}
	return taskPollResult, nil
}

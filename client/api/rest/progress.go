package rest

import (
	"encoding/json"
	"fmt"

	"nutanix/client/api/mgmt"
)

func (c *NutanixRESTClient) TrackProgress(taskId string) (*mgmt.ProgressList, error) {
	respBytes, err := c.DoRequest("GET", "/progress_monitors/", map[string][]string{
		"filterCriteria": []string{fmt.Sprintf("%s==%s", "entity_uuid_list", taskId)},
	}, nil)
	if err != nil {
		return nil, err
	}
	vmInfo := &mgmt.ProgressList{}
	err = json.Unmarshal(respBytes, vmInfo)
	if err != nil {
		return nil, err
	}
	return vmInfo, nil
}

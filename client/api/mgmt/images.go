package mgmt

import (
	"encoding/json"
)

func (c *NutanixMGMTClient) GetImageList() (*ImageList, error) {
	respBytes, err := c.DoRequest("GET", "/images/", nil, nil)
	if err != nil {
		return nil, err
	}
	imageList := &ImageList{}
	err = json.Unmarshal(respBytes, imageList)
	if err != nil {
		return nil, err
	}
	return imageList, nil
}

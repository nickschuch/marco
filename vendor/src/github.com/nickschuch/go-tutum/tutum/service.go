package tutum

import (
	"encoding/json"
)

type ServiceListResponse struct {
	Objects []ServiceOverview `json:"objects"`
}

type ServiceOverview struct {
	Uuid string `json:"uuid"`
}

type ServiceDetails struct {
	ContainerEnvVars []EnvVar `json:"container_envvars"`
	ContainerUris    []string `json:"containers"`  
}

func ListServices() ([]ServiceOverview, error) {
	url := "service/"
	request := "GET"

	var response ServiceListResponse
	data, err := TutumCall(url, request)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return response.Objects, nil
}

func GetService(uuid string) (ServiceDetails, error) {
	url := "service/" + uuid + "/"
	request := "GET"

	var response ServiceDetails
	data, err := TutumCall(url, request)
	if err != nil {
		return ServiceDetails{}, err
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return ServiceDetails{}, err
	}
	return response, nil
}

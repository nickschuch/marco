package tutum

import (
	"encoding/json"
)

type CListResponse struct {
	Objects []Container `json:"objects"`
}

type Container struct {
	Application          string     `json:"application"`
	Autodestroy          string     `json:"autodestroy"`
	Autoreplace          string     `json:"autoreplace"`
	Autorestart          string     `json:"autorestart"`
	ContainerPorts       []PortInfo `json:"container_ports"`
	ContainerSize        string     `json:"container_size"`
	CurrentNumContainers int        `json:"current_num_containers"`
	DeployedDatetime     string     `json:"deployed_datetime"`
	DestroyedDatetime    string     `json:"destroyed_datetime"`
	Entrypoint           string     `json:"entrypoint"`
	ExitCode             int        `json:"exit_code"`
	ExitCodeMessage      string     `json:"exit_code_message"`
	ImageName            string     `json:"image_name"`
	ImageTag             string     `json:"image_tag"`
	Name                 string     `json:"name"`
	PublicDNS            string     `json:"public_dns"`
	ResourceUri          string     `json:"resource_uri"`
	RunCommand           string     `json:"run_command"`
	StartedDatetime      string     `json:"started_datetime"`
	State                string     `json:"state"`
	StoppedDatetime      string     `json:"stopped_datetime"`
	UniqueName           string     `json:"unique_name"`
	Uuid                 string     `json:"uuid"`
}

func ListContainers() ([]Container, error) {
	url := "container/"
	request := "GET"

	var response CListResponse
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

func GetContainer(uuid string) (Container, error) {
	url := "container/" + uuid + "/"
	request := "GET"

	var response Container
	data, err := TutumCall(url, request)
	if err != nil {
		return Container{}, err
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return Container{}, err
	}
	return response, nil
}

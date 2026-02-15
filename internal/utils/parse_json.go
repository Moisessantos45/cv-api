package utils

import (
	"cv_api/internal/models"
	"encoding/json"
)

func ParseJSONToStringArray(data []byte) ([]models.Project, error) {
	var projects []models.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, err
	}

	for i := range projects {
		p := &projects[i]

		if len(p.Tecnologies) > 0 {
			var techs []string
			if err := json.Unmarshal(p.Tecnologies, &techs); err == nil {
				p.Tecnologies = json.RawMessage{}
			}
		}

		if len(p.Links) > 0 {
			var tempLink models.Link
			if err := json.Unmarshal(p.Links, &tempLink); err == nil {
			}
		}
	}

	return projects, nil
}

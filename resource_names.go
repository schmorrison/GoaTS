package GoaTS

import (
	"fmt"

	"github.com/goadesign/goa/design"
)

type resourceType struct {
	Name        string
	Path        string
	Media       string
	URI         string
	Payloads    []dataType
	Description string
}

func (g *Generator) genResourceNames(api *design.APIDefinition) ([]resourceType, error) {
	if api == nil {
		return nil, fmt.Errorf("GoaTS: design.Design is not initialized")
	}

	r := []resourceType{}

	err := api.IterateResources(func(res *design.ResourceDefinition) error {
		nr := resourceType{}
		nr.Name = res.Name
		nr.Path = res.BasePath
		nr.Media = res.MediaType
		nr.Description = res.Description
		nr.URI = res.URITemplate()

		ut := res.UserTypes()
		for _, a := range ut {
		Prime:
			for _, b := range nr.Payloads {
				if b.Name == a.TypeName {
					continue Prime
				}
			}
			tn, err := g.genTNames(a)
			if err != nil {
				return err
			}
			nr.Payloads = append(nr.Payloads, tn)
		}

		r = append(r, nr)

		return nil
	})
	if err != nil {
		return []resourceType{}, err
	}
	return r, nil
}

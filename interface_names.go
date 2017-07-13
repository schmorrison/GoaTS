package GoaTS

import (
	"github.com/goadesign/goa/design"
)

//datatype struct defines the base structure for each payload
type dataType struct {
	Name       string //Name of the payload
	MediaType  string //outgoing Mediatype
	FieldCount string //number of fields
	Type       string //Type name of the payload

	Fields map[string]fieldType //map of field names and data
}

//fieldtype struct defines the base structure for a field
type fieldType struct {
	Description string      //field description
	Type        string      //data type
	Required    bool        //required validation
	Nested      bool        //if value is object
	Value       interface{} //either a string of the datatype name or a map[string] fieldtype
}

//genTypeNames creates fills the dataType structure for each payload
func (g *Generator) genTNames(a *design.UserTypeDefinition) (dataType, error) {
	n := dataType{}
	//type name
	n.Name = a.TypeName
	n.Type = "UserType"
	//	fmt.Printf("Generating Interface for Usertype %s\n", a.TypeName)
	//get the fields for a payload
	i, err := g.getUserTypeFields(a)
	if err != nil {
		return dataType{}, err
	}
	//set the fields field in datatype
	n.Fields = i
	return n, nil
}

// func (g *Generator) genTypeNames(api *design.APIDefinition) ([]dataType, error) {
// 	if api == nil {
// 		return nil, fmt.Errorf("GoaTS: design.Design is not initialized")
// 	}

// 	inf := []dataType{}
// 	err := api.IterateResources(func(res *design.ResourceDefinition) error {
// 		fmt.Printf("Getting Usertypes from Resource %s\n", res.Name)
// 		ut := res.UserTypes()
// 	Prime:
// 		for _, a := range ut {
// 			for _, b := range inf {
// 				if b.Name == a.TypeName {
// 					continue Prime
// 				}
// 			}
// 			n := dataType{}
// 			n.Name = a.TypeName
// 			n.Type = "UserType"
// 			fmt.Printf("Generating Interface for Usertype %s\n", a.TypeName)
// 			i, err := g.getUserTypeFields(a)
// 			if err != nil {
// 				return err
// 			}
// 			n.Fields = i
// 			inf = append(inf, n)
// 			fmt.Println()
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return []dataType{}, err
// 	}

// 	fmt.Printf("Getting Mediatypes for API %s\n", api.Name)
// 	mt := api.MediaTypes
// 	for _, a := range mt {
// 		n := dataType{}
// 		n.Name = a.TypeName
// 		n.Type = "MediaType"
// 		fmt.Printf("Generating Interface for Mediatype %s\n", a.TypeName)
// 		i, err := g.getMediaTypeFields(a)
// 		if err != nil {
// 			return nil, err
// 		}
// 		n.Fields = i
// 		inf = append(inf, n)
// 		fmt.Println()
// 	}

// 	return inf, nil
// }

func (g *Generator) getUserTypeFields(a *design.UserTypeDefinition) (map[string]fieldType, error) {
	fields := make(map[string]fieldType)
	o := a.Type.ToObject()
	for j, b := range o {
		//		fmt.Printf("%s: %s\n", j, b.Type.Name())
		newField := fieldType{}
		if b.Type.IsObject() {
			f, err := g.getFieldTypes(b)
			if err != nil {
				return make(map[string]fieldType), err
			}
			newField.Nested = true
			newField.Value = f
		} else {
			newField.Value = b.Type.Name()
		}
		newField.Description = b.Description
		fields[j] = newField
	}
	return fields, nil
}

func (g *Generator) getMediaTypeFields(a *design.MediaTypeDefinition) (map[string]fieldType, error) {
	fields := make(map[string]fieldType)
	o := a.Type.ToObject()
	for j, b := range o {
		//		fmt.Printf("%s: %s\n", j, b.Type.Name())
		newField := fieldType{}
		if b.Type.IsObject() {
			f, err := g.getFieldTypes(b)
			if err != nil {
				return make(map[string]fieldType), err
			}
			newField.Nested = true
			newField.Value = f
		} else {
			newField.Value = b.Type.Name()
		}
		newField.Description = b.Description
		newField.Required = b.IsRequired(j)
		fields[j] = newField
	}
	return fields, nil
}

func (g *Generator) getFieldTypes(a *design.AttributeDefinition) (map[string]fieldType, error) {
	fields := make(map[string]fieldType)
	o := a.Type.ToObject()
	for j, b := range o {
		//		fmt.Printf("  %s: %s\n", j, b.Type.Name())
		newField := fieldType{}
		if b.Type.IsObject() {
			f, err := g.getFieldTypes(b)
			if err != nil {
				return make(map[string]fieldType), err
			}
			newField.Nested = true
			newField.Value = f
		} else {
			newField.Value = b.Type.Name()
		}
		newField.Description = b.Description
		newField.Required = b.IsRequired(j)
		fields[j] = newField
	}
	return fields, nil
}

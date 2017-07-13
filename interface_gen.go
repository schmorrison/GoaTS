package GoaTS

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

//template.FuncMap for html/template functions
var tFuncs = template.FuncMap{
	"writeValue": writeTypeValue,
}

type output struct {
	name  string
	parts []string
}

//writeTypes creates a file of the payloads for a given resource and returns the name
//of the new file
func (g *Generator) writeTypes(resource string, d []dataType) ([]string, error) {
	if len(d) == 0 {
		return []string{}, nil
	}
	u := output{}
	m := output{}

	//range over the payloads
	for _, a := range d {
		//if the payload has no fields then skip
		if len(a.Fields) == 0 {
			continue
		}

		//the template for a typescript interface
		str := `
interface {{ .Name }} {
{{ range $k, $v := .Fields }}    {{ $k }}: {{ $v | writeValue }};
{{ end }}
}
`
		//parse the template from the string
		t, err := template.New("").Funcs(tFuncs).Parse(str)
		if err != nil {
			log.Fatal("GoaTS: Failed to parse the template")
		}
		//create new buffer
		buf := bytes.NewBuffer([]byte{})
		//execute the template
		err = t.Execute(buf, a)
		if err != nil {
			log.Fatal("GoaTS: Failed to parse the template")
		}

		//if the type is user
		if a.Type == "UserType" {
			//write buffer to string and append to usertype list
			u.parts = append(u.parts, buf.String())
		}
		//if the type is media
		if a.Type == "MediaType" {
			//write buffer to string and append to mediatype list
			m.parts = append(m.parts, buf.String())
		}
	}

	outputFile := filepath.Join(g.outDir, fmt.Sprintf("%s.ts", resource))
	content := strings.Join(u.parts, "\n")
	//write the file to the name '{{resource}}.ts'
	if err := ioutil.WriteFile(outputFile, []byte(content), 0755); err != nil {
		return []string{}, err
	}

	return []string{outputFile}, nil
}

//writeTypeValue is template.Func
//takes the fields value and returns the absolute representation of it
func writeTypeValue(a fieldType) string {
	//return strings as is
	if b, ok := a.Value.(string); ok {
		return b
	}
	//if type is map[string]fieldType then nested object exists
	if b, ok := a.Value.(map[string]fieldType); ok {
		//send the map into writeNestedType
		return writeNestedType(b, "    ")
	}
	log.Fatal("GoaTS: Did not have a string or map[string]interface for value.")
	return ""
}

//writeNestedType iterates over each entry in the map
//and writes the value to the string
func writeNestedType(a map[string]fieldType, indent string) string {
	//start the string
	str := "{\n"
	//iterate over the map of fields
	for k, v := range a {
		//if the value is string
		if b, ok := v.Value.(string); ok {
			//write to the string
			str += fmt.Sprintf("%s%s%s: %s;\n", indent, indent, k, b)
		}
		//if the field is map[string]fieldType
		//execute this function within the string
		if b, ok := v.Value.(map[string]fieldType); ok {
			str += fmt.Sprintf("%s%s%s: %s;\n", indent, indent, k, writeNestedType(b, indent+"    "))
		}
	}
	//finish the string
	str += indent + "}"
	return str
}

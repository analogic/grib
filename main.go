package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/nilsmagnus/grib/grib"
)

func main() {
	//category := flag.Int("category", 0, "Category. Default is temperature") // temperature
	//product := flag.Int("product", 6, "Product. Default is temperature") // temperature
	filename := flag.String("file", "", "Grib filename")
	exportType := flag.Int("export", 0, "Export format. Valid types are 0 (none) 1 (json) ")
	maxNum := flag.Int("maxmsg", math.MaxInt32, "Maximum number of messages to parse.")

	flag.Parse()

	gribFile, err := os.Open(*filename)

	if err != nil {
		fmt.Printf("\nFile [%s] not found.\n", *filename)
	}
	defer gribFile.Close()

	messages, err := grib.ReadMessages(gribFile, *maxNum)

	if err != nil {
		fmt.Printf("Error reading all messages in gribfile: %s", err.Error())
	}

	switch *exportType {
	case 0:
	case 1:
		exportJSONConsole(messages)
	}

}

func exportJSONConsole(messages []grib.Message) {
	for _, message := range messages {
		export(&message)
	}
}

func export(m *grib.Message) {
	templateNumber := int(m.Section4.ProductDefinitionTemplateNumber)
	template := m.Section4.ProductDefinitionTemplate
	category := int(template.ParameterCategory)
	number := int(template.ParameterNumber)

	d := make(map[string]interface{})

	d["type"] = grib.ReadDataType(int(m.Section1.Type))
	d["template"] = grib.ReadProductDefinitionTemplateNumber(templateNumber)
	d["category"] = grib.ReadProductDisciplineParameters(templateNumber, category)
	d["parameter"] = grib.ReadProductDisciplineCategoryParameters(templateNumber, category, number)
	d["grid"] = grib.ReadGridDefinitionTemplateNumber(int(m.Section3.TemplateNumber))
	d["surface1"] = grib.ReadSurfaceTypesUnits(int(m.Section4.ProductDefinitionTemplate.FirstSurface.Type))
	d["surface1value"] = m.Section4.ProductDefinitionTemplate.FirstSurface.Value
	d["surface1scale"] = m.Section4.ProductDefinitionTemplate.FirstSurface.Scale
	d["surface2"] = grib.ReadSurfaceTypesUnits(int(m.Section4.ProductDefinitionTemplate.SecondSurface.Type))
	d["surface2value"] = m.Section4.ProductDefinitionTemplate.SecondSurface.Value
	d["data"] = m.Section7.Data

	for k, v := range m.Section3.Definition.Export() {
		d[k] = v
	}

	// json print
	js, _ := json.Marshal(d)
	var out bytes.Buffer
	json.Indent(&out, js, "", "\t")
	out.WriteTo(os.Stdout)
	fmt.Println("")
}

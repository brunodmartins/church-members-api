package file

type Data struct {
	Value interface{}
}
type DataMapper func(Data) []string

func TransformToCSVData(data []Data, columns []string, dataMapper DataMapper) [][]string {
	var finalData [][]string
	finalData = append(finalData, columns)

	for _, line := range data {
		finalData = append(finalData, dataMapper(line))
	}

	return finalData
}

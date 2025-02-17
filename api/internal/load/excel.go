package load

import (
	"errors"

	"github.com/xuri/excelize/v2"
)

func ExcelImport(file string) ([][]string, error) {
	var rows [][]string
	f, err := excelize.OpenFile(file)
	if err != nil {
		return rows, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return rows, errors.New("any sheet not found")
	}
	firstSheet := sheets[0]

	// Tüm satırları oku
	rows, err = f.GetRows(firstSheet)
	if err != nil {
		return rows, err
	}

	return rows, nil

	/* for _, row := range rows {
		err := config.App().DB.ExistsInTable("towns", map[string]any{
			"id":          row[0],
			"district_id": row[1],
		})
		if err == nil {
			config.App().DB.DynamicUpdate(load.Param{
				Table: "towns",
				Fields: map[string]any{
					"latitude":  row[2],
					"longitude": row[3],
				},
				Conditions: map[string]any{
					"id": row[0],
				},
			})
		}
	} */
}

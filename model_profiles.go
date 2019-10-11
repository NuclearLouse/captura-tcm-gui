package main

import (
	"fmt"
	l "log"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"

	"redits.oculeus.com/asorokin/tcm/structs"
)

type modelProfiles struct {
	quantityRows int
	checkStates  []int
	cellValue    [][]string
}

func newModelProfiles() *modelProfiles {
	m := new(modelProfiles)
	rows, values, err := dataTableProfiles()
	if err != nil {
		l.Println("Error reading from database. Err=", err)
	}
	m.quantityRows = rows - 1
	m.checkStates = make([]int, m.quantityRows)
	m.cellValue = values
	return m
}

func dataTableProfiles() (int, [][]string, error) {
	var profiles []structs.ItestProfiles
	if err := pg.Find(&profiles).Error; err != nil {
		return 0, nil, err
	}
	rows := len(profiles) + 1
	cellValue := make([][]string, rows)
	cellValue[1] = make([]string, rows)
	cellValue[2] = make([]string, rows)
	cellValue[3] = make([]string, rows)
	for i := range profiles {
		cellValue[1][i] = profiles[i].ProfileID
		cellValue[2][i] = profiles[i].ProfileName
		cellValue[3][i] = profiles[i].ProfileIP
	}
	return rows, cellValue, nil
}

func (mp *modelProfiles) ButtAddProfile() {
	for i := 0; i < mp.quantityRows; i++ {
		if mp.checkStates[i] == 1 {
			fmt.Printf("Added row %d. Profile=%s. ID=%s\n", i+1, mp.cellValue[2][i], mp.cellValue[1][i])
			profile := fmt.Sprintf("Profile: %s", mp.cellValue[2][i])
			switch mp.cellValue[2][i] {
			case "AMVTS":
				newTest.SystemName = "amvts"
			case "BMVTS":
				newTest.SystemName = "bmvts"
			case "Avys_S2":
				newTest.SystemName = "fmvts"
			}
			newTest.ProfileID = mp.cellValue[1][i]
			textRequest = fmt.Sprintf("Request: %s?t=%d&profid=%s&%s=%s&ndbccgid=%s&ndbcgid=%s",
				itestAPI.ApiURL, apiRequest, newTest.ProfileID, venPref, newTest.SupOrPref, newTest.CountryID, newTest.BreakoutID)
			entryProfile.SetText(profile)
			entryRequest.SetText(textRequest)
			return
		}
	}

}

func (mp *modelProfiles) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableInt(0), // column 4 checkbox state
	}
}
func (mp *modelProfiles) NumRows(m *ui.TableModel) int {
	return mp.quantityRows
}

func (mp *modelProfiles) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	if column == 0 {
		return ui.TableString(fmt.Sprintf("%d", row+1))
	}
	if column == 4 {
		return ui.TableInt(mp.checkStates[row])
	}
	return ui.TableString(mp.cellValue[column][row])
}

func (mp *modelProfiles) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	if column == 4 { // checkboxes
		mp.checkStates[row] = int(value.(ui.TableInt))
	}
}

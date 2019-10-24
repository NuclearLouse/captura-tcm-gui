package main

import (
	"fmt"
	l "log"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
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
	var profiles []itestProfiles
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
			l.Printf("Select Profile:%s. Profile_ID:%s\n", mp.cellValue[2][i], mp.cellValue[1][i])
			entry.Profile = fmt.Sprintf("Profile: %s", mp.cellValue[2][i])
			switch mp.cellValue[2][i] {
			case "AMVTS":
				newtest.SystemName = "amvts"
			case "BMVTS":
				newtest.SystemName = "bmvts"
			case "Avys_S2":
				newtest.SystemName = "fmvts"
			}
			newtest.ProfileID = mp.cellValue[1][i]
			entry.Request = fmt.Sprintf("%s?t=%d&profid=%s&%s=%s&ndbccgid=%s&ndbcgid=%s",
				itest.URL, apiRequest, newtest.ProfileID, venPref, newtest.SupOrPref, newtest.CountryID, newtest.BreakoutID)
			entryProfile.SetText(entry.Profile)
			entryRequest.SetText(entry.Request)
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

package main

import (
	"fmt"
	l "log"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"

	"redits.oculeus.com/asorokin/tcm/structs"
)

type modelSuppliers struct {
	quantityRows int
	checkStates  []int
	cellValue    [][]string
}

func newModelSuppliers() *modelSuppliers {
	m := new(modelSuppliers)
	rows, values, err := dataTableSuppliers()
	if err != nil {
		l.Println("Error reading from database. Err=", err)
	}
	m.quantityRows = rows - 1
	m.checkStates = make([]int, m.quantityRows)
	m.cellValue = values
	return m
}

func dataTableSuppliers() (int, [][]string, error) {
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

func (ms *modelSuppliers) ButtAddSupplier() {
	for i := 0; i < ms.quantityRows; i++ {
		if ms.checkStates[i] == 1 {
			fmt.Printf("Added row %d. Profile=%s. ID=%s\n", i+1, ms.cellValue[2][i], ms.cellValue[1][i])
			profile := fmt.Sprintf("Profile: %s", ms.cellValue[2][i])
			switch ms.cellValue[2][i] {
			case "AMVTS":
				newTest.SystemName = "amvts"
			case "BMVTS":
				newTest.SystemName = "bmvts"
			case "Avys_S2":
				newTest.SystemName = "fmvts"
			}
			entryProfile.SetText(profile)
			return
		}
	}

}

func (ms *modelSuppliers) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableInt(0), // column 3 checkbox state
	}
}
func (ms *modelSuppliers) NumRows(m *ui.TableModel) int {
	return ms.quantityRows
}

func (ms *modelSuppliers) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	if column == 0 {
		return ui.TableString(fmt.Sprintf("%d", row+1))
	}
	if column == 4 {
		return ui.TableInt(ms.checkStates[row])
	}
	return ui.TableString(ms.cellValue[column][row])
}

func (ms *modelSuppliers) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	if column == 4 { // checkboxes
		ms.checkStates[row] = int(value.(ui.TableInt))
	}
}

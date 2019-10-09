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
	var suppliers []structs.ItestSuppliers
	if err := pg.Find(&suppliers).Error; err != nil {
		return 0, nil, err
	}
	rows := len(suppliers) + 1
	cellValue := make([][]string, rows)
	cellValue[1] = make([]string, rows)
	cellValue[2] = make([]string, rows)
	cellValue[3] = make([]string, rows)
	for i := range suppliers {
		cellValue[1][i] = suppliers[i].SupplierID
		cellValue[2][i] = suppliers[i].SupplierName
		cellValue[3][i] = suppliers[i].Prefix
	}
	return rows, cellValue, nil
}

func (ms *modelSuppliers) ButtAddSupplier() {
	for i := 0; i < ms.quantityRows; i++ {
		if ms.checkStates[i] == 1 {
			fmt.Printf("Added row %d. Prefix=%s. ID=%s\n", i+1, ms.cellValue[3][i], ms.cellValue[1][i])
			var supplier string
			switch newTest.CallType {
			case "CLI":
				newTest.SupplierID = ms.cellValue[1][i]
				supplier = fmt.Sprintf("SupplierID: %s", ms.cellValue[1][i])
			case "Voice":
				newTest.Prefix = ms.cellValue[2][i]
				supplier = fmt.Sprintf("Prefix: %s", ms.cellValue[3][i])
			}
			entrySupplier.SetText(supplier)
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

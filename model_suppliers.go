package main

import (
	"fmt"
	l "log"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
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
	var suppliers []itestSuppliers
	searchTemplate := fmt.Sprintf("%s%%", newtest.SystemName)
	if err := pg.Where("supplier_name LIKE ?", searchTemplate).Find(&suppliers).Error; err != nil {
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
			l.Printf("Select supplier:%s. Supplier_ID:%s .Prefix:%s\n", ms.cellValue[2][i], ms.cellValue[1][i], ms.cellValue[3][i])
			switch newtest.CallType {
			case "CLI":
				newtest.SupOrPref = ms.cellValue[1][i]
				entry.Supplier = fmt.Sprintf("SupplierID: %s", ms.cellValue[1][i])
			case "Voice":
				newtest.SupOrPref = ms.cellValue[3][i]
				entry.Supplier = fmt.Sprintf("Prefix: %s", ms.cellValue[3][i])

			}
			entry.Request = fmt.Sprintf("%s?t=%d&profid=%s&%s=%s&ndbccgid=%s&ndbcgid=%s",
				itest.URL, apiRequest, newtest.ProfileID, venPref, newtest.SupOrPref, newtest.CountryID, newtest.BreakoutID)
			entrySupplier.SetText(entry.Supplier)
			entryRequest.SetText(entry.Request)
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
		ui.TableInt(0), // column 4 checkbox state
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

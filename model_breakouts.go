package main

import (
	"fmt"
	l "log"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"

	"redits.oculeus.com/asorokin/tcm/structs"
)

var breakCLI structs.ItestBreakoutsCli
var breakSTD structs.ItestBreakoutsStd

type modelBreakouts struct {
	quantityRows int
	checkStates  []int
	cellValue    [][]string
}

func newModelBreakouts() *modelBreakouts {
	m := new(modelBreakouts)
	rows, values, err := dataTableBreakouts()
	if err != nil {
		l.Println("Error reading from database. Err=", err)
	}
	m.quantityRows = rows - 1
	m.checkStates = make([]int, m.quantityRows)
	m.cellValue = values
	return m
}

func dataTableBreakouts() (int, [][]string, error) {
	var rows int
	var cellValue [][]string
	var err error
	switch newTest.CallType {
	case "CLI":
		var breakouts []structs.ItestBreakoutsCli
		if err = pg.Find(&breakouts).Error; err != nil {
			return 0, nil, err
		}
		rows = len(breakouts) + 1
		cellValue := make([][]string, rows)
		cellValue[1] = make([]string, rows)
		cellValue[2] = make([]string, rows)
		cellValue[3] = make([]string, rows)
		cellValue[4] = make([]string, rows)
		for i := range breakouts {
			cellValue[1][i] = breakouts[i].CountryID
			cellValue[2][i] = breakouts[i].CountryName
			cellValue[3][i] = breakouts[i].BreakoutName
			cellValue[4][i] = breakouts[i].BreakoutID
		}
		return rows, cellValue, nil
	case "Voice":
		var breakouts []structs.ItestBreakoutsStd
		if err = pg.Find(&breakouts).Error; err != nil {
			return 0, nil, err
		}
		rows = len(breakouts) + 1
		cellValue := make([][]string, rows)
		cellValue[1] = make([]string, rows)
		cellValue[2] = make([]string, rows)
		cellValue[3] = make([]string, rows)
		cellValue[4] = make([]string, rows)
		for i := range breakouts {
			cellValue[1][i] = breakouts[i].CountryID
			cellValue[2][i] = breakouts[i].CountryName
			cellValue[3][i] = breakouts[i].BreakoutName
			cellValue[4][i] = breakouts[i].BreakoutID
		}
		return rows, cellValue, nil
	}
	return rows, cellValue, nil
}

func (mb *modelBreakouts) ButtAddBreakout() {
	for i := 0; i < mb.quantityRows; i++ {
		if mb.checkStates[i] == 1 {
			fmt.Printf("Added row %d. CountryID=%s. Breakout=%s\n", i+1, mb.cellValue[1][i], mb.cellValue[4][i])
			newTest.CountryID = mb.cellValue[1][i]
			newTest.BreakoutID = mb.cellValue[4][i]
			country := fmt.Sprintf("Country: %s", mb.cellValue[2][i])
			breakout := fmt.Sprintf("Breakout: %s", mb.cellValue[3][i])
			entryCountry.SetText(country)
			entryBreakout.SetText(breakout)
			return
		}
	}
}

func (mb *modelBreakouts) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableInt(0), // column 5 checkbox state
	}
}
func (mb *modelBreakouts) NumRows(m *ui.TableModel) int {
	return mb.quantityRows
}

func (mb *modelBreakouts) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	if column == 0 {
		return ui.TableString(fmt.Sprintf("%d", row+1))
	}
	if column == 5 {
		return ui.TableInt(mb.checkStates[row])
	}
	return ui.TableString(mb.cellValue[column][row])
}

func (mb *modelBreakouts) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	if column == 5 { // checkboxes
		mb.checkStates[row] = int(value.(ui.TableInt))
	}
}

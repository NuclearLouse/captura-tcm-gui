package main

import (
	"fmt"
	l "log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"

	"github.com/jinzhu/gorm"
	"redits.oculeus.com/asorokin/tcm/config"
	"redits.oculeus.com/asorokin/tcm/database"
	"redits.oculeus.com/asorokin/tcm/structs"
)

var (
	itestAPI                                                            structs.ItestApi
	mainwin                                                             *ui.Window
	pg                                                                  *gorm.DB
	slash                                                               string
	absPath                                                             string
	entryType, entryProfile, entrySupplier, entryCountry, entryBreakout *ui.Entry
)

const (
	settingsFile = "tcm.ini"
)

func init() {
	var err error
	switch runtime.GOOS {
	case "linux":
		slash = "/"
	case "windows":
		slash = "\\"
	default:
		l.Fatalf("%s not support operation system", runtime.GOOS)
	}
	ex, err := os.Executable()
	if err != nil {
		l.Fatalln("FATAL! Cann't get the absolute path of the executive file. Error=", err)
	}

	absPath = filepath.Dir(ex)

	cfg, err := config.ReadConfig(settingsFile)
	if err != nil {
		l.Fatalln("FATAL! Failed to load config. Error=", err)
	}
	db, err := database.NewDB(cfg, absPath+slash)
	if err != nil {
		l.Fatal("FATAL! Could not connect to the database. Error=", err)
	}
	pg = db.Connect
	if err := apiSettings(); err != nil {
		l.Fatal("FATAL! Cann't obtained iTest API settings. Error=", err)
	}
}

func main() {
	ui.Main(setupUI)
}

func apiSettings() error {
	if err := pg.Take(&itestAPI).Error; err != nil {
		return err
	}
	return nil
}

func tSystems() (int, [][]string, error) {
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

type myModelHandler struct {
	quantityRows int
	checkStates  []int
	cellValue    [][]string
}

func newNodelHandler() *myModelHandler {
	m := new(myModelHandler)
	rows, systems, err := tSystems()
	if err != nil {
		l.Println("Error reading from database. Err=", err)
	}
	m.quantityRows = rows - 1
	m.checkStates = make([]int, m.quantityRows)
	m.cellValue = systems
	return m
}

func (mh *myModelHandler) Butt() {
	for i := 0; i < mh.quantityRows; i++ {
		if mh.checkStates[i] == 1 {
			fmt.Printf("Added row %d. Profile=%s. ID=%s\n", i+1, mh.cellValue[2][i], mh.cellValue[1][i])
			profile := fmt.Sprintf("Profile: %s", mh.cellValue[2][i])
			entryProfile.SetText(profile)
			return
		}
	}

}

func (mh *myModelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableString(""),
		ui.TableInt(0), // column 3 checkbox state
	}
}
func (mh *myModelHandler) NumRows(m *ui.TableModel) int {
	return mh.quantityRows
}

func (mh *myModelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	if column == 0 {
		return ui.TableString(fmt.Sprintf("%d", row+1))
	}
	if column == 4 {
		return ui.TableInt(mh.checkStates[row])
	}
	return ui.TableString(mh.cellValue[column][row])
}

func (mh *myModelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	if column == 4 { // checkboxes
		mh.checkStates[row] = int(value.(ui.TableInt))
	}
}

func setupUI() {
	mainwin = ui.NewWindow("Test Calls System Manage - Initiate New Tests", 780, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})
	// Создание основного бокса
	mainVbox := ui.NewVerticalBox()
	mainVbox.SetPadded(true)
	mainwin.SetChild(mainVbox)
	mainwin.SetMargined(true)

	// Создание бокса для размещения вкладок
	tabsHbox := ui.NewHorizontalBox()
	tabsHbox.SetPadded(true)
	mainVbox.Append(tabsHbox, true)

	// Размещение вкладок на боксе
	tabProfiles := ui.NewTab()
	tabProfiles.Append("Profiles", makeProfilesPage())
	tabProfiles.SetMargined(0, true)
	tabProfiles.Append("Suppliers", makeSuppliersPage())
	tabProfiles.SetMargined(1, true)
	tabProfiles.Append("Destinations", makeDestinationsPage())
	tabProfiles.SetMargined(2, true)
	tabProfiles.Append("Results", makeResultsPage())
	tabProfiles.SetMargined(3, true)
	tabsHbox.Append(tabProfiles, true)

	mainVbox.Append(ui.NewHorizontalSeparator(), false)

	// Создание бокса для информации о добавляемых значениях
	entrysHbox := ui.NewHorizontalBox()
	entrysHbox.SetPadded(true)
	mainVbox.Append(entrysHbox, false)

	// Создание и управление ячейками
	entryType = ui.NewEntry()
	entryType.SetReadOnly(true)
	entryType.SetText("Test Type:")
	entryProfile = ui.NewEntry()
	entryProfile.SetReadOnly(true)
	entryProfile.SetText("Profile:")
	entrySupplier = ui.NewEntry()
	entrySupplier.SetReadOnly(true)
	entrySupplier.SetText("Supplier or Prefix:")
	entryCountry = ui.NewEntry()
	entryCountry.SetReadOnly(true)
	entryCountry.SetText("CountryID:")
	entryBreakout = ui.NewEntry()
	entryBreakout.SetReadOnly(true)
	entryBreakout.SetText("Breakout:")

	// Размещение ячеек для добавленой информации
	entrysHbox.Append(entryType, true)
	entrysHbox.Append(entryProfile, true)
	entrysHbox.Append(entrySupplier, true)
	entrysHbox.Append(entryCountry, true)
	entrysHbox.Append(entryBreakout, true)

	// Кнопка старт тестов
	buttonStart := ui.NewButton("Start Test")
	buttonStart.OnClicked(func(*ui.Button) {
		textType := testType()
		if textType == "" {
			textType = "(cancelled)"
		}
		entryType.SetText(textType)
	})
	// entrysHbox.Append(entryForm, true)
	entrysHbox.Append(buttonStart, true)

	mainwin.Show()
}

func testType() string {
	return "Test Type: CLI"
}

func makeProfilesPage() ui.Control {
	// vbox := ui.NewVerticalBox()
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	group := ui.NewGroup("Test Type")
	group.SetMargined(true)
	hbox.Append(group, false)
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	rb := ui.NewRadioButtons()
	rb.Append("Test CLI")
	rb.Append("Test Voice")
	vbox.Append(rb, false)

	rb.OnSelected(func(*ui.RadioButtons) {
		var typeTest string
		switch rb.Selected() {
		case 0:
			typeTest = "CLI"
		case 1:
			typeTest = "Voice"
		}
		text := fmt.Sprintf("Test Type: %s", typeTest)
		entryType.SetText(text)
	})

	group = ui.NewGroup("Profiles")
	group.SetMargined(true)
	hbox.Append(group, true)

	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	mh := newNodelHandler()
	model := ui.NewTableModel(mh)
	table := ui.NewTable(&ui.TableParams{
		Model: model,
	})

	table.AppendTextColumn("№", 0, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Profile ID", 1, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Profile Name", 2, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Profile IP", 3, ui.TableModelColumnNeverEditable, nil)
	table.AppendCheckboxColumn("Select", 4, ui.TableModelColumnAlwaysEditable)

	button := ui.NewButton("Add Profile")
	vbox.Append(table, true)
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	button.OnClicked(func(*ui.Button) {
		mh.Butt()
	})

	grid.Append(button,
		1, 0, 1, 1, //left, top int, xspan, yspan int
		false, ui.AlignCenter, false, ui.AlignCenter)

	return hbox
}
func makeSuppliersPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	return hbox
}
func makeDestinationsPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	return hbox
}
func makeResultsPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	return hbox
}

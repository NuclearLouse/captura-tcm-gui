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
	newTest                                                             structs.NewInitTest
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
	defer fmt.Println(newTest)
}

func apiSettings() error {
	if err := pg.Take(&itestAPI).Error; err != nil {
		return err
	}
	return nil
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
	// Создание вкладки типа теста и профилей
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
		newTest.CallType = typeTest
		entryType.SetText(text)
	})

	group = ui.NewGroup("Profiles")
	group.SetMargined(true)
	hbox.Append(group, true)

	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	mp := newModelProfiles()
	model := ui.NewTableModel(mp)
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
		mp.ButtAddProfile()
	})

	grid.Append(button,
		1, 0, 1, 1, //left, top int, xspan, yspan int
		false, ui.AlignCenter, false, ui.AlignCenter)

	return hbox
}
func makeSuppliersPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	group := ui.NewGroup("Suppliers")
	group.SetMargined(true)
	hbox.Append(group, true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	ms := newModelSuppliers()
	model := ui.NewTableModel(ms)
	table := ui.NewTable(&ui.TableParams{
		Model: model,
	})

	table.AppendTextColumn("№", 0, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Supplier ID", 1, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Supplier Name", 2, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Prefix", 3, ui.TableModelColumnNeverEditable, nil)
	table.AppendCheckboxColumn("Select", 4, ui.TableModelColumnAlwaysEditable)

	button := ui.NewButton("Add Supplier")
	vbox.Append(table, true)
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	button.OnClicked(func(*ui.Button) {
		ms.ButtAddSupplier()
	})

	grid.Append(button,
		1, 0, 1, 1, //left, top int, xspan, yspan int
		false, ui.AlignCenter, false, ui.AlignCenter)

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

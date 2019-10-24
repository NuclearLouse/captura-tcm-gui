package main

import (
	"fmt"
	l "log"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/jinzhu/gorm"
)

var (
	itest                            itestAPI
	newtest                          newInitTest
	mainwin                          *ui.Window
	pg                               *gorm.DB
	apiRequest                       int
	slash, absPath, venPref, prefSup string
	entryType, entryProfile, entrySupplier,
	entryCountry, entryBreakout, entryRequest *ui.Entry
	entry entrys
)

const (
	settingsFile = "tcm.ini"
)

type entrys struct {
	TestType string
	Profile  string
	Supplier string
	Country  string
	Breakout string
	Request  string
}

func newEntrys() entrys {
	return entrys{
		TestType: "Test Type:",
		Profile:  "Profile:",
		Supplier: "Supplier or Prefix:",
		Country:  "Country:",
		Breakout: "Breakout:",
		Request:  itest.URL,
	}
}

func newTest() newInitTest {
	return newInitTest{}
}

func init() {
	var err error
	cfg, err := readConfig(settingsFile)
	if err != nil {
		l.Fatalln("FATAL! Failed to load config. Error=", err)
	}

	db, err := newDB(cfg)
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
	defer l.Println("Exit program")
}

func apiSettings() error {
	if err := pg.Take(&itest).Error; err != nil {
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

	// Создание бокса для размещения вкладки
	tabsHbox := ui.NewHorizontalBox()
	tabsHbox.SetPadded(true)
	mainVbox.Append(tabsHbox, true)

	// Размещение вкладки на боксе
	tabProfiles := ui.NewTab()
	tabProfiles.Append("Profiles", makeProfilesPage())
	tabProfiles.SetMargined(0, true)
	tabsHbox.Append(tabProfiles, true)

	mainVbox.Append(ui.NewHorizontalSeparator(), false)

	// Создание бокса для информации о добавляемых значениях
	entrysHbox := ui.NewHorizontalBox()
	entrysHbox.SetPadded(true)
	mainVbox.Append(entrysHbox, false)

	// Создание и управление ячейками
	entry = newEntrys()
	entryType = ui.NewEntry()
	entryType.SetReadOnly(true)
	entryType.SetText(entry.TestType)
	entryProfile = ui.NewEntry()
	entryProfile.SetReadOnly(true)
	entryProfile.SetText(entry.Profile)
	entrySupplier = ui.NewEntry()
	entrySupplier.SetReadOnly(true)
	entrySupplier.SetText(entry.Supplier)
	entryCountry = ui.NewEntry()
	entryCountry.SetReadOnly(true)
	entryCountry.SetText(entry.Country)
	entryBreakout = ui.NewEntry()
	entryBreakout.SetReadOnly(true)
	entryBreakout.SetText(entry.Breakout)

	// Размещение ячеек для добавленой информации
	entrysHbox.Append(entryType, true)
	entrysHbox.Append(entryProfile, true)
	entrysHbox.Append(entrySupplier, true)
	entrysHbox.Append(entryCountry, true)
	entrysHbox.Append(entryBreakout, true)

	requestHbox := ui.NewHorizontalBox()
	requestHbox.SetPadded(true)
	mainVbox.Append(requestHbox, false)

	entryRequest = ui.NewEntry()
	entryRequest.SetReadOnly(true)
	entryRequest.SetText(entry.Request)
	requestHbox.Append(entryRequest, true)
	// Кнопка старт тестов
	buttonStart := ui.NewButton("Start Test")
	buttonStart.OnClicked(func(*ui.Button) {
		startTest()

	})
	// entrysHbox.Append(entryForm, true)
	requestHbox.Append(buttonStart, false)

	mainwin.Show()
}

func makeProfilesPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	butSup := ui.NewButton("Select Suppliers")
	butSup.OnClicked(func(*ui.Button) {
		makeSuppliersPage()
	})
	butDes := ui.NewButton("Select Destination")
	butDes.OnClicked(func(*ui.Button) {
		makeDestinationsPage()
	})

	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)
	grid.Append(butSup,
		1, 0, 1, 1, //left, top int, xspan, yspan int
		false, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(butDes,
		2, 0, 1, 1, //left, top int, xspan, yspan int
		false, ui.AlignCenter, false, ui.AlignCenter)

	// Создание вкладки типа теста и профилей
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	vbox.Append(hbox, true)

	group := ui.NewGroup("Test Type")
	group.SetMargined(true)
	hbox.Append(group, false)
	vboxg := ui.NewVerticalBox()
	vboxg.SetPadded(true)
	group.SetChild(vboxg)

	rb := ui.NewRadioButtons()
	rb.Append("Test CLI")
	rb.Append("Test Voice")
	vboxg.Append(rb, false)

	rb.OnSelected(func(*ui.RadioButtons) {
		newtest = newTest()
		entry = newEntrys()
		var typeTest string
		switch rb.Selected() {
		case 0:
			typeTest = "CLI"
			apiRequest = itest.TestInitCli
			venPref = "vended"
		case 1:
			typeTest = "Voice"
			apiRequest = itest.TestInit
			venPref = "prefix"
		}
		newtest.CallType = typeTest
		entry.TestType = "Test Type: " + typeTest
		entry.Request = fmt.Sprintf("%s?t=%d&profid=%s&%s=%s&ndbccgid=%s&ndbcgid=%s",
			itest.URL, apiRequest, newtest.ProfileID, venPref, newtest.SupOrPref, newtest.CountryID, newtest.BreakoutID)
		entryType.SetText(entry.TestType)
		entryRequest.SetText(entry.Request)
		l.Println("Select Type test:", typeTest)
	})

	group = ui.NewGroup("Profiles")
	group.SetMargined(true)
	hbox.Append(group, true)

	vboxt := ui.NewVerticalBox()
	vboxt.SetPadded(true)
	group.SetChild(vboxt)

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
	vboxt.Append(table, true)
	grid = ui.NewGrid()
	grid.SetPadded(true)
	vboxt.Append(grid, false)

	button.OnClicked(func(*ui.Button) {
		mp.ButtAddProfile()
	})

	grid.Append(button,
		1, 0, 1, 1, //left, top int, xspan, yspan int
		false, ui.AlignCenter, false, ui.AlignCenter)

	return vbox
}

func makeSuppliersPage() {
	if newtest.SystemName != "" {
		win := ui.NewWindow("Suppliers", 780, 480, true)
		win.OnClosing(func(*ui.Window) bool {
			win.Destroy()
			return false
		})
		ui.OnShouldQuit(func() bool {
			win.Destroy()
			return false
		})

		vbox := ui.NewVerticalBox()
		vbox.SetPadded(true)
		win.SetChild(vbox)
		win.SetMargined(true)

		group := ui.NewGroup("Suppliers")
		group.SetMargined(true)
		vbox.Append(group, true)

		vboxt := ui.NewVerticalBox()
		vboxt.SetPadded(true)
		group.SetChild(vboxt)

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
		vboxt.Append(table, true)
		grid := ui.NewGrid()
		grid.SetPadded(true)
		vboxt.Append(grid, false)

		button.OnClicked(func(*ui.Button) {
			ms.ButtAddSupplier()
			win.Destroy()
		})

		grid.Append(button,
			1, 0, 1, 1, //left, top int, xspan, yspan int
			false, ui.AlignCenter, false, ui.AlignCenter)

		win.Show()
	}
	// иначе надо показывать окно с предупреждением о выборе системы
}

func makeDestinationsPage() {
	// os.Unsetenv("TEST_TYPE")
	if newtest.CallType != "" && newtest.SystemName != "" {
		// os.Setenv("TEST_TYPE", newTest.CallType)
		win := ui.NewWindow("Destinations", 780, 480, true)
		win.OnClosing(func(*ui.Window) bool {
			win.Destroy()
			return false
		})
		ui.OnShouldQuit(func() bool {
			win.Destroy()
			return false
		})

		vbox := ui.NewVerticalBox()
		vbox.SetPadded(true)
		win.SetChild(vbox)
		win.SetMargined(true)

		group := ui.NewGroup("Breakouts")
		group.SetMargined(true)
		vbox.Append(group, true)

		vboxt := ui.NewVerticalBox()
		vboxt.SetPadded(true)
		group.SetChild(vboxt)

		ms := newModelBreakouts()
		model := ui.NewTableModel(ms)
		table := ui.NewTable(&ui.TableParams{
			Model: model,
		})
		table.AppendTextColumn("№", 0, ui.TableModelColumnNeverEditable, nil)
		table.AppendTextColumn("Country ID", 1, ui.TableModelColumnNeverEditable, nil)
		table.AppendTextColumn("Country Name", 2, ui.TableModelColumnNeverEditable, nil)
		table.AppendTextColumn("Breakout", 3, ui.TableModelColumnNeverEditable, nil)
		table.AppendTextColumn("Breakout ID", 4, ui.TableModelColumnNeverEditable, nil)
		table.AppendCheckboxColumn("Select", 5, ui.TableModelColumnAlwaysEditable)

		button := ui.NewButton("Add Destination")
		vboxt.Append(table, true)
		grid := ui.NewGrid()
		grid.SetPadded(true)
		vboxt.Append(grid, false)

		button.OnClicked(func(*ui.Button) {
			ms.ButtAddBreakout()
			win.Destroy()
		})

		grid.Append(button,
			1, 0, 1, 1, //left, top int, xspan, yspan int
			false, ui.AlignCenter, false, ui.AlignCenter)

		win.Show()
	}
	// иначе надо показывать окно с предупреждением о выборе типа теста
}

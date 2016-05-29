package main

import (
	"fmt"
	"github.com/s-rah/recoil"
	"github.com/visualfc/goqt/ui"
	"log"
	"os"
)

type Widget struct {
	*ui.QWidget
	Label *ui.QLabel
	text  string
}

type MainWindow struct {
	*ui.QMainWindow
	targetSectionLabel *ui.QLabel
	localSectionLabel  *ui.QLabel

	labelTarget   *ui.QLabel
	labelAction   *ui.QLabel
	labelName     *ui.QLabel
	labelKey      *ui.QLabel
	labelHostname *ui.QLabel
	labelMessage  *ui.QLabel
	labelLogs     *ui.QLabel

	editTarget   *ui.QLineEdit
	editAction   *ui.QComboBox
	editHostname *ui.QLineEdit
	editName     *ui.QLineEdit
	logsBox      *ui.QPlainTextEdit
	editMessage  *ui.QPlainTextEdit
	editKey      *ui.QPlainTextEdit

	sendButton      *ui.QPushButton
	clearLogsButton *ui.QPushButton

	currentMessageFile string
}

func main() {
	ui.RunEx(os.Args, main_ui)
}

func (w *Widget) SetText(text string) {
	w.text = text
}

func main_ui() {
	app := ui.Application()
	app.SetOrganizationName("?")
	app.SetApplicationName("Recoil")

	w := &MainWindow{}
	w.QMainWindow = ui.NewMainWindow()

	// Widgets
	w.targetSectionLabel = ui.NewLabel()
	w.targetSectionLabel.SetText("Target Test Settings")
	font := w.Font()
	font.SetPointSize(font.PointSize() + 1)
	w.targetSectionLabel.SetFont(font)

	w.labelTarget = ui.NewLabel()
	w.labelTarget.SetText(fmt.Sprintf("Target Address"))
	w.labelTarget.SetFixedWidth(400)
	w.editTarget = ui.NewLineEdit()
	w.editTarget.SetPlaceholderText(fmt.Sprintf("target.onion"))

	w.labelAction = ui.NewLabel()
	w.labelAction.SetText(fmt.Sprintf("Action"))
	w.editAction = ui.NewComboBox()
	w.editAction.AddItem(fmt.Sprintf("ping"))
	w.editAction.AddItem(fmt.Sprintf("contact-request"))
	w.editAction.AddItem(fmt.Sprintf("send-messages"))
	w.editAction.AddItem(fmt.Sprintf("spamchannel"))
	w.editAction.SetFixedWidth(400)

	w.localSectionLabel = ui.NewLabel()
	w.localSectionLabel.SetText(fmt.Sprintf("Local Test Settings"))
	w.localSectionLabel.SetFont(font)

	w.labelHostname = ui.NewLabel()
	w.labelHostname.SetText(fmt.Sprintf("Onion Hostname"))
	w.labelHostname.SetFixedWidth(400)
	w.editHostname = ui.NewLineEdit()
	w.editHostname.SetPlaceholderText(fmt.Sprintf("something.onion"))
	w.editHostname.SetFixedWidth(400)

	w.labelName = ui.NewLabel()
	w.labelName.SetText(fmt.Sprintf("Ricochet Name"))
	w.labelName.SetFixedWidth(400)
	w.editName = ui.NewLineEdit()
	w.editName.SetText("recoil")
	w.editName.SetPlaceholderText(fmt.Sprintf("Ricochet Name"))
	w.editName.SetFixedWidth(400)

	w.labelKey = ui.NewLabel()
	w.labelKey.SetText(fmt.Sprintf("Onion Private Key"))
	w.labelKey.SetFixedWidth(400)
	w.editKey = ui.NewPlainTextEdit()
	w.editKey.SetFixedWidth(400)
	w.editKey.SetFixedHeight(99)

	w.labelMessage = ui.NewLabel()
	w.labelMessage.SetText(fmt.Sprintf("Message"))
	w.editMessage = ui.NewPlainTextEdit()
	w.editMessage.SetFixedHeight(100)
	w.editMessage.AppendPlainText(fmt.Sprintf("I am the recoil testing tool"))

	w.sendButton = ui.NewPushButton()
	w.sendButton.SetText("Send")
	w.clearLogsButton = ui.NewPushButton()
	w.clearLogsButton.SetText("Clear Logs")

	w.sendButton.OnClicked(func() {
		// Need to add back in debug
		target := w.editTarget.DisplayText()
		//debug := false
		action := w.editAction.CurrentText()
		hostname := w.editHostname.DisplayText()
		privateKey := w.editKey.ToPlainText()
		name := w.editName.DisplayText()
		message := w.editMessage.ToPlainText()

		if target == "" {
			w.logsBox.AppendPlainText(fmt.Sprintf("[ERROR] Target must be specified."))
			w.StatusBar().ShowMessage("Error: Target must be specified.")
			return
		} else if hostname == "" {
			w.logsBox.AppendPlainText(fmt.Sprintf("[ERROR] Hostname must be specified."))
			w.StatusBar().ShowMessage("Error: Hostname must be specified.")
			return
		} else if privateKey == "" {
			w.logsBox.AppendPlainText(fmt.Sprintf("[ERROR] Onion key must be specified."))
			w.StatusBar().ShowMessage("Error: Onion key must be specified.")
			return
		}

		recoil := new(recoil.Recoil)
		recoil.Ready = make(chan bool)

		if action == "ping" {
			online := recoil.Ping(privateKey, hostname, target)
			if online == true {
				w.logsBox.AppendPlainText(fmt.Sprintf("[INFO] Target appears to be online."))
				w.StatusBar().ShowMessage(fmt.Sprintf("%s appers to be online.", target))
			} else {
				w.logsBox.AppendPlainText(fmt.Sprintf("[INFO] Target appears to be offline."))
				w.StatusBar().ShowMessage(fmt.Sprintf("%s appers to be offline.", target))
			}
		} else {
			go recoil.Authenticate(privateKey, hostname, target)
			log.Printf("Running Recoil...")
			ready := <-recoil.Ready
			log.Printf("Received Authentication Result %v", ready)
			if ready == true {
				if action == "contact-request" {
					recoil.SendContactRequest(name, message)
					w.logsBox.AppendPlainText(fmt.Sprintf("[INFO] Sent contact request to %s.", target))
					w.StatusBar().ShowMessage(fmt.Sprintf("Sent contact request to %s.", target))
				} else if action == "spamchannel" {
					// go recoil.SpamChannel()
					// w.logsBox.AppendPlainText("[ERROR] spamchannel action not functional.")
					// w.StatusBar().ShowMessage("Error: spamchannel action not functional.")
				}
			}
		}
	})
	w.labelLogs = ui.NewLabel()
	w.labelLogs.SetText(fmt.Sprintf("Logs"))
	w.logsBox = ui.NewPlainTextEdit()
	w.logsBox.SetFixedWidth(800)
	w.logsBox.AppendPlainText(fmt.Sprintf("[Recoil] Ricochet testing toolkit"))
	w.logsBox.SetReadOnly(true)

	targetRow := ui.NewHBoxLayout()
	targetSection := ui.NewVBoxLayout()
	targetSection.AddWidget(w.labelTarget)
	targetSection.AddWidget(w.editTarget)
	actionSection := ui.NewVBoxLayout()
	actionSection.AddWidget(w.labelAction)
	actionSection.AddWidget(w.editAction)
	targetRow.AddLayout(targetSection)
	targetRow.AddLayout(actionSection)

	localRow := ui.NewHBoxLayout()
	localCol1 := ui.NewVBoxLayout()
	localCol1.AddWidget(w.labelName)
	localCol1.AddWidget(w.editName)
	localCol1.AddWidget(w.labelHostname)
	localCol1.AddWidget(w.editHostname)
	localRow.AddLayout(localCol1)
	localCol2 := ui.NewVBoxLayout()
	localCol2.AddWidget(w.labelKey)
	localCol2.AddWidget(w.editKey)
	localRow.AddLayout(localCol2)

	buttonRow := ui.NewHBoxLayout()
	buttonRow.AddWidget(w.sendButton)
	buttonRow.AddWidget(w.clearLogsButton)

	inputLayout := ui.NewVBoxLayout()
	inputLayout.AddWidget(w.localSectionLabel)
	inputLayout.AddLayout(localRow)
	inputLayout.AddWidget(w.targetSectionLabel)
	inputLayout.AddLayout(targetRow)
	inputLayout.AddWidget(w.labelMessage)
	inputLayout.AddWidget(w.editMessage)
	inputLayout.AddWidget(w.labelLogs)
	inputLayout.AddWidget(w.logsBox)
	inputLayout.AddLayout(buttonRow)

	centralWidget := ui.NewWidget()
	centralWidget.SetLayout(inputLayout)

	w.SetCentralWidget(centralWidget)

	w.clearLogsButton.OnClicked(func() {
		w.logsBox.Clear()
	})

	//splash := ui.NewSplashScreen()
	//splash.ShowMessage("test")

	openMessageAct := ui.NewActionWithTextParent("&Open Message File...", w)
	openMessageAct.OnTriggered(func() { w.openMessage() })
	openOnionHostnameAct := ui.NewActionWithTextParent("&Open Onion Service Hostname...", w)
	openOnionHostnameAct.OnTriggered(func() { w.openOnionHostname() })
	openOnionKeyAct := ui.NewActionWithTextParent("&Open Onion Service Key...", w)
	openOnionKeyAct.OnTriggered(func() { w.openOnionKey() })
	saveMessageAct := ui.NewActionWithTextParent("&Save Message File...", w)
	saveMessageAct.OnTriggered(func() { w.saveMessage() })
	saveMessageAsAct := ui.NewActionWithTextParent("&Save Message File As...", w)
	saveMessageAsAct.OnTriggered(func() { w.saveMessageAs() })
	exitAct := ui.NewActionWithTextParent("&Exit", w)
	exitAct.OnTriggered(func() { os.Exit(0) })

	fileMenu := w.MenuBar().AddMenuWithTitle("&File")
	fileMenu.AddAction(openMessageAct)
	fileMenu.AddAction(openOnionHostnameAct)
	fileMenu.AddAction(openOnionKeyAct)
	fileMenu.AddSeparator()
	fileMenu.AddAction(saveMessageAct)
	fileMenu.AddAction(saveMessageAsAct)
	fileMenu.AddSeparator()
	fileMenu.AddAction(exitAct)

	w.setCurrentMessageFile("")
	w.StatusBar().ShowMessage("Ready")
	w.Show()
}

func (w *MainWindow) openMessage() {
	filename := ui.QFileDialogGetOpenFileName()
	if filename != "" {
		w.loadFile(filename, "message")
	}
}

func (w *MainWindow) openOnionKey() {
	filename := ui.QFileDialogGetOpenFileName()
	if filename != "" {
		w.loadFile(filename, "onion-key")
	}
}

func (w *MainWindow) openOnionHostname() {
	filename := ui.QFileDialogGetOpenFileName()
	if filename != "" {
		w.loadFile(filename, "onion-hostname")
	}
}

func (w *MainWindow) saveMessage() bool {
	if w.currentMessageFile == "" {
		return w.saveMessageAs()
	}
	return w.saveFile(w.currentMessageFile)
}

func (w *MainWindow) saveMessageAs() bool {
	fileName := ui.QFileDialogGetSaveFileName()
	if fileName == "" {
		return false
	}
	return w.saveFile(fileName)
}

func (w *MainWindow) saveFile(fileName string) bool {
	file := ui.NewFileWithName(fileName)
	defer file.Delete()
	if !file.Open(ui.QIODevice_WriteOnly | ui.QIODevice_Text) {
		return false
	}
	file.Write([]byte(w.editMessage.ToPlainText()))
	w.setCurrentMessageFile(fileName)
	w.StatusBar().ShowMessage("Message file saved.")
	return true
}

func (w *MainWindow) loadFile(filename, fileType string) {
	file := ui.NewFileWithName(filename)
	defer file.Delete()
	if !file.Open(ui.QIODevice_ReadOnly) {
		w.StatusBar().ShowMessage("Error: Failed to load message file.")
		return
	} else if fileType == "message" {
		w.editMessage.SetPlainText(string(file.ReadAll()))
		w.setCurrentMessageFile(filename)
		w.StatusBar().ShowMessage("Message loaded.")
	} else if fileType == "onion-hostname" {
		w.editHostname.SetText(string(file.ReadAll()))
		w.StatusBar().ShowMessage("Onion service hostname loaded.")
	} else if fileType == "onion-key" {
		w.editKey.SetPlainText(string(file.ReadAll()))
		w.StatusBar().ShowMessage("Onion service key loaded.")
	}
}

func (w *MainWindow) setCurrentMessageFile(filename string) {
	w.currentMessageFile = filename
	//w.SetWindowFilePath(filename)
}

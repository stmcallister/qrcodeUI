package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	myApp := app.New()
	formWindow := myApp.NewWindow("qrcodeUI")
	formWindow.Resize(fyne.NewSize(800, 800))
	formWindow.SetMaster()

	title := widget.NewLabelWithStyle("QR Code Generator", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	desc := widget.NewLabelWithStyle("Convert text to a QR code", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	strInput := widget.NewEntry()
	strInput.SetPlaceHolder("Text to convert to QR Code...")

	imgInput := widget.NewEntry()
	imgInput.SetPlaceHolder("Name of the QR Code Image...")

	saveDirPath := ""
	saveDirButton := widget.NewButton("Select Save Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil && err == nil {
				saveDirPath = uri.Path()
				fmt.Println("Selected directory:", saveDirPath)
			}
		}, formWindow)
	})
	tabContainer := container.NewAppTabs() // Define and initialize tabContainer here

	createButton := widget.NewButton("Create QR Code", func() {
		if len(strInput.Text) > 0 {
			imgName := imgInput.Text
			if imgName == "" {
				dialog.ShowInformation("Information", "Image name cannot be blank. \nPlease enter an image name.", formWindow)
				return
			}
			imgFileExt := imgName[len(imgName)-4:]
			if imgFileExt != ".png" {
				imgName = fmt.Sprintf("%s.png", imgName)
			}
			createQrCode(imgName, strInput.Text, saveDirPath)
			imagePath := fmt.Sprintf("%s/%s", saveDirPath, imgName)
			// imagePathLabel := widget.NewLabelWithStyle(imagePath, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

			// Button for opening qr code image file
			openFile := func() {
				var cmd *exec.Cmd
				switch runtime.GOOS {
				case "windows":
					cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", imagePath)
				case "darwin":
					cmd = exec.Command("open", imagePath)
				case "linux":
					cmd = exec.Command("xdg-open", imagePath)
				default:
					return
				}
				err := cmd.Start()
				if err != nil {
					println("Failed to open file:", err)
				}
			}
			qrCodeButton := widget.NewButton(imagePath, openFile)

			image := canvas.NewImageFromFile(imagePath)
			image.FillMode = canvas.ImageFillOriginal

			imageContent := container.NewVBox(
				image,
				qrCodeButton,
			)

			codeTab := container.NewTabItem("QR Code: "+strInput.Text, imageContent)
			tabContainer.Append(codeTab)
			tabContainer.Select(codeTab)

			log.Println("Content was:", strInput.Text)
		} else {
			dialog.ShowInformation("Information", "Text cannot be blank. \nPlease enter text to create QR code.", formWindow)
		}
	})

	inputTab := container.NewTabItem("Input", container.NewVBox(title, desc, strInput, imgInput, saveDirButton, createButton))

	tabContainer = container.NewAppTabs(inputTab)

	formWindow.SetContent(tabContainer)
	formWindow.ShowAndRun()
}

func createQrCode(imgFileName, qrString, imgDir string) {
	qrc, err := qrcode.New(qrString)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}

	imgPath := fmt.Sprintf("%s/%s", imgDir, imgFileName)

	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		if err := os.Mkdir(imgDir, os.ModePerm); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	w, err := standard.New(imgPath)
	if err != nil {
		fmt.Printf("standard.New failed: %v", err)
		return
	}

	if err = qrc.Save(w); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}

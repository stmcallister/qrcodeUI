package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("qrcodeUI")
	myWindow.Resize(fyne.NewSize(300, 100))

	strInput := widget.NewEntry()
	strInput.SetPlaceHolder("Enter text to convert to QR Code...")

	imgInput := widget.NewEntry()
	imgInput.SetPlaceHolder("Enter the name of the QR Code Image...")

	content := container.NewVBox(strInput, imgInput, widget.NewButton("Create QR Code", func() {
		if len(strInput.Text) == 0 {

		}
		imgName := imgInput.Text
		imgFileExt := imgName[len(imgName)-4:]
		if imgFileExt != ".png" {
			imgName = fmt.Sprintf("%s.png", imgName)
		}
		createQrCode(imgName, strInput.Text)
		image := canvas.NewImageFromFile("img/" + imgName)
		myWindow.SetContent(image)
		myWindow.Resize(fyne.NewSize(200, 200))
		log.Println("Content was:", strInput.Text)
	}))

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func createQrCode(imgFileName, qrString string) {
	qrc, err := qrcode.New(qrString)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}

	imgDir := "img"
	imgPath := fmt.Sprintf("%s/%s", imgDir, imgFileName)

	// check if folder exists
	if _, err := os.Stat(imgDir); err == nil {
		fmt.Printf("imgDir exists\n")
	} else {
		if err := os.Mkdir(imgDir, os.ModePerm); err != nil {
			fmt.Println(err.Error())
		}
		// create folder
		fmt.Printf("imgDir created\n")
	}

	// check if file exists
	if _, err := os.Stat(imgPath); err == nil {
		fmt.Printf("imgPath exists\n")
	} else {
		// write to file
		w, err := standard.New(imgPath)
		if err != nil {
			fmt.Printf("standard.New failed: %v", err)
			return
		}

		// save file
		if err = qrc.Save(w); err != nil {
			fmt.Printf("could not save image: %v", err)
		} else {
			_, err := os.Open(imgPath)
			if err != nil {
				fmt.Printf("could not open saved image: %v", err)
			}
			// defer f.Close()
		}
	}

}

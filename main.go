package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
    "github.com/kbinani/screenshot"
    "image"
    "image/png"
    "os"
    "time"
)

func captureScreen() (image.Image, error) {
    bounds := screenshot.GetDisplayBounds(0)

    img, err := screenshot.CaptureRect(bounds)
    if err != nil {
        return nil, err
    }
    return img, nil
}

func saveImage(img image.Image, filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    return png.Encode(file, img)
}

func captureAndSave(label *widget.Label, parent fyne.Window) {
    parent.Hide()

    time.Sleep(1 * time.Second)

    // Capture the screen based on the option
    img, err := captureScreen()
    if err != nil {
        label.SetText(fmt.Sprintf("Error: %v", err))
        parent.Show()
        return
    }

    defer parent.Show()

    dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
        if err != nil {
            label.SetText(fmt.Sprintf("Error: %v", err))
            return
        }
        if file == nil {
            label.SetText("Save operation canceled")
            return
        }

        defer file.Close()

        if err := png.Encode(file, img); err != nil {
            label.SetText(fmt.Sprintf("Error saving file: %v", err))
            return
        }

        label.SetText(fmt.Sprintf("Screenshot saved: %s", file.URI().Path()))
    }, parent)
}

func readBytes(file *os.File) []byte {
    stat, _ := file.Stat()
    bytes := make([]byte, stat.Size())
    file.Read(bytes)
    return bytes
}

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Screenshot App")

    iconFile, err := os.Open("icon.png")
    if err == nil {
        defer iconFile.Close()
        iconResource := fyne.NewStaticResource("icon.png", readBytes(iconFile))
        myApp.SetIcon(iconResource)
    } else {
        // Option 2: Use a built-in theme icon as a fallback
        myApp.SetIcon(theme.FyneLogo())
    }

    label := widget.NewLabel("Tale a screenshot")
    buttonFull := widget.NewButton("Capture", func() {
        captureAndSave(label, myWindow)
    })

    content := container.NewVBox(label, buttonFull)
    myWindow.SetContent(content)
    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()
}

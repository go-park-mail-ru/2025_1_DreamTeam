package sertificate

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func GenerateCertificate(name, course, dateStr string, outputFile string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Фон
	pdf.SetFillColor(212, 227, 250)
	pdf.Rect(0, 0, 297, 210, "F")

	// Рамка
	pdf.SetLineWidth(2)
	pdf.SetDrawColor(34, 63, 151)
	pdf.Rect(5, 5, 287, 200, "D")

	// Шрифт
	fontPath := "./pkg/sertificate/fonts/DejaVuSerif.ttf"
	pdf.AddUTF8Font("DejaVu", "", fontPath)

	// Логотип
	opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	imageInfo := pdf.RegisterImageOptions("./pkg/sertificate/logo1.png", opt)
	if imageInfo == nil {
		return fmt.Errorf("не удалось зарегистрировать изображение логотипа")
	}
	pdf.ImageOptions("./pkg/sertificate/logo1.png", 20, 20, 20, 20, false, opt, 0, "")

	pdf.SetTextColor(34, 63, 151)
	pdf.SetFont("DejaVu", "", 28)
	pdf.SetXY(40, 20)
	pdf.LinkString(40, 20, 60, 10, "http://217.16.21.64/")
	pdf.CellFormat(0, 10, "SkillForce", "", 1, "L", false, 0, "")

	pdf.SetFont("DejaVu", "", 14)
	pdf.SetXY(40, 30)
	pdf.CellFormat(0, 10, "Маркетплейс курсов", "", 0, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	// Заголовок с декоративными элементами
	pdf.SetFont("DejaVu", "", 32)
	pdf.SetXY(0, 50)
	pdf.CellFormat(297, 20, "~ СЕРТИФИКАТ ~", "", 1, "C", false, 0, "")

	pdf.SetFont("DejaVu", "", 20)
	pdf.SetXY(0, 65)
	pdf.CellFormat(0, 20, "Настоящим подтверждается, что", "", 0, "C", false, 0, "")

	pdf.SetFont("DejaVu", "", 27)
	pdf.SetTextColor(34, 63, 151)
	pdf.SetXY(0, 80)
	pdf.CellFormat(0, 20, name, "", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("DejaVu", "", 20)

	pdf.SetXY(0, 95)
	pdf.CellFormat(0, 20, "успешно завершил(а) курс", "", 0, "C", false, 0, "")

	pdf.SetFont("DejaVu", "", 27)
	pdf.SetTextColor(34, 63, 151)
	pdf.SetXY(0, 110)
	pdf.CellFormat(0, 20, course, "", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	pdf.SetFont("DejaVu", "", 18)
	pdf.SetXY(55, -50)
	pdf.CellFormat(0, 10, fmt.Sprintf("Дата: %s   Подпись:", dateStr), "", 0, "L", false, 0, "")

	imageInfo = pdf.RegisterImageOptions("./pkg/sertificate/signature.png", opt)
	if imageInfo == nil {
		return fmt.Errorf("не удалось зарегистрировать изображение подписи")
	}
	pdf.ImageOptions("./pkg/sertificate/signature.png", 182, 137, 80, 30, false, opt, 0, "")

	pdf.SetXY(160, -50)
	pdf.CellFormat(0, 10, "______________________________________", "", 1, "C", false, 0, "")

	pdf.SetFont("DejaVu", "", 13)
	pdf.SetXY(160, -42)
	pdf.CellFormat(0, 10, "Д.А. Санталов - Генеральный директор ", "", 1, "C", false, 0, "")

	return pdf.OutputFileAndClose(outputFile)
}

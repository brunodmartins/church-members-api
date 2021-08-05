package file

import (
	"bufio"
	"bytes"
	"fmt"
	i18n2 "github.com/BrunoDM2943/church-members-api/platform/i18n"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/signintech/gopdf"
	"github.com/spf13/viper"
)

//Builder interface
//go:generate mockgen -source=./pdf.go -destination=./mock/pdf_mock.go
type Builder interface {
	BuildFile(title string, data []*domain.Member) ([]byte, error)
}

type pdfBuilder struct {
	messageService *i18n2.MessageService
}

func NewPDFBuilder() *pdfBuilder {
	return &pdfBuilder{
		messageService: i18n2.GetMessageService(),
	}
}

func (pdfBuilder *pdfBuilder) buildFirstPageSection(title string, builder *gopdf.GoPdf) {
	builder.AddPage()
	builder.SetMargins(30, 30, 30, 30)
	rect := &gopdf.Rect{H: gopdf.PageSizeA4.H, W: gopdf.PageSizeA4.W}
	options := gopdf.CellOption{
		Align: gopdf.Center,
	}

	builder.SetX(0)
	builder.SetY(100)
	builder.CellWithOption(rect, viper.GetString("church.name"), options)
	builder.SetY(120)
	builder.SetX(0)
	builder.CellWithOption(rect, title, options)
}

func (pdfBuilder *pdfBuilder) setFont(builder *gopdf.GoPdf) error {
	builder.AddTTFFont("arial", viper.GetString("pdf.font.path"))

	err := builder.SetFont("arial", "", 14)
	if err != nil {
		return err
	}
	return nil
}

func (pdfBuilder *pdfBuilder) toBytes(builder *gopdf.GoPdf) []byte {
	byteArr := &bytes.Buffer{}
	buffer := bufio.NewWriter(byteArr)
	builder.Write(buffer)
	buffer.Flush()
	return byteArr.Bytes()
}

func (pdfBuilder *pdfBuilder) setField(field string, builder *gopdf.GoPdf) {
	builder.Cell(nil, fmt.Sprintf("%s:", field))
	builder.SetX(builder.GetX() + 10)
}

func (pdfBuilder *pdfBuilder) setValue(value string, builder *gopdf.GoPdf) {
	builder.Cell(nil, value)
}

func (pdfBuilder *pdfBuilder) buildRowSection(data *domain.Member, builder *gopdf.GoPdf) {
	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.Name", "Name"), builder)
	pdfBuilder.setValue(data.Person.GetFullName(), builder)
	builder.Br(15)
	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.Classification", "Classification"), builder)
	classification := data.Classification()
	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.Classification."+classification.String(), classification.String()), builder)
	builder.Br(15)
	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.Address", "Address"), builder)
	pdfBuilder.setValue(data.Person.Address.String(), builder)
	builder.Br(15)

	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.BirthDate", "Birth date"), builder)
	pdfBuilder.setValue(data.Person.BirthDate.Format("02/01/2006"), builder)
	builder.SetX(builder.GetX() + 10)
	if data.Person.MarriageDate != nil {
		pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.MarriageDate", "Marriage Date"), builder)
		pdfBuilder.setValue(data.Person.MarriageDate.Format("02/01/2006"), builder)
	}
	builder.Br(15)

	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.Phone", "Phone"), builder)

	pdfBuilder.setValue(data.Person.Contact.GetFormattedPhone(), builder)

	builder.SetX(builder.GetX() + 10)
	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.CellPhone", "CellPhone"), builder)

	pdfBuilder.setValue(data.Person.Contact.GetFormattedCellPhone(), builder)
	builder.Br(15)

	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.Email", "Email"), builder)
	pdfBuilder.setValue(data.Person.Contact.Email, builder)

	builder.Br(10)
	builder.SetLineWidth(2)
	builder.SetLineType("dashed")
	builder.Line(builder.GetX(), builder.GetY()+10, gopdf.PageSizeA4.W-30, builder.GetY()+10)
	builder.Br(25)

}

func (pdfBuilder *pdfBuilder) buildSummarySection(data []*domain.Member, builder *gopdf.GoPdf) {
	pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.MembersQuantity", "Number of members"), builder)
	pdfBuilder.setValue(fmt.Sprintf("%d", len(data)), builder)
	builder.Br(15)

	summary := map[string]int{}
	for _, member := range data {
		count := summary[member.Classification().String()]
		count++
		summary[member.Classification().String()] = count
	}

	for key, value := range summary {
		pdfBuilder.setField(pdfBuilder.messageService.GetMessage("Domain.Classification." + key, key), builder)
		pdfBuilder.setValue(fmt.Sprintf("%d", value), builder)
		builder.Br(15)
	}

}

func (pdfBuilder *pdfBuilder) BuildFile(title string, data []*domain.Member) ([]byte, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err := pdfBuilder.setFont(pdf); err != nil {
		return nil, err
	}
	pdfBuilder.buildFirstPageSection(title, pdf)
	pdf.AddPage()
	const maxPerPage int = 7
	count := 0
	for _, member := range data {
		if count == maxPerPage {
			pdf.AddPage()
			count = 0
		}
		count++
		pdfBuilder.buildRowSection(member, pdf)
	}
	pdf.AddPage()
	pdfBuilder.buildSummarySection(data, pdf)
	return pdfBuilder.toBytes(pdf), nil
}
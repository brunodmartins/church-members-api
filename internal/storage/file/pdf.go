package file

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	tr "github.com/BrunoDM2943/church-members-api/internal/infra/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/signintech/gopdf"
	"github.com/spf13/viper"
)

//Builder interface
//go:generate mockgen -source=./pdf.go -destination=./mock/pdf_mock.go
type Builder interface {
	BuildFile(title string, data []*model.Member) ([]byte, error)
}

type pdfBuilder struct {

}

func NewPDFBuilder() *pdfBuilder {
	return &pdfBuilder{}
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
	builder.AddTTFFont("arial", "./arial.ttf")
	builder.AddTTFFont("arial", "./fonts/arial.ttf")

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

func (pdfBuilder *pdfBuilder) buildRowSection(data *model.Member, builder *gopdf.GoPdf) {
	pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Name",
			Other: "Name",
		},
	}), builder)
	pdfBuilder.setValue(data.Person.GetFullName(), builder)
	builder.Br(15)
	pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Classification",
			Other: "Classification",
		},
	}), builder)
	pdfBuilder.setValue(data.ClassificationLocalized(tr.Localizer), builder)
	builder.Br(15)
	pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Address",
			Other: "Address",
		},
	}), builder)
	pdfBuilder.setValue(data.Person.Address.GetFormatted(), builder)
	builder.Br(15)

	pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.BirthDate",
			Other: "Birth date",
		},
	}), builder)
	pdfBuilder.setValue(data.Person.BirthDate.Format("02/01/2006"), builder)
	builder.SetX(builder.GetX() + 10)
	if !data.Person.MarriageDate.IsZero() {
		pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Domain.MarriageDate",
				Other: "Marriage Date",
			},
		}), builder)
		pdfBuilder.setValue(data.Person.MarriageDate.Format("02/01/2006"), builder)
	}
	builder.Br(15)

	pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Phone",
			Other: "Phone",
		},
	}), builder)
	if data.Person.Contact.Phone == 0 {
		pdfBuilder.setValue("N/A", builder)
	} else {
		pdfBuilder.setValue(data.Person.Contact.GetFormattedPhone(), builder)
	}
	builder.SetX(builder.GetX() + 10)
	pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.CellPhone",
			Other: "CellPhone",
		},
	}), builder)
	if data.Person.Contact.CellPhone == 0 {
		pdfBuilder.setValue("N/A", builder)
	} else {
		pdfBuilder.setValue(data.Person.Contact.GetFormattedCellPhone(), builder)
	}
	builder.Br(15)

	pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Email",
			Other: "Email",
		},
	}), builder)
	pdfBuilder.setValue(data.Person.Contact.Email, builder)

	builder.Br(10)
	builder.SetLineWidth(2)
	builder.SetLineType("dashed")
	builder.Line(builder.GetX(), builder.GetY()+10, gopdf.PageSizeA4.W-30, builder.GetY()+10)
	builder.Br(25)

}

func (pdfBuilder *pdfBuilder) buildSummarySection(data []*model.Member, builder *gopdf.GoPdf) {
	pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.MembersQuantity",
			Other: "Number of members",
		},
	}), builder)
	pdfBuilder.setValue(fmt.Sprintf("%d", len(data)), builder)
	builder.Br(15)

	summary := map[string]int{}
	for _, member := range data {
		count := summary[member.Classification()]
		count++
		summary[member.Classification()] = count
	}

	for key, value := range summary {
		pdfBuilder.setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Domain.Classification." + key,
				Other: key,
			},
		}), builder)
		pdfBuilder.setValue(fmt.Sprintf("%d", value), builder)
		builder.Br(15)
	}

}

func (pdfBuilder *pdfBuilder) BuildFile(title string, data []*model.Member) ([]byte, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err := pdfBuilder.setFont(pdf); err != nil {
		panic(err)
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

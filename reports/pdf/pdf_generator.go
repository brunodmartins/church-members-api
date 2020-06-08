package pdf

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/BrunoDM2943/church-members-api/entity"
	tr "github.com/BrunoDM2943/church-members-api/infra/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/signintech/gopdf"
	"github.com/spf13/viper"
)

func buildFirstPageSection(title string, builder *gopdf.GoPdf) {
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

func setFont(builder *gopdf.GoPdf) error {
	builder.AddTTFFont("arial", "./pdf/arial.ttf")
	builder.AddTTFFont("arial", "./fonts/arial.ttf")

	err := builder.SetFont("arial", "", 14)
	if err != nil {
		return err
	}
	return nil
}

func toBytes(builder *gopdf.GoPdf) []byte {
	byteArr := &bytes.Buffer{}
	buffer := bufio.NewWriter(byteArr)
	builder.Write(buffer)
	buffer.Flush()
	return byteArr.Bytes()
}

func setField(field string, builder *gopdf.GoPdf) {
	builder.Cell(nil, fmt.Sprintf("%s:", field))
	builder.SetX(builder.GetX() + 10)
}

func setValue(value string, builder *gopdf.GoPdf) {
	builder.Cell(nil, value)
}

func buildRowSection(data *entity.Member, builder *gopdf.GoPdf) {
	setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Name",
			Other: "Name",
		},
	}), builder)
	setValue(data.Person.GetFullName(), builder)
	builder.Br(15)
	setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Classification",
			Other: "Classification",
		},
	}), builder)
	setValue(data.ClassificationLocalized(tr.Localizer), builder)
	builder.Br(15)
	setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Address",
			Other: "Address",
		},
	}), builder)
	setValue(data.Person.Address.GetFormatted(), builder)
	builder.Br(15)

	setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.BirthDate",
			Other: "Birth date",
		},
	}), builder)
	setValue(data.Person.BirthDate.Format("02/01/2006"), builder)
	builder.SetX(builder.GetX() + 10)
	if !data.Person.MarriageDate.IsZero() {
		setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Domain.MarriageDate",
				Other: "Marriage Date",
			},
		}), builder)
		setValue(data.Person.MarriageDate.Format("02/01/2006"), builder)
	}
	builder.Br(15)

	setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Phone",
			Other: "Phone",
		},
	}), builder)
	if data.Person.Contact.Phone == 0 {
		setValue("N/A", builder)
	} else {
		setValue(data.Person.Contact.GetFormattedPhone(), builder)
	}
	builder.SetX(builder.GetX() + 10)
	setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.CellPhone",
			Other: "CellPhone",
		},
	}), builder)
	if data.Person.Contact.CellPhone == 0 {
		setValue("N/A", builder)
	} else {
		setValue(data.Person.Contact.GetFormattedCellPhone(), builder)
	}
	builder.Br(15)

	setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Email",
			Other: "Email",
		},
	}), builder)
	setValue(data.Person.Contact.Email, builder)

	builder.Br(10)
	builder.SetLineWidth(2)
	builder.SetLineType("dashed")
	builder.Line(builder.GetX(), builder.GetY()+10, gopdf.PageSizeA4.W-30, builder.GetY()+10)
	builder.Br(25)

}

func buildSummarySection(data []*entity.Member, builder *gopdf.GoPdf) {
	setField(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.MembersQuantity",
			Other: "Number of members",
		},
	}), builder)
	setValue(fmt.Sprintf("%d", len(data)), builder)
	builder.Br(15)

	summary := map[string]int{}
	for _, member := range data {
		count := summary[member.Classification()]
		count++
		summary[member.ClassificationLocalized(tr.Localizer)] = count
	}

	for key, value := range summary {
		setField(key, builder)
		setValue(fmt.Sprintf("%d", value), builder)
		builder.Br(15)
	}

}

func BuildPdf(title string, data []*entity.Member) ([]byte, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err := setFont(pdf); err != nil {
		panic(err)
	}
	buildFirstPageSection(title, pdf)
	pdf.AddPage()
	const maxPerPage int = 7
	count := 0
	for _, membro := range data {
		if count == maxPerPage {
			pdf.AddPage()
			count = 0
		}
		count++
		buildRowSection(membro, pdf)
	}
	pdf.AddPage()
	buildSummarySection(data, pdf)
	return toBytes(pdf), nil
}

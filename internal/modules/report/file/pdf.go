package file

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"github.com/brunodmartins/church-members-api/platform/i18n"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/signintech/gopdf"
)

//go:embed fonts/*.ttf
var fontFS embed.FS

// Builder interface
//
//go:generate mockgen -source=./pdf.go -destination=./mock/pdf_mock.go
type Builder interface {
	BuildFile(ctx context.Context, title string, church *domain.Church, members []*domain.Member) ([]byte, error)
}

type pdfBuilder struct {
}

func NewPDFBuilder() *pdfBuilder {
	return &pdfBuilder{}
}

func (pdfBuilder *pdfBuilder) buildFirstPageSection(title string, church *domain.Church, builder *gopdf.GoPdf) {
	builder.AddPage()
	builder.SetMargins(30, 30, 30, 30)
	rect := &gopdf.Rect{H: gopdf.PageSizeA4.H, W: gopdf.PageSizeA4.W}
	options := gopdf.CellOption{
		Align: gopdf.Center,
	}

	builder.SetX(0)
	builder.SetY(100)
	builder.CellWithOption(rect, church.Name, options)
	builder.SetY(120)
	builder.SetX(0)
	builder.CellWithOption(rect, title, options)
}

func (pdfBuilder *pdfBuilder) setFont(builder *gopdf.GoPdf) error {
	file, _ := fontFS.Open("fonts/Arial.ttf")
	builder.AddTTFFontByReader("arial", file)
	_ = builder.SetFont("arial", "", 14)
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

func (pdfBuilder *pdfBuilder) buildRowSection(ctx context.Context, data *domain.Member, builder *gopdf.GoPdf) {
	pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.Name"), builder)
	pdfBuilder.setValue(data.Person.GetFullName(), builder)
	builder.Br(15)
	pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.Classification"), builder)
	classification := data.Classification()
	pdfBuilder.setValue(i18n.GetMessage(ctx, "Domain.Classification."+classification.String()), builder)
	builder.Br(15)
	pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.Address"), builder)
	pdfBuilder.setValue(data.Person.Address.String(), builder)
	builder.Br(15)

	pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.BirthDate"), builder)
	pdfBuilder.setValue(data.Person.BirthDate.Format("02/01/2006"), builder)
	builder.SetX(builder.GetX() + 10)
	if data.Person.MarriageDate != nil {
		pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.MarriageDate"), builder)
		pdfBuilder.setValue(data.Person.MarriageDate.Format("02/01/2006"), builder)
	}
	builder.Br(15)

	pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.Phone"), builder)

	pdfBuilder.setValue(data.Person.Contact.GetFormattedPhone(), builder)

	builder.SetX(builder.GetX() + 10)
	pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.CellPhone"), builder)

	pdfBuilder.setValue(data.Person.Contact.GetFormattedCellPhone(), builder)
	builder.Br(15)

	pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.Email"), builder)
	pdfBuilder.setValue(data.Person.Contact.Email, builder)

	builder.Br(10)
	builder.SetLineWidth(2)
	builder.SetLineType("dashed")
	builder.Line(builder.GetX(), builder.GetY()+10, gopdf.PageSizeA4.W-30, builder.GetY()+10)
	builder.Br(25)

}

func (pdfBuilder *pdfBuilder) buildSummarySection(ctx context.Context, data []*domain.Member, builder *gopdf.GoPdf) {
	pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.MembersQuantity"), builder)
	pdfBuilder.setValue(fmt.Sprintf("%d", len(data)), builder)
	builder.Br(15)

	summary := map[string]int{}
	for _, member := range data {
		count := summary[member.Classification().String()]
		count++
		summary[member.Classification().String()] = count
	}

	for key, value := range summary {
		pdfBuilder.setField(i18n.GetMessage(ctx, "Domain.Classification."+key), builder)
		pdfBuilder.setValue(fmt.Sprintf("%d", value), builder)
		builder.Br(15)
	}

}

func (pdfBuilder *pdfBuilder) BuildFile(ctx context.Context, title string, church *domain.Church, data []*domain.Member) ([]byte, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err := pdfBuilder.setFont(pdf); err != nil {
		return nil, err
	}
	pdfBuilder.buildFirstPageSection(title, church, pdf)
	pdf.AddPage()
	const maxPerPage int = 7
	count := 0
	for _, member := range data {
		if count == maxPerPage {
			pdf.AddPage()
			count = 0
		}
		count++
		pdfBuilder.buildRowSection(ctx, member, pdf)
	}
	pdf.AddPage()
	pdfBuilder.buildSummarySection(ctx, data, pdf)
	return pdfBuilder.toBytes(pdf), nil
}

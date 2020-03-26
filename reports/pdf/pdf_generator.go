package pdf

import (
	"bufio"
	"bytes"

	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/signintech/gopdf"
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
	builder.CellWithOption(rect, "Igreja Evangélica de Pinheiros na Vila São Francisco", options)
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
	builder.Cell(nil, field)
	builder.SetX(builder.GetX() + 10)
}

func setValue(value string, builder *gopdf.GoPdf) {
	builder.Cell(nil, value)
}

func buildRowSection(data *entity.Membro, builder *gopdf.GoPdf) {
	setField("Nome:", builder)
	setValue(data.Pessoa.GetFullName(), builder)
	builder.Br(15)
	setField("Classificacao:", builder)
	setValue(data.Classificacao(), builder)
	builder.Br(15)
	setField("Endereco:", builder)
	setValue(data.Pessoa.Endereco.GetFormatted(), builder)
	builder.Br(15)

	setField("Dt. Nascimento:", builder)
	setValue(data.Pessoa.DtNascimento.Format("02/01/2006"), builder)
	builder.SetX(builder.GetX() + 10)
	if !data.Pessoa.DtCasamento.IsZero() {
		setField("Dt. Casamento:", builder)
		setValue(data.Pessoa.DtCasamento.Format("02/01/2006"), builder)
	}
	builder.Br(15)

	setField("Telefone:", builder)
	if data.Pessoa.Contato.Telefone == 0 {
		setValue("N/A", builder)
	} else {
		setValue(data.Pessoa.Contato.GetFormattedPhone(), builder)
	}
	builder.SetX(builder.GetX() + 10)
	setField("Celular:", builder)
	if data.Pessoa.Contato.Celular == 0 {
		setValue("N/A", builder)
	} else {
		setValue(data.Pessoa.Contato.GetFormattedCellPhone(), builder)
	}
	builder.Br(15)

	setField("Email:", builder)
	setValue(data.Pessoa.Contato.Email, builder)

	builder.Br(10)
	builder.SetLineWidth(2)
	builder.SetLineType("dashed")
	builder.Line(builder.GetX(), builder.GetY()+10, gopdf.PageSizeA4.W-30, builder.GetY()+10)
	builder.Br(25)

}

func BuildPdf(title string, data []*entity.Membro) ([]byte, error) {
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
	return toBytes(pdf), nil
}

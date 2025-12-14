package render

import (
	"final-project/intern/payment-system/application"
	"final-project/intern/payment-system/config"
	"io"
	"text/template"
)

type HTMLRender struct {
	templateDir string
}

func (render *HTMLRender) New(config *config.Config) {
	render.templateDir = config.TemplateDir
}

func (render *HTMLRender) RenderPaymentPage(writer io.Writer, data application.PaymentPageData) error {
	temp, err := template.ParseFiles(render.templateDir + "/payment.html")

	if err != nil {
		return err
	}

	err = temp.Execute(writer, data)
	return err
}

func (render *HTMLRender) RenderMessagePage(writer io.Writer, data application.MessagePageData) error {
	temp, err := template.ParseFiles(render.templateDir + "/message.html")

	if err != nil {
		return err
	}

	err = temp.Execute(writer, data)
	return err
}

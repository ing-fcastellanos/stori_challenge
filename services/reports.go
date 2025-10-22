package services

import (
	"fintech-backend/utils"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmailReport envía un correo electrónico con el resumen de la migración
func SendEmailReport(success bool, totalTransactions int, totalCredits float64, totalDebits float64, mailer utils.Mailer) error {
	// Configura el mensaje del correo
	subject := "Informe de Migración de Transacciones"
	body := fmt.Sprintf("La migración de transacciones ha %s.\n\nTotal de transacciones: %d\nTotal de créditos: %.2f\nTotal de débitos: %.2f",
		func() string {
			if success {
				return "sido exitosa"
			}
			return "fallado"
		}(),
		totalTransactions, totalCredits, totalDebits)

	var SMTPUser = os.Getenv("SMTPUser")
	var ReportDistList = os.Getenv("REPORT_DIST_LIST")

	// Configuración del correo
	mailerMessage := gomail.NewMessage()
	// TODO: Agregar un template en html para los envios de correos
	mailerMessage.SetHeader("From", SMTPUser)
	mailerMessage.SetHeader("To", ReportDistList)
	mailerMessage.SetHeader("Subject", subject)
	mailerMessage.SetBody("text/plain", body)

	// Enviar el correo
	err := mailer.DialAndSend(*mailerMessage)
	if err != nil {
		return fmt.Errorf("error al enviar el correo: %v", err)
	}
	return nil
}

package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/gomail.v2"
)

// MockMailer es un mock de la interfaz Mailer
type MockMailer struct {
	mock.Mock
}

// DialAndSend simula el método DialAndSend de Mailer
func (m *MockMailer) DialAndSend(message gomail.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestSendEmailReport_Success(t *testing.T) {
	// Crear un mock de Mailer
	mockMailer := new(MockMailer)

	// Configurar las expectativas del mock
	mockMailer.On("DialAndSend", mock.Anything).Return(nil)

	// Llamar a la función SendEmailReport con los datos de prueba
	err := SendEmailReport(true, 10, 500.00, -200.00, mockMailer)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que el mock haya sido llamado
	mockMailer.AssertExpectations(t)
}

func TestSendEmailReport_ErrorSendingEmail(t *testing.T) {
	// Crear un mock de Mailer
	mockMailer := new(MockMailer)

	// Configurar el mock para simular un error al enviar el correo
	mockMailer.On("DialAndSend", mock.Anything).Return(fmt.Errorf("SMTP error"))

	// Llamar a la función SendEmailReport
	err := SendEmailReport(true, 10, 500.00, -200.00, mockMailer)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, "error al enviar el correo: SMTP error", err.Error())

	// Verificar que el mock haya sido llamado
	mockMailer.AssertExpectations(t)
}

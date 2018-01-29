package model

//Email - representa um e-mail na aplicação
type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

package erro

import (
	"errors"
)

var (
	ErrServer		 	= errors.New("erro não identificado")
	ErrNotFound 		= errors.New("item não encontrado")
	ErrUnauthorized 	= errors.New("erro de autorização")
	ErrDecode			= errors.New("erro na decodificação do Base64")
	ErrUpdate			= errors.New("erro no update do dado")
	ErrEvent			= errors.New("erro no tratamento de evento")
)

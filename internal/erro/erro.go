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
	ErrUnmarshal 		= errors.New("unmarshal json error")
	ErrInsert 			= errors.New("insert data error")
	ErrEvent			= errors.New("erro no tratamento de evento")
)

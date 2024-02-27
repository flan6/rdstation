# RDStaton

Lib de integração com a RDStation para buscar, editar e deletar um Lead
## Get started

Para instanciar o client basta passar as credenciais para o construtor:
```go
	ClientID := "<client_id>"
	ClientSecret := "<client_secret>"
	Token := rdstation.Token{
		AccessToken:  "<access_token>",
		RefreshToken: "<refresh_token>",
		ExpiresIn:    86400,
		CreationDate: "2022-07-06 11:54:26",
	}

	rd := rdstation.NewRDStation(ClientID, ClientSecret, Token)
```
Não é necessário atualizar o access-token, pois o client gerencia o OAuth de forma independente.

## Exemplo

Para executar o exemplo edite as credenciais do RDStation no `examples/example.go` e execute o comando:

```sh
    go run ./examples/example.go
```

### Documentação auxiliar

[Documentaçao da API da RDStation](https://developers.rdstation.com/reference/contatos)

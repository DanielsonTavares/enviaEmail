package main

/*
  Aplicação para enviar um email toda a vez que o windows for iniciado.
  A configuração do servidor SMTP a ser utilizado deve ser colocada
  no arquivo configsmtp.txt, no mesmo diretório do executável.
  Formato do arquivo configsmtp.txt:
  	server = nome_do_servidor_smtp
	port = porta
	username = nome_usuario
	password = senha
	from = remetente@email.com
	to = destinatario@email.com

	Autor: caosdan@gmail.com
	Data: 20/11/2017
	Versão: 1.0

*/

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/DanielsonTavares/utildan"
)

const nomeArquivoRegEnvio = "regEnvio.txt"

func main() {
	dDataHora := time.Now().Format("02/01/2006 15:04:05")
	dData := time.Now().Format("02/01/2006")
	var comando string

	if verificaEnvio(dData) == false {
		sSMTPConfig := utildan.RecuperaConfiguracao("configsmtp.txt")

		to := []string{sSMTPConfig["to"]}

		msg := []byte("To: " + to[0] + "\r\n" +
			"subject: Horario de entrada - " + dData + "\r\n" +
			"\r\n" +
			"Bom dia,\r\n" +
			"segue meu horario de entrada " + dDataHora)

		auth := unencryptedAuth{
			smtp.PlainAuth(
				"",
				sSMTPConfig["username"],
				sSMTPConfig["password"],
				sSMTPConfig["server"],
			),
		}

		err := smtp.SendMail(sSMTPConfig["server"]+":"+sSMTPConfig["port"], auth, sSMTPConfig["from"], to, msg)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Email enviado com sucesso!")
		geraRegEnvio()

	} else {
		fmt.Println("Email já enviado hoje!")
	}

	fmt.Println("Pressione qualquer tecla para terminar...")
	fmt.Scanln(&comando)
	os.Exit(0)

}

/*unencryptedAuth é utilizado para realizar uma autenticação não encriptada.

fonte: https://stackoverflow.com/a/11066064/6786759
*/
type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

//Gera o registro do email enviado no arquivo definido na constante nomeArquivoRegEnvio
func geraRegEnvio() {
	arquivo, err := os.OpenFile(nomeArquivoRegEnvio, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		fmt.Println("Erro ao ler arquivo: ", err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - Email enviado \r\n")

	arquivo.Close()
}

// func geraLog(txtLog string) {
// 	arquivo, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

// 	if err != nil {
// 		fmt.Println("Erro ao ler arquivo: ", err)
// 	}

// 	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + "\r\n" +
// 		txtLog + "\r\n" +
// 		"=====================================FIM===========================================\r\n")

// 	arquivo.Close()
// }

//verificaEnvio verifica se já houve emissão de email no dia informado no parâmetro sDataEnvio
func verificaEnvio(sDataEnvio string) bool {
	arquivo, err := os.OpenFile(nomeArquivoRegEnvio, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	var linha string

	if err != nil {
		fmt.Println("Erro ao ler arquivo: ", err)
	}

	leitor := bufio.NewReader(arquivo)

	for {

		linhaArq, err := leitor.ReadString('\n')
		linhaArq = strings.TrimSpace(linhaArq)

		if err == io.EOF {
			break
		}

		if linhaArq != "" {
			linha = linhaArq
		}

	}

	if len(linha) >= 10 {
		sDataLog := linha[0:10]
		if sDataLog == sDataEnvio {
			arquivo.Close()
			return true
		} else {
			arquivo.Close()
			return false
		}

	} else {
		return false
	}

}

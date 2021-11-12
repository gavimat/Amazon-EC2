package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Message:  Define the Message structure, with 4 properties.  Structure tags are used by encoding/json library
type MessageInfos struct {
	Remetente      string `json:"Remetente"`
	NumeroMensagem string `json:"NumeroMensagem"`
	Mensagem       string `json:"Mensagem"`
	Status         string `json:"Status"`
}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("fabcar_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryMensagem":
		return s.queryMensagem(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "escreveMensagem":
		return s.escreveMensagem(APIstub, args)
	case "queryTodasMensagens":
		return s.queryTodasMensagens(APIstub)
	case "queryMensagensDaOutraOrg":
		return s.queryMensagensDaOutraOrg(APIstub, args)
	case "queryMensagensNaoLidas":
		return s.queryMensagensNaoLidas(APIstub, args)
	case "alterarStatusMensagem":
		return s.alterarStatusMensagem(APIstub, args)
	case "geraNotificacao":
		return s.geraNotificacao(APIstub)

	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

func (s *SmartContract) queryMensagem(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Número incorreto de argumentos. Esperado 1 argumento.")
	}

	messageAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(messageAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	mensagens := []MessageInfos{
		MessageInfos{Remetente: "ORG1", NumeroMensagem: "0", Mensagem: "Mensagem de inicialização do Ledger", Status: "Inicializador"}}

	i := 0
	for i < len(mensagens) {
		messageAsBytes, _ := json.Marshal(mensagens[i])
		APIstub.PutState("Mensagem"+strconv.Itoa(i), messageAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) escreveMensagem(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Número incorreto de argumentos. Esperado 4 argumentos.")
	}

	var mensagem = MessageInfos{Remetente: args[1], NumeroMensagem: args[2], Mensagem: args[3], Status: args[4]}

	messageAsBytes, _ := json.Marshal(mensagem)
	APIstub.PutState(args[0], messageAsBytes)

	indexName := "owner~key"
	colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{mensagem.NumeroMensagem, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	APIstub.PutState(colorNameIndexKey, value)

	return shim.Success(messageAsBytes)
}

func (s *SmartContract) queryTodasMensagens(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Mensagem0"
	endKey := "Mensagem999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("\n")
	buffer.WriteString("Mensagens enviadas no Ledger ")
	buffer.WriteString(": \n")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}

	fmt.Printf("- queryAll:\n %s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryMensagensDaOutraOrg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	startKey := "Mensagem0"
	endKey := "Mensagem999"

	var organizacao = args[0]

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString(": \n")
	buffer.WriteString("MENSAGENS ENCAMINHADAS PELA ")
	buffer.WriteString(organizacao)
	buffer.WriteString(": \n")

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var testeString = string(queryResponse.Value)

		for i := 0; i < len(testeString); i++ {
			if string(testeString[i]) == string(organizacao[0]) && string(testeString[i+1]) == string(organizacao[1]) && string(testeString[i+2]) == string(organizacao[2]) && string(testeString[i+3]) == string(organizacao[3]) {

				buffer.WriteString("{\"Key\":")
				buffer.WriteString("\"")
				buffer.WriteString(queryResponse.Key)
				buffer.WriteString("\"")

				buffer.WriteString(", \"Record\":")
				// Record is a JSON object, so we write as-is
				buffer.WriteString(string(queryResponse.Value))
				buffer.WriteString("}")
				buffer.WriteString(",")
				buffer.WriteString("\n")

			}
		}
	}

	fmt.Printf("- queryAll:\n %s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryMensagensNaoLidas(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	startKey := "Mensagem1"
	endKey := "Mensagem999"

	var organizacao = args[0]

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("\n")
	buffer.WriteString("MENSAGENS ENCAMINHADAS PELA ")
	buffer.WriteString(organizacao)
	buffer.WriteString(": \n")

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var testeString = string(queryResponse.Value)

		for i := 0; i < len(testeString); i++ {
			if string(testeString[i]) == string(organizacao[0]) && string(testeString[i+1]) == string(organizacao[1]) && string(testeString[i+2]) == string(organizacao[2]) && string(testeString[i+3]) == string(organizacao[3]) {

				for j := 0; j < len(testeString); j++ {
					if string(testeString[j]) == "N" && string(testeString[j+1]) == "L" && string(testeString[j+2]) == "I" && string(testeString[j+3]) == "D" && string(testeString[j+4]) == "O" {

						buffer.WriteString("{\"Key\":")
						buffer.WriteString("\"")
						buffer.WriteString(queryResponse.Key)
						buffer.WriteString("\"")

						buffer.WriteString(", \"Record\":")
						// Record is a JSON object, so we write as-is
						buffer.WriteString(string(queryResponse.Value))
						buffer.WriteString("}")
						buffer.WriteString(",")
						buffer.WriteString("\n")

						// // ============== Mudar Status ==============
						// messageAsBytes, _ := APIstub.GetState(queryResponse.Key)
						// message := MessageInfos{}

						// json.Unmarshal(messageAsBytes, &message)
						// message.Status = "LIDO"
						// messageAsBytes, _ = json.Marshal(message)
						// APIstub.PutState(queryResponse.Key, messageAsBytes)
						// return shim.Success(messageAsBytes)
					}
				}

			}
		}

	}

	fmt.Printf("- queryAll:\n %s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) alterarStatusMensagem(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Número incorreto de argumentos. Esperado 1 argumento.")
	}

	startKey := "Mensagem1"
	endKey := "Mensagem999"

	var organizacao = args[0]

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var testeString = string(queryResponse.Value)

		for i := 0; i < len(testeString); i++ {
			if string(testeString[i]) == string(organizacao[0]) && string(testeString[i+1]) == string(organizacao[1]) && string(testeString[i+2]) == string(organizacao[2]) && string(testeString[i+3]) == string(organizacao[3]) {

				for j := 0; j < len(testeString); j++ {
					if string(testeString[j]) == "N" && string(testeString[j+1]) == "L" && string(testeString[j+2]) == "I" && string(testeString[j+3]) == "D" && string(testeString[j+4]) == "O" {

						// ============== Mudar Status ==============
						messageAsBytes, _ := APIstub.GetState(queryResponse.Key)
						message := MessageInfos{}

						json.Unmarshal(messageAsBytes, &message)
						message.Status = "LIDO"
						messageAsBytes, _ = json.Marshal(message)
						APIstub.PutState(queryResponse.Key, messageAsBytes)
						//return shim.Success(messageAsBytes)
					}
				}

			}
		}

	}
	return shim.Success(nil)
}

func (s *SmartContract) geraNotificacao(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Mensagem0"
	endKey := "Mensagem999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	var contadorMensagem = 0

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var testeString = string(queryResponse.Value)

		for i := 0; i < len(testeString); i++ {
			if string(testeString[i]) == "N" && string(testeString[i+1]) == "L" && string(testeString[i+2]) == "I" && string(testeString[i+3]) == "D" && string(testeString[i+4]) == "O" {
				contadorMensagem++
			}
		}
	}

	if contadorMensagem > 0 {
		buffer.WriteString("Existem ")
		buffer.WriteString(strconv.Itoa(contadorMensagem))
		buffer.WriteString(" mensagens 'Não lidas' no Ledger")
	}

	if contadorMensagem == 0 {
		buffer.WriteString("Não existem mensagens 'Não Lidas' no Ledger")
	}

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

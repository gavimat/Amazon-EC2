export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export PEER0_ORG2_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/artifacts/channel/config/

#export PRIVATE_DATA_CONFIG=${PWD}/artifacts/private-data/collections_config.json

export CHANNEL_NAME=mychannel

setGlobalsForOrderer() {
    export CORE_PEER_LOCALMSPID="OrdererMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp

}

setGlobalsForPeer0Org1() {
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}

setGlobalsForPeer1Org1() {
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:8051

}

setGlobalsForPeer0Org2() {
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051

}

setGlobalsForPeer1Org2() {
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:10051

}

CHANNEL_NAME="mychannel"
CC_RUNTIME_LANGUAGE="golang"
VERSION="1"
CC_SRC_PATH="./artifacts/src/github.com/fabcar/go"
CC_NAME="fabcar"


#===============================================
#FUNCAO A SER CRIADA AQUI
#===============================================
escreverMensagem() {

    numeroMensagens=0

  #Solicita a escolha da organiza????o a qual ser?? realizada a opera????o e faz a leitura
	echo "Escolha a Organizacao com a qual deseja realizar a operacao:"
	
	read -p "Digite 1 para ORG1 ou 2 para ORG2 (Usar _ nos espa??os) : " Org
	
	if [ "$Org" -eq 1 ]; then

		setGlobalsForPeer0Org1
		
		#Define o nome da organiza????o, o qual ser?? enviado como par??metro para a execu????o da fun????o no chaincode
		numeroOrg=$( printf '%d' $Org )
		nomeOrgRemetente="ORG"
		nomeOrgRemetente+=$numeroOrg
		echo "Organizacao escolhida: ${nomeOrgRemetente}"

		#Solicita a escrita da mensagem a ser enviada e realiza a inser????o do identificador
		read -p "Digite a mensagem que voc?? deseja enviar para a Org 2: " mensagemOrg1
		echo "A mensagem adicionada no Ledger ser??: "${mensagemOrg1}""
		
		#Define o ID da opera????o, este que ser?? utilizado como key para a inser????o da mensagem no Ledger
		((numeroMensagens=numeroMensagens+1))
		numeroMensagens=$( printf '%d' $numeroMensagens )
		idMensagem="Mensagem"
		idMensagem+=$numeroMensagens
		echo "O ID da mensagem ?? ${idMensagem}"

		 # Chama func escreveMensagem no ChainCode
    /home/ubuntu/fabric-samples/bin/peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA \
        -C $CHANNEL_NAME -n ${CC_NAME}  \
        --peerAddresses localhost:7051 \
        --tlsRootCertFiles $PEER0_ORG1_CA \
        --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG2_CA   \
        -c '{"function": "escreveMensagem","Args":["'${idMensagem}'", "'${nomeOrgRemetente}'", "'${numeroMensagens}'", "'${mensagemOrg1}'", "NLIDO"]}'

	elif [ "$Org" -eq 2 ]; then
		setGlobalsForPeer0Org2

		numeroOrg=$( printf '%d' $Org )
		nomeOrgRemetente="ORG"
		nomeOrgRemetente+=$numeroOrg
		echo "Organizacao escolhida: ${nomeOrgRemetente}"
		
		read -p "Digite a mensagem que voc?? deseja enviar para a Org 1 (Usar _ nos espa??os) : " mensagemOrg2
		echo "A mensagem adicionada no Ledger ser??: "${mensagemOrg2}""

		((numeroMensagens=numeroMensagens+1))
		numeroMensagens=$( printf '%d' $numeroMensagens )
		idMensagem="Mensagem"
		idMensagem+=7
		echo "O ID da mensagem ?? ${idMensagem}"

		 # Chama func escreveMensagem no ChainCode
    /home/ubuntu/fabric-samples/bin/peer chaincode invoke -o localhost:7050 \
        --ordererTLSHostnameOverride orderer.example.com \
        --tls $CORE_PEER_TLS_ENABLED \
        --cafile $ORDERER_CA \
        -C $CHANNEL_NAME -n ${CC_NAME}  \
        --peerAddresses localhost:7051 \
        --tlsRootCertFiles $PEER0_ORG1_CA \
        --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_ORG2_CA   \
        -c '{"function": "escreveMensagem","Args":["'${idMensagem}'", "'${nomeOrgRemetente}'", "'${numeroMensagens}'", "'${mensagemOrg1}'", "NLIDO"]}'

	else 
    echo "O valor digitado est?? incorreto"
	fi

}


#numeroMensagensOrg1 = 0
#numeroMensagensOrg2 = 0

escreverMensagem

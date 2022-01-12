/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package util

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	token "github.com/wangyi/fishpond/eth"
	"log"
	"math/big"
	"os"
	"testing"

	//	token "./contracts_erc20" // for demo
)

func TestOne(t *testing.T) {

	wei := new(big.Int)
	wei.SetString("174093039657465489", 10)

	eth := ToDecimal(wei, 18)
	fmt.Println(eth) // 0.02
	b := decimal.NewFromFloat(3217.54)
	fmt.Println(eth.Mul(b)) // 0.02
	c := eth.Mul(b)
	fmt.Println(c.IntPart())
}

func TestTwo(t *testing.T) {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/a73f8de0d9974dd7a35c5d241e24e853")
	if err != nil {
		fmt.Println("455")
		return
	}
	BPrivateKey := "3da233bb9e8629df0f4fb199c4e6529b6a89c6a5ba659ecdbc8c2f18509dbcce" //b的私钥
	privateKey, err := crypto.HexToECDSA(BPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//接下来我们可以读取我们应该用于帐户交易的随机数。
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fromAddress, nonce)
	//下一步是设置我们将要转移的ETH数量。 但是我们必须将ETH以太转换为wei，因为这是以太坊区块链所使用的。 以太网支持最多18个小数位，因此1个ETH为1加18个零。
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	//ETH转账的燃气应设上限为“21000”单位。
	gasLimit := uint64(21000) // in units
	//燃气价格必须以wei为单位设定。 在撰写本文时，将在一个区块中比较快的打包交易的燃气价格为30 gwei。
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value, gasLimit, gasPrice)

	//接下来我们弄清楚我们将ETH发送给谁。
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	//现在我们最终可以通过导入go-ethereumcore/types包并调用NewTransaction来生成我们的未签名以太坊事务，这个函数需要接收nonce，地址，值，燃气上限值，燃气价格和可选发的数据。 发送ETH的数据字段为“nil”。 在与智能合约进行交互时，我们将使用数据字段，仅仅转账以太币是不需要数据字段的。

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	//下一步是使用发件人的私钥对事务进行签名。 为此，我们调用SignTx方法，该方法接受一个未签名的事务和我们之前构造的私钥。 SignTx方法需要EIP155签名者，这个也需要我们先从客户端拿到链ID。
	chainID, err := client.NetworkID(context.Background())

	fmt.Println("+++++")
	fmt.Println(chainID)
	fmt.Println("+++++")
	if err != nil {
		log.Fatal(err)

		fmt.Println("8888")

	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
		fmt.Println("00-0")
		return
	}
	//现在我们终于准备通过在客户端上调用“SendTransaction”来将已签名的事务广播到整个网络。
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
		fmt.Println("000")
		return
	}

	fmt.Println("====")

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0x77006fcb3938f648e2cc65bafd27dec30b9bfbe9df41f78498b9c8b7322a249e

}

func TestTo(t *testing.T) {

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/a73f8de0d9974dd7a35c5d241e24e853")
	if err != nil {
		fmt.Println("455")
		return
	}

	tokenAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") //usDT
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0xbf8F13fFAAffE93DB052AFC50339c6fcEaaF691F")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println(bal)  //100000030

}

func TestTiXian(t *testing.T) {
	type Address struct {
		Type    string
		City    string
		Country string
	}

	type VCard struct {
		FirstName string
		LastName  string
		Addresses []*Address
		Remark    string
	}
	pa := &Address{"private", "Aartselaar", "Belgium"}

	js, _ := json.Marshal(pa)

	fmt.Printf("JSON format: %s", js)


	wa := &Address{"work", "Boom", "Belgium"}
	vc := VCard{"Jan", "Kersschot", []*Address{pa, wa}, "none"}
	// fmt.Printf("%v: \n", vc) // {Jan Kersschot [0x126d2b80 0x126d2be0] none}:
	// JSON format:
	//js, _ := json.Marshal(vc)
	//fmt.Printf("JSON format: %s", js)
	// using an encoder:
	file, _ := os.OpenFile("vcard.json", os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	enc := json.NewEncoder(file)
	err := enc.Encode(vc)
	if err != nil {
		log.Println("Error in encoding json")
	}



}

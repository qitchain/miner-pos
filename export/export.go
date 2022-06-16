package export

import (
	"chia-miner/pkg/bls"
	"encoding/hex"
	bip39 "github.com/tyler-smith/go-bip39"
)

func CreateFarmerKey(masterKey *bls.PrivateKey) *bls.PrivateKey {
	coinType := int(0x800020fc & 0x7fffffff)
	return masterKey.DeriveChild([]int{12381, coinType, 0, 0})
}

func GetFarmerPrivateKeyByMnemonic(minerMnemonic string) (string, string, error) {
	privateKeySeed := bip39.NewSeed(minerMnemonic, "")
	masterPrivateKey, _ := bls.GenerateKeyFromSeed(privateKeySeed)
	//poolPrivateKey := chiaNetwork.CreatePoolKey(masterPrivateKey)
	farmerKey := CreateFarmerKey(masterPrivateKey)
	farmerPrivateKey, _ := farmerKey.MarshalBinary()
	farmerPk, _ := farmerKey.Public().(*bls.PublicKey).MarshalBinary()

	return hex.EncodeToString(farmerPk), hex.EncodeToString(farmerPrivateKey), nil
}

func GetFarmerPublicKey(privateKeyData string) (string, error) {
	data, err := hex.DecodeString(privateKeyData)
	if err != nil {
		return "", err
	}
	privateKey := bls.PrivateKey{}
	if err := privateKey.UnmarshalBinary(data); err != nil {
		return "", err
	}
	farmerPk, _ := privateKey.Public().(*bls.PublicKey).MarshalBinary()
	return hex.EncodeToString(farmerPk), nil
}

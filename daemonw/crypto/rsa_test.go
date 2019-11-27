package crypto

import (
	"io/ioutil"
	"os"
	"testing"
)

const (
	PubKeyStr  = "30819f300d06092a864886f70d010101050003818d0030818902818100c32d95711a79c20bd118a804035b9decf2610c078f04ca607cfd47cb45bb2470e981bc747196dc9e82ff09a4f9747c5ab86c806cc97aff5d6ed9cf3326372962a85b32ab7f06dac5e21c8fbb83666f60064f1cce188d8f95f9f413f381548a65fb547cf446bfd3b712201ca39359a0fc142c1ab88ebc412b54db4fefa0eeab5f0203010001"
	PriKeyStr  = "30820276020100300d06092a864886f70d0101010500048202603082025c02010002818100c32d95711a79c20bd118a804035b9decf2610c078f04ca607cfd47cb45bb2470e981bc747196dc9e82ff09a4f9747c5ab86c806cc97aff5d6ed9cf3326372962a85b32ab7f06dac5e21c8fbb83666f60064f1cce188d8f95f9f413f381548a65fb547cf446bfd3b712201ca39359a0fc142c1ab88ebc412b54db4fefa0eeab5f0203010001028180019412cb515d7eb63dc0a74b3103c2a931e49b064c9fcb4b307793e70eef4e8223414b56ddbc4c987392c06f371d3c0c37c598dc485cd847d727a7c7501dce15dc8cecf31d78324ea20c210cf353d2807cb4aad06729d71d88d863d5376c804de1cf4d8381ede7c6d8caf03ebcb40d4a685788ea9fa8b7b94ddb0d89e57bfa41024100e456e06ecfd7e3913ee06e5da28f9aeb41067505deb8b38f0f876436e0a104674652708160590076fceb3567bb0e321ff54402163aca378bfbcabcf8f5ea4f71024100dad254264351d9d2e585fc22c47369139563efb7d06a75edd1581276774c49ea5d7e073a2c3d6d73f94b963e48f0a0df32415473ae1ae8c442ac66113915dfcf024007675dbd9a15a8f653b0994385fadc0e16bcabdcb5dd0d4d6d3679a461e236af569363bc7e3044989718d55178ba3c1bd72b0e5a442bbc635895b94ce54464b10240470d4485a8164bc44f02d27bf78e01f079f19944f6479b433d567ba83bd469abcfefbc2907746ab5dd9ccd8c2be2f65982d7f3b41c3460a05985e11e768ac3c10241009a5d2ed268f5c75886ee1950ee482af10176ac08a41e8a254f968b802481991cfafc098dede91cffaad3b326f56033a934de878215423cbe74b6d23553fab8a1"
	PlainText  = "Hello,World!"
	CipherText = "ab10fdca4e12e990fdbe66297597a17f6b30308c51fdce2da5002911ac35c3ddf6d8c2c5145b9ff87ee0901c95cf94ed74ba307553e2346cede4de5bcf5177960e03bebeed305d6f6b2fcda70b5ad81af5d8864b395b534a375d5348fcb2ce0d8daab98da13ffcec65cfb295ab0988fe8a9d6e75f91aa71915be0445510c3b02"
)

func TestGenRsaKeyPair(t *testing.T) {
	pub, pri, err := GenRsaKeyPair(2048)
	if err != nil {
		t.Error(err)
	}
	dir := os.TempDir()
	err = ioutil.WriteFile(dir+"/rsa_pub", []byte(Bytes2Hex(pub)), 0666)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(dir+"/rsa_pri", []byte(Bytes2Hex(pri)), 0666)
	if err != nil {
		t.Error(err)
	}

	b, err := ioutil.ReadFile(dir + "/rsa_pub")
	if err != nil {
		t.Error(err)
	}
	pubBytes, err := Hex2Bytes(string(b))
	if err != nil {
		t.Error(err)
	}
	_, err = ParsePublicKey(pubBytes)
	if err != nil {
		t.Error(err)
	}


	b, err = ioutil.ReadFile(dir + "/rsa_pri")
	if err != nil {
		t.Error(err)
	}
	priBytes, err := Hex2Bytes(string(b))
	if err != nil {
		t.Error(err)
	}
	_, err = ParsePrivateKey(priBytes)
	if err != nil {
		t.Error(err)
	}
}

func TestRsaEnc(t *testing.T) {
	b, err := Hex2Bytes(PubKeyStr)
	if err != nil {
		t.Error(err)
	}
	pub, err := ParsePublicKey(b)
	if err != nil {
		t.Error(err)
	}

	b, err = Hex2Bytes(PriKeyStr)
	if err != nil {
		t.Error(err)
	}
	pri, err := ParsePrivateKey(b)
	if err != nil {
		t.Error(err)
	}
	enc, err := RsaEnc(pub, []byte(PlainText))
	if err != nil {
		t.Error(err)
	}
	dec, err := RsaDec(pri, enc)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != PlainText {
		t.Error("result not match")
	}
}

func TestRsaDec(t *testing.T) {
	b, err := Hex2Bytes(PriKeyStr)
	if err != nil {
		t.Error(err)
	}
	key, err := ParsePrivateKey(b)
	if err != nil {
		t.Error(err)
	}
	b, err = Hex2Bytes(CipherText)
	if err != nil {
		t.Error(err)
	}
	dec, err := RsaDec(key, b)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != PlainText {
		t.Error("result not match")
	}
}

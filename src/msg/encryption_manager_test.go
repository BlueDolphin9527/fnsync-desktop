package msg

import (
	"encoding/base64"
	"testing"
)

func TestDecrypt(t *testing.T) {
	code := "297e8fdf-0b77-4f2c-a508-43d1edc412023709bb57-fe94-441c-bf4e-90c61ddc25a4"
	encryptText, _ := base64.StdEncoding.DecodeString("w0c/4bKdOKW/TUV5C8XolgtMEXbZ1AHwA+cLQkMwQCfUmsw49VtTeXaSCD+YZIngn8NNngdL5KjPpgmA59fbAtCMbAahlLFyKnAEqKCp87dQ3jtfAVRGoSZKL2/PXXontBbmDHqMl6fcprpe938DnPdQn2U5Twrk9es6I8AdB2PCTX06ZjFIl0ieXZapt3SkXPlqhsi2sykF9VsiySsNKVrvkc5M532QXpvWgfW1vAYyJpLLIQeEQ38RkNY4mF1Da/DSkm+LIaeBbkOzNSgSM2IpG6IssSyGYeaffhFYRC51Zg9uICohZ02LjY4qOSfh78rrlu9XHtuRQQ3Ab6oEzKw=")
	expect := `{"data":"{\"phoneid\":\"461caca8-70c1-4e48-b753-38a32bf9a0e4\",\"oldconnection\":true,\"phonename\":\"GM1910\",\"wait_accept\":true}","hash":"80832ff0422797148e785af02b76649cfbf14320749de9f45ba1ad3a91459fcf","time":1633592823573}`

	encryMgr := NewEncryptionManager([]byte(code))
	result, err := encryMgr.Decrypt(encryptText)
	if err != nil {
		t.Error(err)
	}

	if string(result) != expect {
		t.Errorf("decrypt not correct.  result: %s  expect: %s", string(result), expect)
	}
}

package certs

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/daticahealth/cli/test"
)

const (
	certsCreateCommandName    = "certs"
	certsCreateSubcommandName = "create"
	certName                  = "example.com"
	pubKey                    = `-----BEGIN CERTIFICATE-----
MIICATCCAWoCCQCsoDP5n7FfzzANBgkqhkiG9w0BAQUFADBFMQswCQYDVQQGEwJB
VTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lkZ2l0
cyBQdHkgTHRkMB4XDTE1MDYwNDE5MTkzNVoXDTE2MDYwMzE5MTkzNVowRTELMAkG
A1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAfBgNVBAoTGEludGVybmV0
IFdpZGdpdHMgUHR5IEx0ZDCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA3+Gz
NFJhBdbUcFUxzlm70DJHXa9+nOAHZ9S6c66T1FXBRF94GfTSq8Qg9U+EOZf5cuhN
6wkLD1LLHMdb/UEjyCVVOqscfeR/nPCT5B9sv881PM8jL8C7grAUezcKiNx7Fng8
Dj9sczwziBR9P5ke5TI1g62LhHc0KGgMa8oNY7UCAwEAATANBgkqhkiG9w0BAQUF
AAOBgQBgTk8C+e13xGEw8qI2xhNfudt+8ffzIjNNWptb8rhGWblyY7EVBuU24LqE
oIOS7EH2aRhgvZjPUEQCNl+foQBRnRkYBeBhfUTl8QAUQNIyRUAHlQcPct9+VYcz
7OeuMetZkluMG3w62ooiufaGC/8orztDEySO4cj1HWssE2h/zw==
-----END CERTIFICATE-----`
	privKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDf4bM0UmEF1tRwVTHOWbvQMkddr36c4Adn1LpzrpPUVcFEX3gZ
9NKrxCD1T4Q5l/ly6E3rCQsPUsscx1v9QSPIJVU6qxx95H+c8JPkH2y/zzU8zyMv
wLuCsBR7NwqI3HsWeDwOP2xzPDOIFH0/mR7lMjWDrYuEdzQoaAxryg1jtQIDAQAB
AoGAeXoVqoYobuqqSmlvpO+7oLQnVQYsRSKp4gTjRnGrdMMzIs5KdIsK5Hh/CZwj
urxjdZ3m6Wj2v1HFM9BYcYouxx5ZYbUWx4tXeQhoVjvu8GxU6uwkDl+kQMjqcvfV
dXEoIm7ejzcvialYlHnsO8HFiB3ayhoQOK3kGcY6dGISWwECQQD/7R8/EIPAP0lU
P97w+I7j2kG79PTvCzoXygqVrmjeW6RJ6FvzT30iCnr5PVmPzHReL+q3i6tMHpGi
eeo0T0atAkEA3/I22OTH2QrKmSaW3EoPNDq78hJzsbSoVaHz+6mMn2ZungzBhJ7i
dOkUzkzuZtftYIcCQ2MtGDeSNIXuohOaKQJADwbVNta5ZahRnejCJlPxz98YzPht
CTwXhR4P0QoUjjnDQ7Oo8nhQWJdU8R1xDMhsbLtThMNmo2mIE4ok/j1JYQJBAJKg
pqSwduF3HVvVVmV54CaUZkaDKlkqLiWTWopmYvpjOP4m3/YTibZ+fe7tlBKmQng3
LZYts3Ltv77ACpT4PLECQQDDql4xPUb6WfsSjyqqfwnzkFLWADTcQQG5MmUX6iNJ
FBlcbW65DK1xPIitnX+jf803WaMPAP5YBoH6jC6VgcVH
-----END RSA PRIVATE KEY-----`
	certsCreateStandardOutput = `Created 'example.com'
To make use of your cert, you need to add a site with the "datica sites create" command
`
)

var (
	pubKeyPath  = "example.pem"
	privKeyPath = "example-key.pem"
	invalidPath = "invalid-file.pem"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if err := CreateCertFiles(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	statusCode := m.Run()
	CleanupCertFiles()
	os.Exit(statusCode)
}

var certCreateTests = []struct {
	env            string
	name           string
	pubKeyPath     string
	privKeyPath    string
	selfSigned     bool
	skipResolve    bool
	expectErr      bool
	expectedOutput string
}{
	{test.Alias, certName, pubKeyPath, privKeyPath, true, false, false, certsCreateStandardOutput},
	{test.Alias, certName, "./invalid-file.pem", privKeyPath, true, false, true, "\033[31m\033[1m[fatal] \033[0mA cert does not exist at path './invalid-file.pem'\n"},
	{test.Alias, certName, pubKeyPath, "./invalid-file.pem", true, false, true, "\033[31m\033[1m[fatal] \033[0mA private key does not exist at path './invalid-file.pem'\n"},
	{test.Alias, certName, pubKeyPath, privKeyPath, false, false, false, "Incomplete certificate chain found, attempting to resolve this\n" + certsCreateStandardOutput},
	{test.Alias, certName, pubKeyPath, privKeyPath, true, true, false, certsCreateStandardOutput},
	{"bad-env", certName, pubKeyPath, privKeyPath, true, false, true, "\033[31m\033[1m[fatal] \033[0mNo environment named \"bad-env\" has been associated. Run \"datica associated\" to see what environments have been associated or run \"datica associate\" from a local git repo to create a new association\n"},
}

func TestCertsCreate(t *testing.T) {
	if err := test.SetUpGitRepo(); err != nil {
		t.Error(err)
		return
	}
	if err := test.SetUpAssociation(); err != nil {
		t.Error(err)
		return
	}

	for _, data := range certCreateTests {
		t.Logf("Data: %+v", data)
		args := []string{"-E", data.env, certsCreateCommandName, certsCreateSubcommandName}
		if len(data.name) != 0 {
			args = append(args, data.name)
		}
		if len(data.pubKeyPath) != 0 {
			args = append(args, data.pubKeyPath)
		}
		if len(data.privKeyPath) != 0 {
			args = append(args, data.privKeyPath)
		}
		if data.selfSigned {
			args = append(args, "-s")
		}
		if data.skipResolve {
			args = append(args, "--resolve=false")
		}
		output, err := test.RunCommand(test.BinaryName, args)
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", output)
			continue
		}
		if output != data.expectedOutput {
			t.Errorf("Expected: %s. Found: %s", data.expectedOutput, output)
			continue
		}
		if err == nil {
			if output, err = test.RunCommand(test.BinaryName, []string{"-E", data.env, certsRmCommandName, certsRmSubcommandName, certName}); err != nil {
				t.Errorf("Unexpected err: %s", output)
				return
			}
		}
	}
}

func TestCertsCreateNoAssociation(t *testing.T) {
	if err := test.ClearAssociations(); err != nil {
		t.Error(err)
		return
	}
	output, err := test.RunCommand(test.BinaryName, []string{certsCreateCommandName, certsCreateSubcommandName, certName, pubKeyPath, privKeyPath})
	if err == nil {
		t.Errorf("Expected error but no error returned: %s", output)
		return
	}
	expectedOutput := "\033[31m\033[1m[fatal] \033[0mNo Datica environment has been associated. Run \"datica associate\" from a local git repo first\n"
	if output != expectedOutput {
		t.Errorf("Expected: %s. Found: %s", expectedOutput, output)
		return
	}
}

func TestCertsCreateTwice(t *testing.T) {
	if err := test.SetUpGitRepo(); err != nil {
		t.Error(err)
		return
	}
	if err := test.SetUpAssociation(); err != nil {
		t.Error(err)
		return
	}
	output, err := test.RunCommand(test.BinaryName, []string{"-E", test.Alias, certsCreateCommandName, certsCreateSubcommandName, certName, pubKeyPath, privKeyPath, "-s"})
	if err != nil {
		t.Errorf("Unexpected error: %s", output)
		return
	}
	if output != certsCreateStandardOutput {
		t.Errorf("Expected: %s. Found: %s", certsCreateStandardOutput, output)
		return
	}
	output, err = test.RunCommand(test.BinaryName, []string{"-E", test.Alias, certsCreateCommandName, certsCreateSubcommandName, certName, pubKeyPath, privKeyPath, "-s"})
	if err == nil {
		t.Errorf("Expected error but no error returned: %s", output)
		return
	}
	expectedOutput := "\033[31m\033[1m[fatal] \033[0m(92003) Cert Already Exists: A Cert already exists with the given name; alter name or delete existing cert to continue.\n"
	if output != expectedOutput {
		t.Errorf("Expected: %s. Found: %s", expectedOutput, output)
		return
	}
	if output, err = test.RunCommand(test.BinaryName, []string{"-E", test.Alias, certsRmCommandName, certsRmSubcommandName, certName}); err != nil {
		t.Errorf("Unexpected err: %s", output)
		return
	}
}

func CreateCertFiles() error {
	cert, err := os.OpenFile(pubKeyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer cert.Close()
	cert.WriteString(pubKey)
	key, err := os.OpenFile(privKeyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer key.Close()
	key.WriteString(privKey)
	return nil
}

func CleanupCertFiles() error {
	err := os.Remove(pubKeyPath)
	if err != nil {
		err = os.Remove(privKeyPath)
	}
	return err
}

package mega

import (
	"encoding/json"
	"fmt"
	MEGA "github.com/t3rm1n4l/go-mega"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// AccountCap amount of permanent upload space for mega account
	AccountCap = 15 * 1073741824
	// EncryptionFileType
	EncryptionFileType = ".maxcrypt"
	// AccountPath file path of stored mega credentials
	AccountPath = ".megaccount"
	TmpDir      = "/tmp/"
)

type Server struct {
	Host        string   `json:"host"`
	Port        int      `json:"ssh-port"`
	ExcludeDirs []string `json:"exclude-dirs"`
}

type Servers struct {
	Servers     map[string]Server `json:"servers"`
	TmpDir      string            `json:"tmp-dir"`
	ExcludeDirs []string          `json:"exclude-dirs"`
	Key         string            `json:"encryption-key"`
}

type CreateServer struct {
	Host        string
	Credentials string
}

type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (server CreateServer) BackupDirectory(dir, name, key string, account Account) (string, error) {
	tmpPath := TmpDir + name + EncryptionFileType
	if err := EncryptCompressDir(dir, tmpPath, key); err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	fi, err := os.Stat(tmpPath)
	if err != nil {
		return "", err
	}
	fileSize := uint64(fi.Size())

	if fileSize > AccountCap {
		// TODO split into smaller pieces
		return "", fmt.Errorf("file too large for upload")
	}

	m := MEGA.New()
	if err := m.Login(account.Email, account.Password); err != nil {
		log.Printf("Failed to log in with: %v\n", account)
		log.Printf("Deleting: %v\n", AccountPath)
		_ = os.Remove(AccountPath)
		return "", err
	}

	q, err := m.GetQuota()
	if err != nil {
		return "", err
	}

	used := q.Cstrg
	if AccountCap-used-fileSize <= 0 {
		log.Printf("%s:%s is full - creating new account!\n", account.Email, account.Password)
		// create new account
		account, err := server.CreateAccount()
		if err != nil {
			return "", err
		}

		// login with new account
		m = MEGA.New()
		if err := m.Login(account.Email, account.Password); err != nil {
			log.Println("Failed 2 login after create")
			return "", err
		}
	}

	log.Printf("Uploading %d bytes with %d bytes left ", fileSize, AccountCap-used)

	err = m.SetUploadWorkers(4)
	if err != nil {
		return "", err
	}

	dirNode, err := m.CreateDir(string(time.Now().UnixNano()), m.FS.GetRoot())
	if err != nil {
		return "", err
	}

	node, err := m.UploadFile(tmpPath, dirNode, "", nil)
	if err != nil {
		return "", err
	}

	link, err := m.Link(node, true)
	if err != nil {
		return "", err
	}
	return link, nil
}

func (server CreateServer) CreateAccount() (account Account, err error) {
	form := url.Values{}
	form.Add("credentials", server.Credentials)
	log.Println(server.Host + "/code")
	resp, err := http.PostForm(server.Host+"/code", form)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &account)
	if err != nil {
		log.Println(string(body))
		return
	}
	err = writeAccountToFile(account)
	return
}

func (server CreateServer) getStoredAccount() (account Account, err error) {
	bytes, _ := ioutil.ReadFile(AccountPath) // ignore missing file will be handled
	err = json.Unmarshal(bytes, &account)
	if err != nil {
		account, err = server.CreateAccount()
		if err != nil {
			return
		}
	}
	return
}

func writeAccountToFile(account Account) error {
	bytes, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return WriteFile(bytes, AccountPath)
}

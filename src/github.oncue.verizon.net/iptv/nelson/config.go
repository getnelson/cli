package main

import (
	"os"
	"log"
	"time"
	"errors"
	"net/http"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

///////////////////////////// CLI ENTRYPOINT //////////////////////////////////

func LoadDefaultConfig() (error, *Config) {
	pth := defaultConfigPath()
	empty := &Config {}

	if _, err := os.Stat(pth); os.IsNotExist(err) {
		return errors.New("No config file existed at " + pth + ". You need to `nelson login` before running other commands."), empty
	}

	err, parsed := readConfigFile(pth)

	if err != nil {
		return errors.New("Unable to read configuration file at '"+pth+"'. Reported error was: " + err.Error()), empty
	}

	return nil, parsed
}

////////////////////////////// CONFIG YAML ///////////////////////////////////

type Config struct {
	Endpoint      string `yaml:endpoint`
	ConfigSession `yaml:"session"`
}

type ConfigSession struct {
	Token     string `yaml:"token"`
	ExpiresAt int64  `yaml:"expires_at"`
}

func (c *Config) GetAuthCookie() *http.Cookie {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:       "nelson.session",
		Value:      c.ConfigSession.Token,
		Path:       "/",
		Domain:     "nelson-beta.oncue.verizon.net",
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     86400,
		Secure:     true,
		HttpOnly:   false,
	}

	return cookie
}

func generateConfigYaml(s Session, url string) string {
	temp := &Config{
		Endpoint: url,
		ConfigSession: ConfigSession{
			Token:     s.SessionToken,
			ExpiresAt: s.ExpiresAt,
		},
	}

	d, err := yaml.Marshal(&temp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return "---\n" + string(d)
}

func parseConfigYaml(yamlAsBytes []byte) *Config {
	temp := &Config{}
	err := yaml.Unmarshal(yamlAsBytes, &temp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return temp
}

func (c *Config) Validate() []error {
	// check that the token has not expired
	errs := []error{}

	if c.ConfigSession.ExpiresAt <= currentTimeMillis() {
		errs = append(errs, errors.New("Your session has expired. Please 'nelson login' again to reactivate your session."))
	}
	return errs
}

/////////////////////////////// CONFIG I/O ////////////////////////////////////

func defaultConfigPath() string {
	targetDir := os.Getenv("HOME") + "/.nelson"
	os.Mkdir(targetDir, 0775)
	return targetDir + "/config.yml"
}

// returns Unit, no error handling. YOLO
func writeConfigFile(s Session, url string, configPath string) {
	yamlConfig := generateConfigYaml(s, url)

	err := ioutil.WriteFile(configPath, []byte(yamlConfig), 0755)
	if err != nil {
		panic(err)
	}
}

func readConfigFile(configPath string) (error, *Config) {
	b, err := ioutil.ReadFile(configPath)
	return err, parseConfigYaml(b) // TIM: parsing never fails, right? ;-)
}

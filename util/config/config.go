package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
)

func LoadConfig(c interface{}) error {
	var fileLoc string
	if len(os.Args) < 2 {
		fileLoc = "conf/default.cfg"
	} else {
		fileLoc = os.Args[1]
	}

	configBytes, err := ioutil.ReadFile(fileLoc)
	if err != nil {
		return errors.Annotate(err, fmt.Sprintf("Error reading from file location %s", fileLoc))
	}

	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		return errors.Annotate(err, fmt.Sprintf("Error parsing config json from %s into config struct", fileLoc))
	}

	return nil
}

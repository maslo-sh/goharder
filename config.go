package prototype

import (
	"bufio"
	"embed"
	"io"
	"strings"
)

//go:embed resources
var resourcesDir embed.FS

type Config map[string]string

func ReadPropertiesBasedConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{}
	if len(filename) == 0 {
		return config, nil
	}
	file, err := resourcesDir.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

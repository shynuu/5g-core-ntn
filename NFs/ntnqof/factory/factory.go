/*
 * AMF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/shynuu/ntn-qof/logger"
)

var (
	QofConfig Config
)

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		QofConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &QofConfig); yamlErr != nil {
			return yamlErr
		}
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := QofConfig.GetVersion()

	if currentVersion != QOF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("QOF config version is [%s], but expected is [%s].",
			currentVersion, QOF_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("QOF config version [%s]", currentVersion)

	return nil
}

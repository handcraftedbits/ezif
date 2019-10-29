package main // import "golang.handcraftedbits.com/ezif/cmd/genhelper"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

//
// Private types
//

type genhelperConfig struct {
	DisabledAccessors []string                 `yaml:"disabledAccessors"`
	DisabledTests     []string                 `yaml:"disabledTests"`
	Packages          map[string]packageConfig `yaml:"packages"`
}

type packageConfig struct {
	Description string `yaml:"description"`
	Family      string `yaml:"family"`
	Groups      string `yaml:"groups"`
	Reference   string `yaml:"reference"`

	disabledAccessors map[string]bool
	disabledTests     map[string]bool
	groupsRegexp      *regexp.Regexp
}

//
// Private variables
//

var commandRoot = &cobra.Command{
	Use:   "genhelper",
	Short: "genhelper is used to generate ezif helper functions and test code from Exiv2 metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		var config genhelperConfig
		var contents []byte
		var err error
		var generatedSource string
		var metadata families
		var pc packageConfig

		// Validate arguments and run templates.

		contents, err = ioutil.ReadFile(flagMetadata)

		if err != nil {
			return errors.Wrap(err, "invalid Exiv2 metadata file provided")
		}

		err = json.Unmarshal(contents, &metadata)

		if err != nil {
			return errors.Wrap(err, "invalid Exiv2 metadata file provided")
		}

		contents, err = ioutil.ReadFile(flagConfig)

		if err != nil {
			return errors.Wrap(err, "invalid genhelper configuration file provided")
		}

		err = yaml.Unmarshal(contents, &config)

		if err != nil {
			return errors.Wrap(err, "invalid genhelper configuration file provided")
		}

		if _, ok := config.Packages[flagPackage]; !ok {
			return fmt.Errorf("invalid package '%s' specified", flagPackage)
		}

		pc = config.Packages[flagPackage]

		pc.disabledAccessors = getDisabledItems(config.DisabledAccessors)
		pc.disabledTests = getDisabledItems(config.DisabledTests)
		pc.groupsRegexp, err = regexp.Compile(pc.Groups)

		if err != nil {
			return errors.Wrap(err, "invalid groups regular expression provided")
		}

		if !flagPackageTest {
			generatedSource, err = generateGroupSource(pc.Family, metadata[pc.Family], flagPackage, pc)
		} else {
			generatedSource, err = generateGroupTestSource(pc.Family, metadata[pc.Family], flagPackage, pc)
		}

		if err != nil {
			return err
		}

		fmt.Print(generatedSource)

		return nil
	},
}

var (
	flagConfig      string
	flagMetadata    string
	flagPackage     string
	flagPackageTest bool
)

//
// Private functions
//

func getDisabledItems(items []string) map[string]bool {
	var result = make(map[string]bool)

	for _, item := range items {
		result[item] = true
	}

	return result
}

func main() {
	commandRoot.Flags().StringVarP(&flagConfig, "config", "c", "", "path to genhelper configuration file")
	commandRoot.Flags().StringVarP(&flagMetadata, "metadata", "m", "", "path to JSON-formatted Exiv2 metadata")
	commandRoot.Flags().StringVarP(&flagPackage, "package", "p", "", "name of package to use for code generation")
	commandRoot.Flags().BoolVarP(&flagPackageTest, "test", "t", false, "whether or not test code, instead of accessor "+
		"code, should be generated")

	_ = commandRoot.MarkFlagRequired("config")
	_ = commandRoot.MarkFlagRequired("metadata")
	_ = commandRoot.MarkFlagRequired("package")

	if err := commandRoot.Execute(); err != nil {
		log.New(os.Stderr, "", 0).Printf("%v\n", err)

		os.Exit(1)
	}
}

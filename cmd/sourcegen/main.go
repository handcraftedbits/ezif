package main // import "golang.handcraftedbits.com/ezif/cmd/sourcegen"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

//
// Private types
//

type sourcegenConfig struct {
	DisabledHelpers []string               `yaml:"disabledHelpers"`
	DisabledTests   []string               `yaml:"disabledTests"`
	Groups          map[string]groupConfig `yaml:"groups"`
}

//
// Private variables
//

var commandAccessor = &cobra.Command{
	Use:   "accessor",
	Short: "Generate metadata property accessors",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var generatedSource string
		var templateSource string

		if flagAccessorImpl {
			templateSource = templateAccessorImplSource
		} else {
			templateSource = templateAccessorIntfSource
		}

		generatedSource, err = generateAccessorSource(flagAccessorPackage, templateSource)

		if err != nil {
			return err
		}

		fmt.Print(generatedSource)

		return nil
	},
}

var commandHelper = &cobra.Command{
	Use:   "helper",
	Short: "Generate ezif helper functions and test code from Exiv2 metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		var config *sourcegenConfig
		var contents []byte
		var err error
		var gc groupConfig
		var generatedSource string
		var metadata families

		// Validate arguments and run templates.

		contents, err = ioutil.ReadFile(flagHelperMetadata)

		if err != nil {
			return errors.Wrap(err, "invalid Exiv2 metadata file provided")
		}

		err = json.Unmarshal(contents, &metadata)

		if err != nil {
			return errors.Wrap(err, "invalid Exiv2 metadata file provided")
		}

		config, err = getSourcegenConfig()

		if err != nil {
			return err
		}

		if _, ok := config.Groups[flagHelperGroup]; !ok {
			return fmt.Errorf("invalid group '%s' specified", flagHelperGroup)
		}

		gc = config.Groups[flagHelperGroup]

		gc.disabledHelpers = getDisabledItems(config.DisabledHelpers)
		gc.disabledTests = getDisabledItems(config.DisabledTests)
		gc.regexp, err = regexp.Compile(gc.Regexp)

		if err != nil {
			return errors.Wrap(err, "invalid groups regular expression provided")
		}

		if !flagHelperTest {
			generatedSource, err = generateGroupSource(gc.Family, metadata[gc.Family], flagHelperGroup, gc)
		} else {
			generatedSource, err = generateGroupTestSource(gc.Family, metadata[gc.Family], flagHelperGroup, gc)
		}

		if err != nil {
			return err
		}

		fmt.Print(generatedSource)

		return nil
	},
}

var commandRoot = &cobra.Command{
	Use:   "sourcegen",
	Short: "sourcegen is used to generate ezif helper and accessor code",
}

var (
	flagAccessorImpl    bool
	flagAccessorPackage string
	flagConfig          string
	flagHelperGroup     string
	flagHelperMetadata  string
	flagHelperTest      bool
)

//
// Private functions
//

func getSourcegenConfig() (*sourcegenConfig, error) {
	var contents []byte
	var err error
	var config sourcegenConfig

	contents, err = ioutil.ReadFile(flagConfig)

	if err != nil {
		return nil, errors.Wrap(err, "invalid sourcegen configuration file provided")
	}

	err = yaml.Unmarshal(contents, &config)

	if err != nil {
		return nil, errors.Wrap(err, "invalid sourcegen configuration file provided")
	}

	return &config, nil
}

func getDisabledItems(items []string) map[string]bool {
	var result = make(map[string]bool)

	for _, item := range items {
		result[item] = true
	}

	return result
}

func main() {
	commandAccessor.Flags().BoolVarP(&flagAccessorImpl, "impl", "i", false, "whether or not implementation code, "+
		"instead of interface code, should be generated")
	commandAccessor.Flags().StringVarP(&flagAccessorPackage, "package", "p", "", "the package for generated accessor "+
		"code")
	commandHelper.Flags().StringVarP(&flagConfig, "config", "c", "", "path to sourcegen configuration file")
	commandHelper.Flags().StringVarP(&flagHelperGroup, "group", "g", "", "name of group to use for code generation")
	commandHelper.Flags().StringVarP(&flagHelperMetadata, "metadata", "m", "", "path to JSON-formatted Exiv2 metadata")
	commandHelper.Flags().BoolVarP(&flagHelperTest, "test", "t", false, "whether or not test code, instead of helper "+
		"code, should be generated")

	_ = commandAccessor.MarkFlagRequired("package")
	_ = commandHelper.MarkFlagRequired("config")
	_ = commandHelper.MarkFlagRequired("group")
	_ = commandHelper.MarkFlagRequired("metadata")

	commandRoot.AddCommand(commandAccessor, commandHelper)

	if err := commandRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

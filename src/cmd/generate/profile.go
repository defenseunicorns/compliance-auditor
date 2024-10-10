package generate

import (
	"fmt"
	"strings"

	pkgCommon "github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var profileExample = `
To generate a profile with included controls:
	lula generate profile -s <catalog/profile source> -i ac-1,ac-2,ac-3

To specify the name and filetype of the generated artifact:
	lula generate profile -s <catalog/profile source> -i ac-1,ac-2,ac-3 -o my_profile.yaml

To generate a profile that includes all controls except a list specified controls:
	lula generate profile -s <catalog/profile source> -e ac-1,ac-2,ac-3
`

var profileLong = `Generation of a Profile OSCAL artifact with controls included or excluded from a source catalog/profile.`

func GenerateProfileCommand() *cobra.Command {
	var (
		source     string
		outputFile string
		include    []string
		exclude    []string
	)

	profilecmd := &cobra.Command{
		Use:     "profile",
		Aliases: []string{"p"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Generate a profile OSCAL template",
		Long:    profileLong,
		Example: profileExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			message.Info("generate profile executed")

			if outputFile == "" {
				outputFile = "profile.yaml"
			}

			// pre-check if the output file exists
			exists, err := pkgCommon.CheckFileExists(outputFile)
			if err != nil {
				return err
			}

			if exists {
				return fmt.Errorf("output File %s currently exist - cannot merge artifacts", outputFile)
			}

			command := fmt.Sprintf("%s --source %s", cmd.CommandPath(), source)

			if len(include) > 0 {
				command += fmt.Sprintf(" --include %s", strings.Join(include, ","))
			}

			if len(exclude) > 0 {
				command += fmt.Sprintf(" --exclude %s", strings.Join(exclude, ","))
			}

			profile, err := oscal.GenerateProfile(command, source, include, exclude)
			if err != nil {
				return err
			}

			// Write the component definition to file
			err = oscal.WriteOscalModelNew(outputFile, profile)
			if err != nil {
				message.Fatalf(err, "error writing component to file")
			}

			return nil
		},
	}

	profilecmd.Flags().StringVarP(&source, "source", "s", "", "the path to the source catalog/profile")
	profilecmd.MarkFlagRequired("source")
	profilecmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be directed to stdout")
	profilecmd.Flags().StringSliceVarP(&include, "include", "i", []string{}, "comma delimited list of controls to include from the source catalog/profile")
	profilecmd.Flags().StringSliceVarP(&exclude, "exclude", "e", []string{}, "comma delimited list of controls to exclude from the source catalog/profile")
	profilecmd.MarkFlagsMutuallyExclusive("include", "exclude")

	return profilecmd
}

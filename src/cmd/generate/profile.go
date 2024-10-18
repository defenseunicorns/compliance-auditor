package generate

import (
	"fmt"
	"strings"

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
		all        bool
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

			/// Check if output file contains a valid OSCAL model
			_, err := oscal.ValidOSCALModelAtPath(outputFile)
			if err != nil {
				return fmt.Errorf("invalid OSCAL model at output: %v", err)
			}

			command := fmt.Sprintf("%s --source %s", cmd.CommandPath(), source)

			if len(include) > 0 {
				command += fmt.Sprintf(" --include %s", strings.Join(include, ","))
			}

			if len(exclude) > 0 {
				command += fmt.Sprintf(" --exclude %s", strings.Join(exclude, ","))
			}

			if all {
				command += " --all"
			}

			profile, err := oscal.GenerateProfile(command, source, include, exclude, all)
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
	err := profilecmd.MarkFlagRequired("source")
	if err != nil {
		message.Fatal(err, "error initializing upgrade command flags")
	}
	profilecmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be directed to stdout")
	profilecmd.Flags().StringSliceVarP(&include, "include", "i", []string{}, "comma delimited list of controls to include from the source catalog/profile")
	profilecmd.Flags().StringSliceVarP(&exclude, "exclude", "e", []string{}, "comma delimited list of controls to exclude from the source catalog/profile")
	profilecmd.Flags().BoolVarP(&all, "all", "a", false, "Include all controls from the source catalog/profile")
	profilecmd.MarkFlagsMutuallyExclusive("include", "exclude", "all")
	profilecmd.MarkFlagsOneRequired("include", "exclude", "all")

	return profilecmd
}

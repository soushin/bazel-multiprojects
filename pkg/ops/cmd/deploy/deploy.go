package deploy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

type buildOptions struct {
	packageName string
}

const rootPath string = "#"

func NewCmdDeploy(out io.Writer) *cobra.Command {

	var o buildOptions

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy is command for deploying app",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

			if err := o.RunBuild(out); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&o.packageName, "package", "p", "", "Package name")
	cmd.MarkFlagRequired("package")

	return cmd
}

func (o *buildOptions) RunBuild(out io.Writer) error {

	resourcePath := fmt.Sprintf("%s/pkg/%s/k8s/base", rootPath, o.packageName)

	yamlBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/kustomization.yaml", resourcePath))
	if err != nil {
		return err
	}

	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		return err
	}

	//kustomization, err := json.Marshal(jsonBytes)
	//if err != nil {
	//	return err
	//}

	out.Write(jsonBytes)

	return nil
}

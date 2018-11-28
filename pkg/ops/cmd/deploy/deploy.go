package deploy

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/k8sdeps"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/target"
)

type buildOptions struct {
	packageName string
	commitHash  string
}

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
	cmd.Flags().StringVarP(&o.commitHash, "commitHash", "c", "", "Commit hash")
	cmd.MarkFlagRequired("package")

	return cmd
}

func (o *buildOptions) RunBuild(out io.Writer) error {

	var kustomizationPath = "https://github.com/soushin/bazel-multiprojects"
	if o.commitHash != "" {
		kustomizationPath = fmt.Sprintf("%s?ref=%s", kustomizationPath, o.commitHash)
	}

	fSys := fs.MakeRealFS()
	ldr, err := loader.NewLoader(kustomizationPath, fSys)
	if err != nil {
		return err
	}

	nldr, err := ldr.New(fmt.Sprintf("pkg/%s/k8s", o.packageName))
	if err != nil {
		return err
	}

	f := k8sdeps.NewFactory()
	kt, err := target.NewKustTarget(nldr, fSys, f.ResmapF, f.TransformerF)
	if err != nil {
		return err
	}

	allResources, err := kt.MakeCustomizedResMap()
	if err != nil {
		return err
	}

	res, err := allResources.EncodeAsYaml()
	if err != nil {
		return err
	}

	out.Write(res)

	//in := os.Stdin
	//in.Write(res)
	//
	//ioErr := os.Stderr
	//
	//ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: ioErr}
	//
	//apply := kubeCmd.NewApplyOptions(ioStreams)
	//if err := apply.Run(); err != nil {
	//	return err
	//}

	return nil
}

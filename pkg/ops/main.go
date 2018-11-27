package main

import (
	"log"
	"os"

	"sigs.k8s.io/kustomize/k8sdeps"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/target"
)

func main() {

	kustomizationPath := "git@github.com:kubernetes-sigs/kustomize.git"

	fSys := fs.MakeRealFS()
	ldr, err := loader.NewLoader(kustomizationPath, fSys)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	nldr, err := ldr.New("examples/helloWorld")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	f := k8sdeps.NewFactory()
	kt, err := target.NewKustTarget(nldr, fSys, f.ResmapF, f.TransformerF)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	allResources, err := kt.MakeCustomizedResMap()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	res, err := allResources.EncodeAsYaml()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	stdOut := os.Stdout
	stdOut.Write(res)
}

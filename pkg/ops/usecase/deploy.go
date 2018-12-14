package usecase

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/ghodss/yaml"

	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/pkg/ops/client"
	"go.uber.org/zap"
)

type DeployUseCase interface {
	GetContents(owner, repo, path string) ([]string, error)
	ExistsContent(owner, repo, path string) error
	ExistsBranch(owner, repo, branch string) error
	CheckoutBranch(owner, repo, branch string) (string, error)
	ReplaceImage(checkoutPath, packagePath, owner, repo, branch string) error
}

const (
	K8s_PATH = "k8s/overlays/default"
)

type deployUseCaseImpl struct {
	appLog    *zap.Logger
	githubCli client.GitHubClient
}

func NewDeployUseCase(appLog *zap.Logger, githubCli client.GitHubClient) DeployUseCase {
	return &deployUseCaseImpl{
		appLog:    appLog,
		githubCli: githubCli,
	}
}

func (u *deployUseCaseImpl) GetContents(owner, repo, path string) ([]string, error) {

	contents, err := u.githubCli.GetContents(owner, repo, path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get from github")
	}

	targets := make([]string, len(contents))
	for i, content := range contents {
		targets[i] = *content.Path
	}

	return targets, nil
}

func (u *deployUseCaseImpl) ExistsContent(owner, repo, path string) error {

	_, err := u.githubCli.GetContents(owner, repo, fmt.Sprintf("%s/%s", path, K8s_PATH))
	if err != nil {
		return errors.Wrap(err, "failed to get from github")
	}

	return nil
}

func (h *deployUseCaseImpl) ExistsBranch(owner, repo, branch string) error {

	//if _, err := h.githubCli.GetBranch(owner, repo, branch); err != nil {
	//	return errors.Wrap(err, "failed to get from github")
	//}

	return nil
}

func (u *deployUseCaseImpl) CheckoutBranch(owner, repo, branch string) (string, error) {

	checkOutDir := fmt.Sprintf("/tmp/deploy/%s", fmt.Sprintf("%s_%s", time.Now().Format("20060102150405"), branch))
	fullPath := fmt.Sprintf("git@github.com:%s/%s", owner, repo)

	cmd := exec.Command("git", "clone", fullPath, "-b", branch, checkOutDir)
	if _, err := u.runCmdOut(cmd); err != nil {
		return "", errors.Wrap(err, "failed to get from github")
	}

	return checkOutDir, nil
}

func (u *deployUseCaseImpl) ReplaceImage(checkoutPath, packagePath, owner, repo, branch string) error {

	deploymentPath := fmt.Sprintf("%s/%s/%s/deployment-patch.yaml", checkoutPath, packagePath, K8s_PATH)
	data, err := ioutil.ReadFile(deploymentPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", deploymentPath)
	}

	patch := make(map[string]interface{})
	yaml.Unmarshal(data, &patch)

	specContainer, ok := u.getValue(patch, "spec.template.spec.containers")
	if !ok {
		u.appLog.With(zap.Any("deploymentPatch", patch)).Error("invalid process")
		return errors.Wrap(err, "failed to get spec of container")
	}
	containers := specContainer.([]interface{})

	var appName string
	for i, v := range strings.Split(packagePath, "/") {
		if i == 1 {
			appName = v
		}
	}

	var originalImage = ""
	for _, container := range containers {
		c := container.(map[string]interface{})
		if c["name"] == appName {
			originalImage = c["image"].(string)
		}
	}
	if originalImage == "" {
		u.appLog.With(zap.Any("deploymentPatch", patch)).Error("invalid process")
		return errors.Wrap(err, "failed to get spec of image")
	}

	var tag = ""
	if branch == "master" {
		tag = "latest"
	} else {
		r := strings.NewReplacer("/", "_")
		tag = r.Replace(branch)
	}

	replaceImage := fmt.Sprintf("%s/%s:%s", owner, repo, tag)
	replaceData := strings.Replace(string(data), originalImage, replaceImage, 1)

	u.appLog.With(zap.String("replaceData", replaceData)).Info("debug")

	return nil
}

func (u *deployUseCaseImpl) getValue(m map[string]interface{}, key string) (interface{}, bool) {
	for _, k := range strings.Split(key, ".") {
		if v, ok := m[k]; ok {
			switch v.(type) {
			case map[string]interface{}:
				m = v.(map[string]interface{})
			default:
				return v, true
			}
		} else {
			return nil, false
		}
	}
	return m, true
}

func (u *deployUseCaseImpl) runCmdOut(cmd *exec.Cmd) ([]byte, error) {
	u.appLog.With(zap.Strings("args", cmd.Args)).Info("Running command")
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrapf(err, "starting command %v", cmd)
	}

	stdout, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return nil, err
	}

	stderr, err := ioutil.ReadAll(stderrPipe)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return stdout, errors.Wrapf(err, "Running %s: stdout %s, stderr: %s, err: %v", cmd.Args, stdout, stderr, err)
	}

	if len(stderr) > 0 {
		u.appLog.With(zap.String("out", string(stdout)),
			zap.String("err", string(stdout))).Info("Command output")
	} else {
		u.appLog.With(zap.String("out", string(stdout))).Info("Command output")
	}

	return stdout, nil
}

package stagezt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hanzozt/fablab/kernel/model"
	"github.com/hanzozt/zt/v2/common/getzt"
	"github.com/hanzozt/zt/v2/zt/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func StageZitiOnce(run model.Run, component *model.Component, version string, source string) error {
	op := "install.zt-"
	if version == "" {
		op += "local"
	} else {
		op += version
	}

	return run.DoOnce(op, func() error {
		return StageZiti(run, component, version, source)
	})
}

func StageZrokOnce(run model.Run, component *model.Component, version string, source string) error {
	op := "install.zrok-"
	if version == "" {
		op += "local"
	} else {
		op += version
	}

	return run.DoOnce(op, func() error {
		return StageZrok(run, component, version, source)
	})
}

func StageCaddyOnce(run model.Run, component *model.Component, version string, source string) error {
	op := "install.caddy-"
	if version == "" {
		op += "local"
	} else {
		op += version
	}

	return run.DoOnce(op, func() error {
		return StageCaddy(run, component, version, source)
	})
}

func StageZitiEdgeTunnelOnce(run model.Run, component *model.Component, version string, source string) error {
	op := "install.zt-edge-tunnel-"
	if version == "" {
		op += "local"
	} else {
		op += version
	}

	return run.DoOnce(op, func() error {
		return StageZitiEdgeTunnel(run, component, version, source)
	})
}

func StageZiti(run model.Run, component *model.Component, version string, source string) error {
	return StageExecutable(run, "zt", component, version, source, func() error {
		return getzt.InstallZiti(version, "linux", "amd64", run.GetBinDir(), false)
	})
}

func StageZrok(run model.Run, component *model.Component, version string, source string) error {
	return StageExecutable(run, "zrok", component, version, source, func() error {
		return getzt.InstallZrok(version, "linux", "amd64", run.GetBinDir(), false)
	})
}

func StageCaddy(run model.Run, component *model.Component, version string, source string) error {
	return StageExecutable(run, "caddy", component, version, source, func() error {
		return getzt.InstallCaddy(version, "linux", "amd64", run.GetBinDir(), false)
	})
}

func StageLocalOnce(run model.Run, executable string, component *model.Component, source string) error {
	op := fmt.Sprintf("install.%s-local", executable)
	return run.DoOnce(op, func() error {
		return StageExecutable(run, executable, component, "", source, func() error {
			return fmt.Errorf("unable to fetch %s, as it a local-only application", executable)
		})
	})
}

func StageExecutable(run model.Run, executable string, component *model.Component, version string, source string, fallbackF func() error) error {
	fileName := executable
	if version != "" {
		fileName += "-" + version
	}

	target := filepath.Join(run.GetBinDir(), fileName)
	if version == "" || version == "latest" {
		_ = os.Remove(target)
	}

	envVar := strings.ToUpper(executable) + "_PATH"

	if version == "" {
		if source != "" {
			logrus.Infof("[%s] => [%s]", source, target)
			return util.CopyFile(source, target)
		}
		if envSource, found := component.GetStringVariable(envVar); found {
			logrus.Infof("[%s] => [%s]", envSource, target)
			return util.CopyFile(envSource, target)
		}
		if ztPath, err := exec.LookPath(executable); err == nil {
			logrus.Infof("[%s] => [%s]", ztPath, target)
			return util.CopyFile(ztPath, target)
		}
		return fmt.Errorf("%s binary not found in path, no path provided and no %s env variable set", executable, envVar)
	}

	found, err := run.FileExists(filepath.Join(model.BuildKitDir, model.BuildBinDir, fileName))
	if err != nil {
		return err
	}

	if found {
		logrus.Infof("%s already present, not downloading again", target)
		return nil
	}

	logrus.Infof("%s not present, attempting to fetch", target)

	return fallbackF()
}

func StageZitiEdgeTunnel(run model.Run, component *model.Component, version string, source string) error {
	fileName := "zt-edge-tunnel"
	if version != "" {
		fileName += "-" + version
	}

	target := filepath.Join(run.GetBinDir(), fileName)
	if version == "" || version == "latest" {
		_ = os.Remove(target)
	}

	if version == "" {
		if source != "" {
			logrus.Infof("[%s] => [%s]", source, target)
			return util.CopyFile(source, target)
		}
		if envSource, found := component.GetStringVariable("zt-edge-tunnel.path"); found {
			logrus.Infof("[%s] => [%s]", envSource, target)
			return util.CopyFile(envSource, target)
		}
		if ztPath, err := exec.LookPath("zt-edge-tunnel"); err == nil {
			logrus.Infof("[%s] => [%s]", ztPath, target)
			return util.CopyFile(ztPath, target)
		}
		return errors.New("zt-edge-tunnel binary not found in path, no path provided and no zt-edge-tunnel.path env variable set")
	}

	found, err := run.FileExists(filepath.Join(model.BuildKitDir, model.BuildBinDir, fileName))
	if err != nil {
		return err
	}

	if found {
		logrus.Infof("%s already present, not downloading again", target)
		return nil
	}
	logrus.Infof("%s not present, attempting to fetch", target)

	return getzt.InstallZitiEdgeTunnel(version, "linux", "amd64", run.GetBinDir(), false)
}

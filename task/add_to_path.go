package task

import (
	"golang.org/x/sys/windows/registry"
	"strings"
)

func AddToPathEnv(add string) Task {
	return &addToPath{add}
}

type addToPath struct {
	add string
}

func (t *addToPath) Name() string {
	return "Adding " + t.add + " to PATH environment variable"
}

func (t *addToPath) Execute() error {
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		"Environment",
		registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer k.Close()

	path, _, err := k.GetStringValue("PATH")
	if err != nil {
		return err
	}

	paths := strings.Split(path, ";")
	for _, p := range paths {
		// remove white space and quotes around path
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, `"`) && strings.HasSuffix(p, `"`) {
			p = strings.TrimSuffix(p[1:], `"`)
		}

		if p == t.add {
			// the path to be added is already there => we are done
			return nil
		}
	}

	return k.SetStringValue("PATH", path+";"+t.add)
}

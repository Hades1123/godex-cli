package runtime

import "os/exec"

var nodeRoots = []string{
	"~/.nvm/versions/node",
	"~/.fnm/node-versions",
	"~/.local/share/fnm/node-versions",
}

func ListNode() ([]Install, error) {
	return listDirs(nodeRoots...)
}

func FindNode(query string) (Install, error) {
	installs, err := ListNode()
	if err != nil {
		return Install{}, err
	}
	return findInstall(query, installs)
}

func CurrentNode() (string, error) {
	out, err := exec.Command("node", "--version").CombinedOutput()
	return string(out), err
}

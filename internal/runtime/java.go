package runtime

import "os/exec"

var javaRoots = []string{
	"/usr/lib/jvm",
	"~/.sdkman/candidates/java",
	"~/.jenv/versions",
}

func ListJava() ([]Install, error) {
	return listDirs(javaRoots...)
}

func FindJava(query string) (Install, error) {
	installs, err := ListJava()
	if err != nil {
		return Install{}, err
	}
	return findInstall(query, installs)
}

func CurrentJava() (string, error) {
	out, err := exec.Command("java", "-version").CombinedOutput()
	return string(out), err
}

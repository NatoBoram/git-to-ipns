package main

import (
	"fmt"
	"net/url"
	"os/exec"

	"github.com/logrusorgru/aurora"
)

func command(cmd *exec.Cmd, errMessage string, cmdMessage ...interface{}) (out []byte, err error) {
	out, err = cmd.Output()
	if err != nil {
		fmt.Println(errMessage)
		fmt.Println(cmdMessage...)

		// Log the `ExitError`
		ee, ok := err.(*exec.ExitError)
		if ok {
			fmt.Println(string(ee.Stderr))
		}

		fmt.Println(string(out))
		return
	}
	return
}

func ipfsClusterAdd(link string, rmin string, rmax string, uuid string) (out []byte, err error) {
	return command(
		exec.Command("ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--cid-version=1", "--name", link, "--replication-min", rmin, "--replication-max", rmax, uuid),
		"Couldn't add the repository to IPFS.",
		aurora.Bold("Command :"), "ipfs-cluster-ctl", "add", "--recursive", "--quieter", "--chunker=rabin", "--cid-version=1", "--name", aurora.Blue(link), "--replication-min", aurora.Bold(rmin), "--replication-max", aurora.Bold(rmax), uuid,
	)
}

func ipfsKeyGen(link string) (out []byte, err error) {
	escaped := url.PathEscape(link)
	return command(
		exec.Command("ipfs", "key", "gen", "--type", "ed25519", escaped),
		"Couldn't generate a new key.",
		aurora.Bold("Command :"), "ipfs", "key", "gen", "--type", "ed25519", aurora.Blue(escaped),
	)
}

func ipfsNamePublish(key string, ipfs string) (out []byte, err error) {
	return command(
		exec.Command("ipfs", "name", "publish", "--key", key, "--quieter", "/ipfs/"+ipfs),
		"Couldn't publish on IPNS.",
		aurora.Bold("Command :"), "ipfs", "name", "publish", "--key", key, "--quieter", aurora.Cyan("/ipfs/"+ipfs),
	)
}

func gitClone(link string, uuid string) (out []byte, err error) {
	return command(
		exec.Command("git", "clone", link, uuid),
		"Couldn't clone the repository.",
		aurora.Bold("Command :"), "git", "clone", aurora.Blue(link), uuid,
	)
}

func gitPull(uuid string) (out []byte, err error) {
	return command(
		exec.Command("git", "-C", uuid, "pull"),
		"Couldn't pull the repository.",
		aurora.Bold("Command :"), "git", "-C", uuid, "pull",
	)
}

func ipfsClusterRm(ipfs string) (out []byte, err error) {
	return command(
		exec.Command("ipfs-cluster-ctl", "pin", "rm", ipfs),
		"Couldn't remove the repository from IPFS.",
		aurora.Bold("Command :"), "ipfs-cluster-ctl", "pin", "rm", aurora.Cyan(ipfs),
	)
}

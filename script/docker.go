/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package script

import (
	"fmt"
	"github.com/yimincai/toolbox/utils"
)

func UbuntuInstallDocker() {
	pkg.PrintConsole("Install Docker ... ")

	pkg.RunCommand("sudo apt-get update -y", false)
	pkg.RunCommand("sudo apt-get install ca-certificates curl gnupg -y", false)
	pkg.RunCommand("sudo install -m 0755 -d /etc/apt/keyrings", false)
	pkg.RunCommand("curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg", false)
	pkg.RunCommand("sudo chmod a+r /etc/apt/keyrings/docker.gpg", false)
	pkg.RunCommand("echo \\\n  \"deb [arch=\"$(dpkg --print-architecture)\" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \\\n  \"$(. /etc/os-release && echo \"$VERSION_CODENAME\")\" stable\" | \\\n  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null", false)
	pkg.RunCommand("sudo apt-get update -y", false)
	pkg.RunCommand("sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y", false)
	pkg.RunCommand("sudo groupadd -f docker", false)
	// Get real username
	user := utils.GetUsername()
	pkg.RunCommand(fmt.Sprintf("sudo usermod -aG docker %v", user), false)
	pkg.RunCommand("sudo newgrp docker", false)
	pkg.RunCommand("sudo groups", false)
	pkg.RunCommand("source ~/.bashrc", false)
	pkg.RunCommand("sudo docker version", false)

	pkg.PrintNotice("Please re-login to leave docker permission error!")
	pkg.PrintConsole("Docker installed successfully!")
}

/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package script

import (
	"fmt"
	"github.com/yimincai/toolbox/utils"

	"github.com/yimincai/toolbox/tools"
)

func UbuntuInstallDocker() {
	tools.PrintConsole("Install Docker ... ")

	tools.RunCommand("sudo apt-get update -y", false)
	tools.RunCommand("sudo apt-get install ca-certificates curl gnupg -y", false)
	tools.RunCommand("sudo install -m 0755 -d /etc/apt/keyrings", false)
	tools.RunCommand("curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg", false)
	tools.RunCommand("sudo chmod a+r /etc/apt/keyrings/docker.gpg", false)
	tools.RunCommand("echo \\\n  \"deb [arch=\"$(dpkg --print-architecture)\" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \\\n  \"$(. /etc/os-release && echo \"$VERSION_CODENAME\")\" stable\" | \\\n  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null", false)
	tools.RunCommand("sudo apt-get update -y", false)
	tools.RunCommand("sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y", false)
	tools.RunCommand("sudo groupadd -f docker", false)
	// Get real username
	user := utils.GetUsername()
	tools.RunCommand(fmt.Sprintf("sudo usermod -aG docker %v", user), false)
	tools.RunCommand("sudo newgrp docker", false)
	tools.RunCommand("sudo groups", false)
	tools.RunCommand("source ~/.bashrc", false)
	tools.RunCommand("sudo docker version", false)

	tools.PrintNotice("Please re-login to leave docker permission error!")
	tools.PrintConsole("Docker installed successfully!")
}

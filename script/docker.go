/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package script

import (
	"fmt"

	"github.com/yimincai/toolbox/internal/command"
	"github.com/yimincai/toolbox/pkg/logger"
	"github.com/yimincai/toolbox/utils"
)

// UbuntuInstallDocker installs docker on ubuntu system
func UbuntuInstallDocker() {
	logger.Console("Install Docker ... ")

	command.Run("sudo apt-get update -y", false)
	command.Run("sudo apt-get install ca-certificates curl gnupg -y", false)
	command.Run("sudo install -m 0755 -d /etc/apt/keyrings", false)
	command.Run("curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg", false)
	command.Run("sudo chmod a+r /etc/apt/keyrings/docker.gpg", false)
	command.Run("echo \\\n  \"deb [arch=\"$(dpkg --print-architecture)\" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \\\n  \"$(. /etc/os-release && echo \"$VERSION_CODENAME\")\" stable\" | \\\n  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null", false)
	command.Run("sudo apt-get update -y", false)
	command.Run("sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y", false)
	command.Run("sudo groupadd -f docker", false)
	// Get real username
	user := utils.GetUsername()
	command.Run(fmt.Sprintf("sudo usermod -aG docker %v", user), false)
	command.Run("sudo newgrp docker", false)
	command.Run("sudo groups", false)
	command.Run("source ~/.bashrc", false)
	command.Run("sudo docker version", false)

	logger.Yellow("Please re-login to leave docker permission error!")
	logger.Console("Docker installed successfully!")
}

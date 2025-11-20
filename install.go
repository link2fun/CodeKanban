package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

type program struct{}

// Start will be called by the service manager, run the actual startup logic in a separate goroutine to avoid blocking.
func (p *program) Start(s service.Service) error {
	runningAsService = true
	go func() {
		run(false, "", 0) // forceMigrate=false, bind="", port=0 (use defaults)
	}()
	return nil
}

// Stop currently requires no additional cleanup, just return.
func (p *program) Stop(service.Service) error {
	return nil
}

// serviceInstall installs or uninstalls the Windows service based on user instruction.
func serviceInstall(isInstall bool) {
	cwd, _ := os.Getwd()
	wd, _ := filepath.Abs(cwd)

	svcConfig := &service.Config{
		Name:             "codekanban",
		DisplayName:      "Code Kanban Service",
		Description:      "Code Kanban runs automatically as a Windows service",
		WorkingDirectory: wd,
	}

	srv, err := service.New(&program{}, svcConfig)
	if err != nil {
		fmt.Printf("Failed to create service: %v\n", err)
		return
	}

	if isInstall {
		fmt.Println("Installing system service...")
		if err := srv.Install(); err != nil {
			fmt.Printf("Installation failed: %v\n", err)
			return
		}
		fmt.Println("Installation completed. Start or stop via the service manager.")
		return
	}

	fmt.Println("Uninstalling system service...")
	if err := srv.Stop(); err != nil {
		fmt.Printf("Failed to stop service: %v\n", err)
	}
	if err := srv.Uninstall(); err != nil {
		fmt.Printf("Uninstallation failed: %v\n", err)
		return
	}
	fmt.Println("Uninstallation completed.")
}

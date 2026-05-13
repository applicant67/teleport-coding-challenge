//go:build linux

package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	geteuid = os.Geteuid
	uname   = unix.Uname
	statfs  = unix.Statfs
)

// SetupCgroups verifies the environment and configures the cgroup root for delegation.
func SetupCgroups() error {
	if geteuid() != 0 {
		return fmt.Errorf("job worker service requires root privileges to configure cgroups")
	}
	if err := checkCgroupMount(CgroupMountPoint); err != nil {
		return err
	}
	if err := configureControllers(CgroupMountPoint); err != nil {
		return err
	}
	return nil
}

func checkCgroupMount(mountPoint string) error {
	var stat unix.Statfs_t
	if err := statfs(mountPoint, &stat); err != nil {
		return fmt.Errorf("failed to stat %s: %w", mountPoint, err)
	}
	// CGROUP2_SUPER_MAGIC is 0x63677270
	if stat.Type != unix.CGROUP2_SUPER_MAGIC {
		return fmt.Errorf("cgroup v2 not mounted at %s", mountPoint)
	}
	return nil
}

func configureControllers(mountPoint string) error {
	rootControllers, err := os.ReadFile(filepath.Join(mountPoint, cgroupControllersFile))
	if err != nil {
		return fmt.Errorf("failed to read root controllers: %w", err)
	}
	sControllers := string(rootControllers)
	// Use a map for exact matching instead of loose string contains
	available := make(map[string]bool)
	for _, c := range strings.Fields(sControllers) {
		available[c] = true
	}
	required := []string{"cpu", "memory", "io"}
	for _, ctrl := range required {
		if !available[ctrl] {
			return fmt.Errorf("missing controller %s in %s", ctrl, filepath.Join(mountPoint, cgroupControllersFile))
		}
	}

	// Enable controllers at the root level
	if err := writeCgroupFile(mountPoint, cgroupSubtreeControlFile, "+"+strings.Join(required, " +")); err != nil {
		return fmt.Errorf("failed to enable root controllers: %w", err)
	}

	// Create Worker Directory
	workerRoot := filepath.Join(mountPoint, CgroupWorkerDir)
	if err := os.MkdirAll(workerRoot, dirMode); err != nil {
		return fmt.Errorf("failed to create worker cgroup root %s: %w", workerRoot, err)
	}

	// Enable Controllers in Worker Directory
	if err := writeCgroupFile(workerRoot, cgroupSubtreeControlFile, "+"+strings.Join(required, " +")); err != nil {
		return fmt.Errorf("failed to enable controllers in %s: %w", workerRoot, err)
	}

	// Feature Detection: Check for cgroup.kill support (Linux 5.14+)
	if _, err := os.Stat(filepath.Join(workerRoot, cgroupKillFile)); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "WARNING: %s not found. Atomic process tree termination is disabled (requires Linux 5.14+ or backports).\n", cgroupKillFile)
	}

	return nil
}

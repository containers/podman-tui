package specgen

import (
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"go.podman.io/common/pkg/config"
	"go.podman.io/storage/pkg/unshare"
)

func (s *SpecGenerator) InitResourceLimits(rtc *config.Config) {
	if s.ResourceLimits == nil || s.ResourceLimits.Pids == nil {
		if s.CgroupsMode != "disabled" {
			limit := rtc.PidsLimit()
			if limit != 0 {
				if s.ResourceLimits == nil {
					s.ResourceLimits = &spec.LinuxResources{}
				}
				s.ResourceLimits.Pids = &spec.LinuxPids{
					Limit: &limit,
				}
			}
		}
		return
	}

	// Rootless with cgroupfs the pids controller is often not delegated, and a
	// new cgroup has no pids limit anyway, so do not ask the runtime to set one.
	if limit := s.ResourceLimits.Pids.Limit; (limit == nil || *limit <= 0) &&
		unshare.IsRootless() && rtc.Engine.CgroupManager != config.SystemdCgroupsManager {
		s.ResourceLimits.Pids = nil
	}
}

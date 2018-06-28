package find

import (
	"regexp"
	"sort"
)

type installs []string

func (i installs) Prioritized() []string {
	ips := installationPriorities{}
	visited := map[string]bool{}
	for _, path := range i {
		if _, ok := visited[path]; ok {
			continue
		}
		visited[path] = true
		ips = append(ips, &installationPriority{
			path:   path,
			weight: calcInstallationPriority(path),
		})
	}

	sort.Sort(ips)

	return ips.Paths()
}

var priorities map[string]int

func calcInstallationPriority(path string) int {
	var cur *int
	for test, priority := range priorities {
		matched, err := regexp.MatchString(test, path)
		if err == nil && matched && (cur == nil || priority < *cur) {
			cur = &priority
		}
	}
	if cur == nil {
		return 10
	}
	return *cur
}

type installationPriority struct {
	path   string
	weight int
}

type installationPriorities []*installationPriority

func (ips installationPriorities) Len() int           { return len(ips) }
func (ips installationPriorities) Less(i, j int) bool { return ips[i].weight < ips[j].weight }
func (ips installationPriorities) Swap(i, j int)      { ips[i], ips[j] = ips[j], ips[i] }
func (ips installationPriorities) Paths() []string {
	paths := make([]string, len(ips))
	for i, ip := range ips {
		paths[i] = ip.path
	}
	return paths
}

package skynet

// 说明:
// 4个字段, 如果没有内容则不比较, 只有有内容时才比较.
//
// 也就是说一个空的Criteria和所有的ServiceInfo都匹配.
type Criteria struct {
	hosts            []string // ServiceInfo.Addr.IPAddress
	regions          []string
	instances        []string // Instance UUID
	ServiceCriterias []*ServiceCriteria
}

func NewCriteria() *Criteria {
	return &Criteria{
		hosts:            make([]string, 0, 0),
		regions:          make([]string, 0, 0),
		instances:        make([]string, 0, 0),
		ServiceCriterias: make([]*ServiceCriteria, 0, 0),
	}
}

func (criteria *Criteria) Matches(si *ServiceInfo) bool {
	// 1. 比较Hosts
	if len(criteria.hosts) > 0 && !exists(criteria.hosts, si.Addr.IPAddress) {
		return false
	}

	// 2. 比较Regions
	if len(criteria.regions) > 0 && !exists(criteria.regions, si.Region) {
		return false
	}

	// 3. 比较Instances
	if len(criteria.instances) > 0 && !exists(criteria.instances, si.InstanceUUID) {
		return false
	}

	// 4. 比较Services
	if len(criteria.ServiceCriterias) > 0 {
		match := false
		for _, sc := range criteria.ServiceCriterias {
			if sc.matches(si.Name, si.Version) {
				match = true
				break
			}
		}

		if !match {
			return false
		}
	}

	// 4个字段都比较过了, 返回true
	return true
}

func (criteria *Criteria) AddInstance(instanceUUID string) {
	if !exists(criteria.instances, instanceUUID) {
		criteria.instances = append(criteria.instances, instanceUUID)
	}
}

func (criteria *Criteria) AddHost(host string) {
	if !exists(criteria.hosts, host) {
		criteria.hosts = append(criteria.hosts, host)
	}
}

func (criteria *Criteria) AddRegion(region string) {
	if !exists(criteria.regions, region) {
		criteria.regions = append(criteria.regions, region)
	}
}

func (criteria *Criteria) AddService(name, version string) {
	for _, sc := range criteria.ServiceCriterias {
		if sc.matches(name, version) {
			return
		}
	}

	criteria.ServiceCriterias = append(criteria.ServiceCriterias, &ServiceCriteria{
		Name:    name,
		Version: version,
	})
}

type ServiceCriteria struct {
	Name    string
	Version string
}

func (sc *ServiceCriteria) matches(name, version string) bool {
	if sc.Name != "" && sc.Name != name {
		return false
	}

	if sc.Version != "" && sc.Version != version {
		return false
	}

	// name & version都相同, 返回true
	return true
}

func exists(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

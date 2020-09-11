package addrs

type ProviderReferences map[LocalProviderConfig]Provider

func (pr ProviderReferences) LocalNamesByAddr(addr Provider) []LocalProviderConfig {
	names := make([]LocalProviderConfig, 0)

	for lName, pAddr := range pr {
		if pAddr.Equals(addr) {
			names = append(names, lName)
		}
	}

	return names
}

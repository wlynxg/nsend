package stun

func Convert(m MappingBehavior, f FilteringBehavior) NatType {
	if m == EndpointIndependentMapping {
		switch f {
		case EndpointIndependentFiltering:
			return FullConeNAT
		case AddressDependentFiltering:
			return RestrictedConeNAT
		case AddressAndPortDependentFiltering:
			return PortRestrictedConeNAT
		}
	}
	return SymmetricNAT
}

package goapstra

import (
	"fmt"
)

//CONFIGLET_OS_SECTION_SUPPORT = {
//'cumulus': ('system', 'interface', 'file', 'frr', 'ospf'),
//'nxos': ('system', 'system_top', 'interface', 'ospf'),
//'eos': ('system', 'system_top', 'interface', 'ospf'),
//'junos': ('system', 'set_based_system', 'interface', 'set_based_interface',
//'delete_based_interface'),
//'sonic': ('system', 'file', 'frr', 'ospf'),
//}

type PlatformOS int
type platformOS string

const (
	PlatformOSCumulus = PlatformOS(iota)
	PlatformOSNxos
	PlatformOSEos
	PlatformOSJunos
	PlatformOSSonic
	PlatformOSUnknown = "unknown os '%s'"

	platformOSCumulus = platformOS("cumulus")
	platformOSNxos    = platformOS("nxos")
	platformOSEos     = platformOS("eos")
	platformOSJunos   = platformOS("junos")
	platformOSSonic   = platformOS("sonic")
	platformOSUnknown = "unknown type %d"
)

func (o PlatformOS) Int() int {
	return int(o)
}

func (o PlatformOS) String() string {
	switch o {
	case PlatformOSCumulus:
		return string(platformOSCumulus)
	case PlatformOSNxos:
		return string(platformOSNxos)
	case PlatformOSEos:
		return string(platformOSEos)
	case PlatformOSJunos:
		return string(platformOSJunos)
	case PlatformOSSonic:
		return string(platformOSSonic)
	default:
		return fmt.Sprintf(platformOSUnknown, o)
	}
}

func (o *PlatformOS) FromString(s string) error {
	i, err := platformOS(s).parse()
	if err != nil {
		return err
	}
	*o = PlatformOS(i)
	return nil
}

func (o PlatformOS) raw() platformOS {
	return platformOS(o.String())
}

func (o PlatformOS) ValidSections() []ConfigletSection {
	switch o {
	case PlatformOSCumulus:
		return []ConfigletSection{
			ConfigletSectionFile,
			ConfigletSectionFRR,
			ConfigletSectionInterface,
			ConfigletSectionOSPF,
			ConfigletSectionSystem,
		}
	case PlatformOSEos:
		return []ConfigletSection{
			ConfigletSectionInterface,
			ConfigletSectionOSPF,
			ConfigletSectionSystem,
			ConfigletSectionSystemTop,
		}
	case PlatformOSJunos:
		return []ConfigletSection{
			ConfigletSectionInterface,
			ConfigletSectionDeleteBasedInterface,
			ConfigletSectionSetBasedInterface,
			ConfigletSectionSystem,
			ConfigletSectionSetBasedSystem,
		}
	case PlatformOSNxos:
		return []ConfigletSection{
			ConfigletSectionSystem,
			ConfigletSectionInterface,
			ConfigletSectionSystemTop,
			ConfigletSectionOSPF,
		}
	case PlatformOSSonic:
		return []ConfigletSection{
			ConfigletSectionFile,
			ConfigletSectionFRR,
			ConfigletSectionOSPF,
			ConfigletSectionSystem,
		}
	}
	return nil
}

func (o platformOS) string() string {
	return string(o)
}

func (o platformOS) parse() (int, error) {
	switch o {
	case platformOSCumulus:
		return int(PlatformOSCumulus), nil
	case platformOSNxos:
		return int(PlatformOSNxos), nil
	case platformOSEos:
		return int(PlatformOSEos), nil
	case platformOSJunos:
		return int(PlatformOSJunos), nil
	case platformOSSonic:
		return int(PlatformOSSonic), nil
	default:
		return 0, fmt.Errorf(PlatformOSUnknown, o)
	}
}

// AllPlatformOS returns the []PlatformOS representing
// each supported PlatformOS
func AllPlatformOSes() []PlatformOS {
	i := 0
	var result []PlatformOS
	for {
		var sec PlatformOS
		err := sec.FromString(PlatformOS(i).String())
		if err != nil {
			return result[:i]
		}
		i++
	}
}

type ConfigletSection int
type configletSection string

const (
	ConfigletSectionSystem = ConfigletSection(iota)
	ConfigletSectionInterface
	ConfigletSectionFile
	ConfigletSectionFRR
	ConfigletSectionOSPF
	ConfigletSectionSystemTop
	ConfigletSectionSetBasedSystem
	ConfigletSectionSetBasedInterface
	ConfigletSectionDeleteBasedInterface
	ConfigletSectionUnknown = "unknown section '%s'"

	configletSectionSystem               = configletSection("system")
	configletSectionInterface            = configletSection("interface")
	configletSectionFile                 = configletSection("file")
	configletSectionFRR                  = configletSection("frr")
	configletSectionOSPF                 = configletSection("ospf")
	configletSectionSystemTop            = configletSection("system_top")
	configletSectionSetBasedSystem       = configletSection("set_based_system")
	configletSectionSetBasedInterface    = configletSection("set_based_interface")
	configletSectionDeleteBasedInterface = configletSection("delete_based_interface")
	configletSectionUnknown              = "unknown section %d"
)

func (o ConfigletSection) Int() int {
	return int(o)
}

func (o ConfigletSection) String() string {
	switch o {
	case ConfigletSectionSystem:
		return string(configletSectionSystem)
	case ConfigletSectionInterface:
		return string(configletSectionInterface)
	case ConfigletSectionFile:
		return string(configletSectionFile)
	case ConfigletSectionFRR:
		return string(configletSectionFRR)
	case ConfigletSectionOSPF:
		return string(configletSectionOSPF)
	case ConfigletSectionSystemTop:
		return string(configletSectionSystemTop)
	case ConfigletSectionSetBasedSystem:
		return string(configletSectionSetBasedSystem)
	case ConfigletSectionSetBasedInterface:
		return string(configletSectionSetBasedInterface)
	case ConfigletSectionDeleteBasedInterface:
		return string(configletSectionDeleteBasedInterface)
	default:
		return fmt.Sprintf(configletSectionUnknown, o)
	}
}

func (o ConfigletSection) raw() configletSection {
	return configletSection(o.String())
}

func (o configletSection) string() string {
	return string(o)
}

func (o configletSection) parse() (int, error) {
	switch o {
	case configletSectionSystem:
		return int(ConfigletSectionSystem), nil
	case configletSectionInterface:
		return int(ConfigletSectionInterface), nil
	case configletSectionFile:
		return int(ConfigletSectionFile), nil
	case configletSectionFRR:
		return int(ConfigletSectionFRR), nil
	case configletSectionOSPF:
		return int(ConfigletSectionOSPF), nil
	case configletSectionSystemTop:
		return int(ConfigletSectionSystemTop), nil
	case configletSectionSetBasedSystem:
		return int(ConfigletSectionSetBasedSystem), nil
	case configletSectionSetBasedInterface:
		return int(ConfigletSectionSetBasedInterface), nil
	case configletSectionDeleteBasedInterface:
		return int(ConfigletSectionDeleteBasedInterface), nil
	default:
		return 0, fmt.Errorf(ConfigletSectionUnknown, o)
	}
}

func (o *ConfigletSection) FromString(s string) error {
	i, err := configletSection(s).parse()
	if err != nil {
		return err
	}
	*o = ConfigletSection(i)
	return nil
}

// AllConfigletSections returns the []ConfigletSection representing
// each supported ConfigletSection
func AllConfigletSections() []ConfigletSection {
	i := 0
	var result []ConfigletSection
	for {
		var sec ConfigletSection
		err := sec.FromString(ConfigletSection(i).String())
		if err != nil {
			return result[:i]
		}
		i++
	}
}

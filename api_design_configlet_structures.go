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

type ApstraPlatformOS int
type apstraPlatformOS string

const (
	ApstraPlatformOSCumulus = ApstraPlatformOS(iota)
	ApstraPlatformOSNxos
	ApstraPlatformOSEos
	ApstraPlatformOSJunos
	ApstraPlatformOSSonic
	ApstraPlatformOSUnknown = "unknown os '%s'"
	apstraPlatformOSCumulus = apstraPlatformOS("cumulus")
	apstraPlatformOSNxos    = apstraPlatformOS("nxos")
	apstraPlatformOSEos     = apstraPlatformOS("eos")
	apstraPlatformOSJunos   = apstraPlatformOS("junos")
	apstraPlatformOSSonic   = apstraPlatformOS("sonic")
	apstraPlatformOSUnknown = "unknown type %d"
)

func (o ApstraPlatformOS) Int() int {
	return int(o)
}

func (o ApstraPlatformOS) String() string {
	switch o {
	case ApstraPlatformOSCumulus:
		return string(apstraPlatformOSCumulus)
	case ApstraPlatformOSNxos:
		return string(apstraPlatformOSNxos)
	case ApstraPlatformOSEos:
		return string(apstraPlatformOSEos)
	case ApstraPlatformOSJunos:
		return string(apstraPlatformOSJunos)
	case ApstraPlatformOSSonic:
		return string(apstraPlatformOSSonic)
	default:
		return fmt.Sprintf(apstraPlatformOSUnknown, o)
	}
}

func (o *ApstraPlatformOS) FromString(s string) error {
	i, err := apstraPlatformOS(s).parse()
	if err != nil {
		return err
	}
	*o = ApstraPlatformOS(i)
	return nil
}

func (o ApstraPlatformOS) raw() apstraPlatformOS {
	return apstraPlatformOS(o.String())
}

func (o ApstraPlatformOS) ValidSections() []ApstraConfigletSection {
	switch o {
	case ApstraPlatformOSCumulus:
		return []ApstraConfigletSection{
			ApstraConfigletSectionFile,
			ApstraConfigletSectionFRR,
			ApstraConfigletSectionInterface,
			ApstraConfigletSectionOSPF,
			ApstraConfigletSectionSystem,
		}
	case ApstraPlatformOSEos:
		return []ApstraConfigletSection{
			ApstraConfigletSectionInterface,
			ApstraConfigletSectionOSPF,
			ApstraConfigletSectionSystem,
			ApstraConfigletSectionSystemTop,
		}
	case ApstraPlatformOSJunos:
		return []ApstraConfigletSection{
			ApstraConfigletSectionInterface,
			ApstraConfigletSectionDeleteBasedInterface,
			ApstraConfigletSectionSetBasedInterface,
			ApstraConfigletSectionSystem,
			ApstraConfigletSectionSetBasedSystem,
		}
	case ApstraPlatformOSNxos:
		return []ApstraConfigletSection{
			ApstraConfigletSectionSystem,
			ApstraConfigletSectionInterface,
			ApstraConfigletSectionSystemTop,
			ApstraConfigletSectionOSPF,
		}
	case ApstraPlatformOSSonic:
		return []ApstraConfigletSection{
			ApstraConfigletSectionFile,
			ApstraConfigletSectionFRR,
			ApstraConfigletSectionOSPF,
			ApstraConfigletSectionSystem,
		}
	}
	return nil
}

func (o apstraPlatformOS) string() string {
	return string(o)
}

func (o apstraPlatformOS) parse() (int, error) {
	switch o {
	case apstraPlatformOSCumulus:
		return int(ApstraPlatformOSCumulus), nil
	case apstraPlatformOSNxos:
		return int(ApstraPlatformOSNxos), nil
	case apstraPlatformOSEos:
		return int(ApstraPlatformOSEos), nil
	case apstraPlatformOSJunos:
		return int(ApstraPlatformOSJunos), nil
	case apstraPlatformOSSonic:
		return int(ApstraPlatformOSSonic), nil
	default:
		return 0, fmt.Errorf(ApstraPlatformOSUnknown, o)
	}
}

// AllPlatformOS returns the []ApstraPlatformOS representing
// each supported ApstraPlatformOS
func AllPlatformOS() []ApstraPlatformOS {
	i := 0
	var result []ApstraPlatformOS
	for {
		var sec ApstraPlatformOS
		err := sec.FromString(ApstraPlatformOS(i).String())
		if err != nil {
			return result[:i]
		}
		i++
	}
}

type ApstraConfigletSection int
type apstraConfigletSection string

const (
	ApstraConfigletSectionSystem = ApstraConfigletSection(iota)
	ApstraConfigletSectionInterface
	ApstraConfigletSectionFile
	ApstraConfigletSectionFRR
	ApstraConfigletSectionOSPF
	ApstraConfigletSectionSystemTop
	ApstraConfigletSectionSetBasedSystem
	ApstraConfigletSectionSetBasedInterface
	ApstraConfigletSectionDeleteBasedInterface
	ApstraConfigletSectionUnknown = "unknown section '%s'"

	apstraConfigletSectionSystem               = apstraConfigletSection("system")
	apstraConfigletSectionInterface            = apstraConfigletSection("interface")
	apstraConfigletSectionFile                 = apstraConfigletSection("file")
	apstraConfigletSectionFRR                  = apstraConfigletSection("frr")
	apstraConfigletSectionOSPF                 = apstraConfigletSection("ospf")
	apstraConfigletSectionSystemTop            = apstraConfigletSection("system_top")
	apstraConfigletSectionSetBasedSystem       = apstraConfigletSection("set_based_system")
	apstraConfigletSectionSetBasedInterface    = apstraConfigletSection("set_based_interface")
	apstraConfigletSectionDeleteBasedInterface = apstraConfigletSection("delete_based_interface")
	apstraConfigletSectionUnknown              = "unknown section %d"
)

func (o ApstraConfigletSection) Int() int {
	return int(o)
}

func (o ApstraConfigletSection) String() string {
	switch o {
	case ApstraConfigletSectionSystem:
		return string(apstraConfigletSectionSystem)
	case ApstraConfigletSectionInterface:
		return string(apstraConfigletSectionInterface)
	case ApstraConfigletSectionFile:
		return string(apstraConfigletSectionFile)
	case ApstraConfigletSectionFRR:
		return string(apstraConfigletSectionFRR)
	case ApstraConfigletSectionOSPF:
		return string(apstraConfigletSectionOSPF)
	case ApstraConfigletSectionSystemTop:
		return string(apstraConfigletSectionSystemTop)
	case ApstraConfigletSectionSetBasedSystem:
		return string(apstraConfigletSectionSetBasedSystem)
	case ApstraConfigletSectionSetBasedInterface:
		return string(apstraConfigletSectionSetBasedInterface)
	case ApstraConfigletSectionDeleteBasedInterface:
		return string(apstraConfigletSectionDeleteBasedInterface)
	default:
		return fmt.Sprintf(apstraConfigletSectionUnknown, o)
	}
}

func (o ApstraConfigletSection) raw() apstraConfigletSection {
	return apstraConfigletSection(o.String())
}

func (o apstraConfigletSection) string() string {
	return string(o)
}

func (o apstraConfigletSection) parse() (int, error) {
	switch o {
	case apstraConfigletSectionSystem:
		return int(ApstraConfigletSectionSystem), nil
	case apstraConfigletSectionInterface:
		return int(ApstraConfigletSectionInterface), nil
	case apstraConfigletSectionFile:
		return int(ApstraConfigletSectionFile), nil
	case apstraConfigletSectionFRR:
		return int(ApstraConfigletSectionFRR), nil
	case apstraConfigletSectionOSPF:
		return int(ApstraConfigletSectionOSPF), nil
	case apstraConfigletSectionSystemTop:
		return int(ApstraConfigletSectionSystemTop), nil
	case apstraConfigletSectionSetBasedSystem:
		return int(ApstraConfigletSectionSetBasedSystem), nil
	case apstraConfigletSectionSetBasedInterface:
		return int(ApstraConfigletSectionSetBasedInterface), nil
	case apstraConfigletSectionDeleteBasedInterface:
		return int(ApstraConfigletSectionDeleteBasedInterface), nil
	default:
		return 0, fmt.Errorf(ApstraConfigletSectionUnknown, o)
	}
}

func (o *ApstraConfigletSection) FromString(s string) error {
	i, err := apstraConfigletSection(s).parse()
	if err != nil {
		return err
	}
	*o = ApstraConfigletSection(i)
	return nil
}

// AllConfigletSections returns the []ApstraConfigletSection representing
// each supported ApstraConfigletSection
func AllConfigletSections() []ApstraConfigletSection {
	i := 0
	var result []ApstraConfigletSection
	for {
		var sec ApstraConfigletSection
		err := sec.FromString(ApstraConfigletSection(i).String())
		if err != nil {
			return result[:i]
		}
		i++
	}
}

// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

func ValidConfigletSections(platform enum.ConfigletStyle) []enum.ConfigletSection {
	switch platform {
	case enum.ConfigletStyleCumulus:
		return []enum.ConfigletSection{
			enum.ConfigletSectionFile,
			enum.ConfigletSectionFrr,
			enum.ConfigletSectionInterface,
			enum.ConfigletSectionOspf,
			enum.ConfigletSectionSystem,
		}
	case enum.ConfigletStyleEos:
		return []enum.ConfigletSection{
			enum.ConfigletSectionInterface,
			enum.ConfigletSectionOspf,
			enum.ConfigletSectionSystem,
			enum.ConfigletSectionSystemTop,
		}
	case enum.ConfigletStyleJunos:
		return []enum.ConfigletSection{
			enum.ConfigletSectionInterface,
			enum.ConfigletSectionDeleteBasedInterface,
			enum.ConfigletSectionSetBasedInterface,
			enum.ConfigletSectionSystem,
			enum.ConfigletSectionSetBasedSystem,
		}
	case enum.ConfigletStyleNxos:
		return []enum.ConfigletSection{
			enum.ConfigletSectionSystem,
			enum.ConfigletSectionInterface,
			enum.ConfigletSectionSystemTop,
			enum.ConfigletSectionOspf,
		}
	case enum.ConfigletStyleSonic:
		return []enum.ConfigletSection{
			enum.ConfigletSectionFile,
			enum.ConfigletSectionFrr,
			enum.ConfigletSectionOspf,
			enum.ConfigletSectionSystem,
		}
	}
	return nil
}

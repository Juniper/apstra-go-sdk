package enum

const (
	Flavor = EnumType(iota)
	Size
	Sauce
)

var enumFuncs = []func() Value{
	FlavorChocolate,
	FlavorStrawberry,
	FlavorVanilla,

	SizeSmall,
	SizeMedium,
	SizeLarge,

	SauceChocolate,
	SauceCaramel,
}

func FlavorChocolate() Value  { return newInstance(Flavor, "chocolate") }
func FlavorStrawberry() Value { return newInstance(Flavor, "strawberry") }
func FlavorVanilla() Value    { return newInstance(Flavor, "vanilla") }

func SizeSmall() Value  { return newInstance(Size, "small") }
func SizeMedium() Value { return newInstance(Size, "medium") }
func SizeLarge() Value  { return newInstance(Size, "large") }

func SauceChocolate() Value { return newInstance(Sauce, "chocolate") }
func SauceCaramel() Value   { return newInstance(Sauce, "caramel") }

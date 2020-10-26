package strs

func ParseSpec(spec ...string) Spec {
	var p Spec
	for _, s := range spec {
		switch s {
		case "lower", "l":
			p |= Lower
		case "upper", "u":
			p |= Upper
		case "num", "n":
			p |= Num
		case "symbol", "s":
			p |= Symbol
		default:
			for i := range s {
				switch s[i] {
				case 'l':
					p |= Lower
				case 'u':
					p |= Upper
				case 'n':
					p |= Num
				case 's':
					p |= Symbol
				}
			}
		}
	}

	if p == 0 {
		p = NoSymbol
	}

	return p
}

package chrome

func analyzeFac(facilities []string) *Facility {
	fac := &Facility{}
	for _, facility := range facilities {
		switch facility {
		case "WADING POOL", "50M LAP POOL", "JACUZZI POOL", "50M FREEFORM POOL", "BEACH SPLASH POOL", "FAMILY POOL", "REFLECTION POOL":
			fac.pool = true
		case "TENNIS COURT":
			fac.tennisCourt = true
		case "READING Corner":
			fac.readingCorner = true
		case "FITNESS STATION", "FITNESS ALCOVE":
			fac.fitnessArea = true
		case "INDOOR GYM", "HYDRO GYM STATION":
			fac.gymnasium = true
		case "BBQ AREA":
			fac.bbqPit = true
		case "24-HOUR SECURITY":
			fac.security = true
		}
	}

	return fac
}

package service

func MakeKey(frefix, idaddress string) string {
	if idaddress != "" {
		idaddress = "/" + idaddress
	}

	return frefix + "/microadmin" + idaddress
}
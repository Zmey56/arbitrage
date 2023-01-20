// Main file to save ifo from huobi to local DB
package getinfohuobi

func GetInfoFromHuobi(fiat string) {
	//get peymont from huobi and save them
	SavePaymentToJSONHuobi(GetPeymontMethodsHuobi(fiat), fiat)
	//Get coin ID and save them
	GetInfoCryptoHuobi(fiat)
}

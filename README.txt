<Add project name>

<Add project introduction>




<Add our name&id here>

<Add description of files, functions and etc.>

DAO:

func ReadData() is for reading the data from xlsx. 

func SetStockData(), SetBMData() and SetVMQScore() are for setting the data from xlsx to struct.

func CalStock() and CalBench() are for calculation.

func BuildGICSList(), BuildRegionList(), BuildRegionMap() and CheckRegion() are for creating&checking the list of GICS&Region code.

func SetValue(), SetBValue(), SetPresentage() and SetBPresentage() are pointer function, for setting the value&presentage in struct.

func StringToFloat() is for converting the string(all data from xlsx are string) to float64, which can be use to calculate. 
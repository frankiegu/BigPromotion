package util



const (
	USER_DETAILS = "USER_DETAIL"
	COURSE_DETAILS = "COURSE_DETAILS"
	GLOBAL_DONATE_DETAILS = "GLOBAL_DONATE_DETAILS"
	DAY_DONATE_COURSE_NUMS = "DAY_DONATE_DETAILS"
	DAY_DONATE_PERSON_NUMS = "DAY_DONATE_PERSON_NUMS"
	ACTIVITY_CENTER = "ACTIVITY_CENTER"
	PUBLIC_DAY_DISTRIBUTE_LOCK = "PUBLIC_DAY_DISTRIBUTE_LOCK"

	PUBLIC_DAY_CODIS_PREFIX = "PUBLIC_DAY_"
	TEST_PUBLIC_DAY_CODIS_PREFIX = "TEST_PUBLIC_DAY_"

	CUR_LOVE_NUMS = "CUR_LOVE_NUMS"
	DONATE_COURSE_NUMS = "DONATE_COURSE_NUMS"


)

type PublicDayCodisUtil struct {

}

var publicDayCodisUtil *PublicDayCodisUtil

func GetCodisUtilInstance() *PublicDayCodisUtil {
	if publicDayCodisUtil == nil {
		publicDayCodisUtil = &PublicDayCodisUtil{}
	}
	return publicDayCodisUtil
}

func (publicDayCodisUtil *PublicDayCodisUtil) GetCodisKey(key string) string {
	var codisKey string
	kind := getCodisTestMode()
	if kind {
		codisKey = TEST_PUBLIC_DAY_CODIS_PREFIX + key
	} else {
		codisKey = PUBLIC_DAY_CODIS_PREFIX + key
	}

	return codisKey
}



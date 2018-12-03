package bean



type UserDetails struct {

	AllloveNumRecordList []int `json:"alllove_num_record_list"`
	DonateCourseTimeRecords []int `json:"donate_course_time_records"`
	CidSet *Set `json:"cid_set"`
	PreviousLearnTime int64 `json:"previous_learn_time"`
	ExchangeLoveNumRecords []int `json:"exchange_love_num_records"`
	LatestShareTime string `json:"latest_share_time"`
	LatestDonateTime string `json:"latest_donate_time"`
	DaySignUpNUms map[string]int `json:"day_sign_up_nums"`
	TimerTime []string `json:"timer_time"`

}

func (u *UserDetails) NewUserDetails() *UserDetails {
	return &UserDetails{
		AllloveNumRecordList: make([]int, 1),
		DonateCourseTimeRecords: make([]int, 1),
		CidSet:NewSet(),
		DaySignUpNUms: map[string]int{},
		TimerTime: make([]string, 1),
	}
}











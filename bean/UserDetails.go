package bean



type UserDetails struct {

	AllloveNumRecordList []int `json:"alllove_num_record_list"`
	DonateCourseTimeRecords []int `json:"donate_course_time_records"`
	cidSet *Set `json:"cid_set"`
	previousLearnTime float64 `json:"previous_learn_time"`
	exchangeLoveNumRecords []int `json:"exchange_love_num_records"`
	latestShareTime string `json:"latest_share_time"`
	latestDonateTime string `json:"latest_donate_time"`
	daySignUpNUms map[string]int `json:"day_sign_up_n_ums"`
	timerTime []string `json:"timer_time"`

}

func (u *UserDetails) NewUserDetails() *UserDetails {
	return &UserDetails{
		daySignUpNUms: map[string]int{},
		cidSet:NewSet(),
	}
}











package jobs

import "github.com/brunodmartins/church-members-api/platform/cdi"

func BuildJob(jobType JobType) Job {
	memberService := cdi.ProvideMemberService()
	userService := cdi.ProvideUserService()
	churchService := cdi.ProvideChurchService()
	switch jobType {
	case WEEKLY_BIRTHDAYS:
		job := newWeeklyBirthDaysJob(memberService, cdi.ProvideEmailService(), userService)
		return newChurchWrapperJob(churchService, job)
	case DAILY_BIRTHDAYS:
		job := newDailyBirthDaysJob(memberService, cdi.ProvideNotificationService(), userService)
		return newChurchWrapperJob(churchService, job)
	}
	return nil
}

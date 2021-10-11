package jobs

import "github.com/BrunoDM2943/church-members-api/platform/cdi"

func BuildJob(jobType JobType) Job {
	memberService := cdi.ProvideMemberService()
	userService := cdi.ProvideUserService()
	switch jobType {
	case WEEKLY_BIRTHDAYS:
		return newWeeklyBirthDaysJob(memberService, cdi.ProvideEmailService(), userService)
	case DAILY_BIRTHDAYS:
		return newDailyBirthDaysJob(memberService, cdi.ProvideNotificationService(), userService)
	}
	return nil
}

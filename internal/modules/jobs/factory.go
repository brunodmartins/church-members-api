package jobs

import "github.com/BrunoDM2943/church-members-api/platform/cdi"

func BuildJob(jobType JobType) Job {
	memberService := cdi.ProvideMemberService()
	notificationService := cdi.ProvideNotificationService()
	switch jobType {
	case WEEKLY_BIRTHDAYS:
		return newWeeklyBirthDaysJob(memberService, notificationService)
	case DAILY_BIRTHDAYS:
		return newDailyBirthDaysJob(memberService, notificationService)
	}
	return nil
}

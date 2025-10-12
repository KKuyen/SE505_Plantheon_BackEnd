package activities

import (
	"time"
)

// ActivityResponse represents activity response
type ActivityResponse struct {
	ID              string     `json:"id"`
	Description     *string    `json:"description"`
	Description2    *string    `json:"description2"`
	Description3    *string    `json:"description3"`
    TimeStart       *time.Time `json:"time_start"`
    TimeEnd         *time.Time `json:"time_end"`
    Day             *bool      `json:"day"`
	Money           *float64   `json:"money"`
    Type            string     `json:"type"`
	Title           string     `json:"title"`
	IsRepeat        *string    `json:"is_repeat"`
    Repeat          *string    `json:"repeat"`
	EndRepeatDay    *time.Time `json:"end_repeat_day"`
	AlertTime       *string    `json:"alert_time"`
	Object          *string    `json:"object"`
	Amount          *int       `json:"amount"`
	Unit            *string    `json:"unit"`
	Purpose         *string    `json:"purpose"`
	TargetPerson    *string    `json:"target_person"`
	SourcePerson    *string    `json:"source_person"`
	AttachedLink    *string    `json:"attached_link"`
	Note            *string    `json:"note"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Minimal activity item for calendar list (only title)
type ActivityCalendarItem struct {
    Title string `json:"title"`
}

// Activity item for day view (minimal info)
type ActivityDayItem struct {
    ID        string     `json:"id"`
    Title     string     `json:"title"`
    Type      string     `json:"type"`
    TimeStart *time.Time `json:"time_start"`
    TimeEnd   *time.Time `json:"time_end"`
    Day       *bool      `json:"day"`
}

// Calendar day item
type ActivityCalendarDay struct {
    Date       string                 `json:"date"`       // YYYY-MM-DD
    Activities []ActivityCalendarItem `json:"activities"`
}

// CreateActivityRequest represents activity creation request
type CreateActivityRequest struct {
	Description     *string    `json:"description"`
	Description2    *string    `json:"description2"`
	Description3    *string    `json:"description3"`
    TimeStart       *time.Time `json:"time_start"`
    TimeEnd         *time.Time `json:"time_end"`
    Day             *bool      `json:"day"`
	Money           *float64   `json:"money"`
    Type            string     `json:"type" binding:"required"`
	Title           string     `json:"title" binding:"required"`
	IsRepeat        *string    `json:"is_repeat"`
    Repeat          *string    `json:"repeat"`
	EndRepeatDay    *time.Time `json:"end_repeat_day"`
	AlertTime       *string    `json:"alert_time"`
	Object          *string    `json:"object"`
	Amount          *int       `json:"amount"`
	Unit            *string    `json:"unit"`
	Purpose         *string    `json:"purpose"`
	TargetPerson    *string    `json:"target_person"`
	SourcePerson    *string    `json:"source_person"`
	AttachedLink    *string    `json:"attached_link"`
	Note            *string    `json:"note"`
}

// UpdateActivityRequest represents activity update request
type UpdateActivityRequest struct {
	Description     *string    `json:"description"`
	Description2    *string    `json:"description2"`
	Description3    *string    `json:"description3"`
    TimeStart       *time.Time `json:"time_start"`
    TimeEnd         *time.Time `json:"time_end"`
    Day             *bool      `json:"day"`
	Money           *float64   `json:"money"`
	Type            *string    `json:"type"`
	Title           *string    `json:"title"`
	IsRepeat        *string    `json:"is_repeat"`
    Repeat          *string    `json:"repeat"`
	EndRepeatDay    *time.Time `json:"end_repeat_day"`
	AlertTime       *string    `json:"alert_time"`
	Object          *string    `json:"object"`
	Amount          *int       `json:"amount"`
	Unit            *string    `json:"unit"`
	Purpose         *string    `json:"purpose"`
	TargetPerson    *string    `json:"target_person"`
	SourcePerson    *string    `json:"source_person"`
	AttachedLink    *string    `json:"attached_link"`
	Note            *string    `json:"note"`
}

// ActivitiesListResponse represents paginated activities list response
type ActivitiesListResponse struct {
	Activities []ActivityResponse `json:"activities"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// ToActivityResponse converts Activity to ActivityResponse
func (a *Activity) ToActivityResponse() ActivityResponse {
	return ActivityResponse{
		ID:              a.ID,
		Description:     a.Description,
		Description2:    a.Description2,
		Description3:    a.Description3,
		TimeStart:       a.TimeStart,
		TimeEnd:         a.TimeEnd,
		Day:             a.Day,
		Money:           a.Money,
		Type:            a.Type,
		Title:           a.Title,
        IsRepeat:        a.IsRepeat,
        Repeat:          a.Repeat,
		EndRepeatDay:    a.EndRepeatDay,
		AlertTime:       a.AlertTime,
		Object:          a.Object,
		Amount:          a.Amount,
		Unit:            a.Unit,
		Purpose:         a.Purpose,
		TargetPerson:    a.TargetPerson,
		SourcePerson:    a.SourcePerson,
		AttachedLink:    a.AttachedLink,
		Note:            a.Note,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

// ToActivitiesListResponse converts activities list to paginated response
func ToActivitiesListResponse(activities []Activity, total int64, page, limit int) ActivitiesListResponse {
	var response []ActivityResponse
	for _, activity := range activities {
		response = append(response, activity.ToActivityResponse())
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return ActivitiesListResponse{
		Activities: response,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

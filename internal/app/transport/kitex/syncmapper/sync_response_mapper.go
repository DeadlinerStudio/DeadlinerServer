package syncmapper

import (
	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1"
)

func ToPullChangesResponse(result *domainSync.PullChangesResult) *v1.PullChangesResponse {
	if result == nil {
		return &v1.PullChangesResponse{}
	}

	return &v1.PullChangesResponse{
		DeadlineChanges: toKitexDeadlineChanges(result.DeadlineChanges),
		HabitChanges:    toKitexHabitChanges(result.HabitChanges),
		NextCursor:      result.NextCursor,
		HasMore:         result.HasMore,
	}
}

func ToPushChangesResponse(result *domainSync.PushChangesResult) *v1.PushChangesResponse {
	if result == nil {
		return &v1.PushChangesResponse{}
	}

	return &v1.PushChangesResponse{
		Results:         toKitexMutationResults(result.Results),
		DeadlineChanges: toKitexDeadlineChanges(result.DeadlineChanges),
		HabitChanges:    toKitexHabitChanges(result.HabitChanges),
		NextCursor:      result.NextCursor,
	}
}

func toKitexMutationResults(results []domainSync.MutationResult) []*v1.MutationResult_ {
	mapped := make([]*v1.MutationResult_, 0, len(results))
	for _, result := range results {
		mapped = append(mapped, &v1.MutationResult_{
			MutationId:      result.MutationID,
			EntityUid:       result.EntityUID,
			Accepted:        result.Accepted,
			RejectionReason: result.RejectionReason,
			ServerVersion:   toKitexServerVersion(result.ServerVersion),
			Replayed:        result.Replayed,
			Status:          stringPtrOrNil(result.Status),
		})
	}
	return mapped
}

func toKitexDeadlineChanges(changes []domainSync.DeadlineChange) []*v1.DeadlineChange {
	mapped := make([]*v1.DeadlineChange, 0, len(changes))
	for _, change := range changes {
		mapped = append(mapped, &v1.DeadlineChange{
			EntityUid:     change.EntityUID,
			Deleted:       change.Deleted,
			ServerVersion: toKitexServerVersion(change.ServerVersion),
			Doc:           toKitexDeadlineDocument(change.Document),
		})
	}
	return mapped
}

func toKitexHabitChanges(changes []domainSync.HabitChange) []*v1.HabitChange {
	mapped := make([]*v1.HabitChange, 0, len(changes))
	for _, change := range changes {
		mapped = append(mapped, &v1.HabitChange{
			EntityUid:     change.EntityUID,
			Deleted:       change.Deleted,
			ServerVersion: toKitexServerVersion(change.ServerVersion),
			Doc:           toKitexHabitDocument(change.Document),
		})
	}
	return mapped
}

func toKitexServerVersion(version domainSync.ServerVersion) *v1.ServerVersion {
	if version.ChangeID == 0 && version.CommittedAt == "" {
		return nil
	}

	return &v1.ServerVersion{
		ChangeId:    version.ChangeID,
		CommittedAt: version.CommittedAt,
	}
}

func toKitexDeadlineDocument(doc domainSync.DeadlineDocument) *v1.DeadlineDocument {
	subTasks := make([]*v1.SubTask, 0, len(doc.SubTasks))
	for _, subTask := range doc.SubTasks {
		subTasks = append(subTasks, &v1.SubTask{
			Id:          subTask.ID,
			Content:     subTask.Content,
			IsCompleted: subTask.IsCompleted,
			SortOrder:   subTask.SortOrder,
			CreatedAt:   subTask.CreatedAt,
			UpdatedAt:   subTask.UpdatedAt,
		})
	}

	return &v1.DeadlineDocument{
		Uid:             doc.UID,
		LegacyId:        doc.LegacyID,
		Name:            doc.Name,
		StartTime:       doc.StartTime,
		EndTime:         doc.EndTime,
		State:           string(doc.State),
		CompleteTime:    doc.CompleteTime,
		Note:            doc.Note,
		IsStared:        doc.IsStared,
		Type:            string(doc.Type),
		HabitCount:      doc.HabitCount,
		HabitTotalCount: doc.HabitTotalCount,
		CalendarEvent:   doc.CalendarEvent,
		Timestamp:       doc.Timestamp,
		SubTasks:        subTasks,
	}
}

func toKitexHabitDocument(doc domainSync.HabitDocument) *v1.HabitDocument {
	records := make([]*v1.HabitRecord, 0, len(doc.Records))
	for _, record := range doc.Records {
		records = append(records, &v1.HabitRecord{
			Date:      record.Date,
			Count:     record.Count,
			Status:    string(record.Status),
			CreatedAt: record.CreatedAt,
		})
	}

	return &v1.HabitDocument{
		DdlUid: doc.DDLUID,
		Habit: &v1.HabitConfig{
			Name:           doc.Habit.Name,
			Description:    doc.Habit.Description,
			Color:          doc.Habit.Color,
			IconKey:        doc.Habit.IconKey,
			Period:         string(doc.Habit.Period),
			TimesPerPeriod: doc.Habit.TimesPerPeriod,
			GoalType:       string(doc.Habit.GoalType),
			TotalTarget:    doc.Habit.TotalTarget,
			CreatedAt:      doc.Habit.CreatedAt,
			UpdatedAt:      doc.Habit.UpdatedAt,
			Status:         string(doc.Habit.Status),
			SortOrder:      doc.Habit.SortOrder,
			AlarmTime:      doc.Habit.AlarmTime,
		},
		Records: records,
	}
}

func stringPtrOrNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

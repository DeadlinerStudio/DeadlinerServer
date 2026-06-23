package syncmapper

import (
	"fmt"

	appsync "github.com/aritxonly/deadlinerserver/internal/app/service/sync"
	domainSync "github.com/aritxonly/deadlinerserver/internal/domain/sync"
	v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1"
)

func ToPullChangesInput(req *v1.PullChangesRequest) appsync.PullChangesInput {
	if req == nil {
		return appsync.PullChangesInput{}
	}

	return appsync.PullChangesInput{
		DeviceUID:      req.DeviceUid,
		Cursor:         req.Cursor,
		Limit:          req.Limit,
		IncludeDeleted: req.IncludeDeleted,
	}
}

func ToPushChangesInput(req *v1.PushChangesRequest) (appsync.PushChangesInput, error) {
	if req == nil {
		return appsync.PushChangesInput{}, nil
	}

	mutations := make([]domainSync.Mutation, 0, len(req.Mutations))
	for _, mutation := range req.Mutations {
		mapped, err := toDomainMutation(mutation)
		if err != nil {
			return appsync.PushChangesInput{}, err
		}
		mutations = append(mutations, mapped)
	}

	return appsync.PushChangesInput{
		DeviceUID:  req.DeviceUid,
		BaseCursor: req.BaseCursor,
		Mutations:  mutations,
	}, nil
}

func toDomainMutation(mutation *v1.Mutation) (domainSync.Mutation, error) {
	if mutation == nil {
		return domainSync.Mutation{}, nil
	}

	payloadCount := 0
	if mutation.Payload != nil {
		payloadCount = mutation.Payload.CountSetFieldsMutationPayload()
	}
	if payloadCount > 1 {
		return domainSync.Mutation{}, fmt.Errorf("mutation %s contains multiple payload variants", mutation.MutationId)
	}

	result := domainSync.Mutation{
		MutationID: mutation.MutationId,
		DeviceUID:  mutation.DeviceUid,
		EntityUID:  mutation.EntityUid,
	}
	if mutation.IsSetClientVersion() {
		result.ClientVersion = domainSync.LogicalVersion{
			TS:  mutation.ClientVersion.GetTs(),
			Ctr: mutation.ClientVersion.GetCtr(),
			Dev: mutation.ClientVersion.GetDev(),
		}
	}
	if mutation.IsSetBaseChangeId() {
		result.BaseChangeID = mutation.GetBaseChangeId()
	}
	if mutation.Payload == nil {
		return result, nil
	}
	if mutation.Payload.IsSetDeadline() {
		result.Deadline = toDomainDeadlinePatch(mutation.Payload.GetDeadline())
	}
	if mutation.Payload.IsSetHabit() {
		result.Habit = toDomainHabitPatch(mutation.Payload.GetHabit())
	}

	return result, nil
}

func toDomainDeadlinePatch(mutation *v1.DeadlineMutation) *domainSync.DeadlinePatch {
	if mutation == nil {
		return nil
	}

	return &domainSync.DeadlinePatch{
		Deleted:  mutation.Deleted,
		Document: toDomainDeadlineDocument(mutation.Doc),
	}
}

func toDomainHabitPatch(mutation *v1.HabitMutation) *domainSync.HabitPatch {
	if mutation == nil {
		return nil
	}

	return &domainSync.HabitPatch{
		Deleted:  mutation.Deleted,
		Document: toDomainHabitDocument(mutation.Doc),
	}
}

func toDomainDeadlineDocument(doc *v1.DeadlineDocument) domainSync.DeadlineDocument {
	if doc == nil {
		return domainSync.DeadlineDocument{}
	}

	subTasks := make([]domainSync.SubTask, 0, len(doc.SubTasks))
	for _, subTask := range doc.SubTasks {
		subTasks = append(subTasks, domainSync.SubTask{
			ID:          subTask.GetId(),
			Content:     subTask.GetContent(),
			IsCompleted: subTask.GetIsCompleted(),
			SortOrder:   subTask.GetSortOrder(),
			CreatedAt:   subTask.GetCreatedAt(),
			UpdatedAt:   subTask.GetUpdatedAt(),
		})
	}

	return domainSync.DeadlineDocument{
		UID:             doc.Uid,
		LegacyID:        doc.LegacyId,
		Name:            doc.Name,
		StartTime:       doc.StartTime,
		EndTime:         doc.EndTime,
		State:           domainSync.DeadlineState(doc.State),
		CompleteTime:    doc.CompleteTime,
		Note:            doc.Note,
		IsStared:        doc.IsStared,
		Type:            domainSync.DeadlineType(doc.Type),
		HabitCount:      doc.HabitCount,
		HabitTotalCount: doc.HabitTotalCount,
		CalendarEvent:   doc.CalendarEvent,
		Timestamp:       doc.Timestamp,
		SubTasks:        subTasks,
	}
}

func toDomainHabitDocument(doc *v1.HabitDocument) domainSync.HabitDocument {
	if doc == nil {
		return domainSync.HabitDocument{}
	}

	records := make([]domainSync.HabitRecord, 0, len(doc.Records))
	for _, record := range doc.Records {
		records = append(records, domainSync.HabitRecord{
			Date:      record.GetDate(),
			Count:     record.GetCount(),
			Status:    domainSync.HabitRecordStatus(record.GetStatus()),
			CreatedAt: record.GetCreatedAt(),
		})
	}

	habitConfig := domainSync.HabitConfig{}
	if doc.IsSetHabit() {
		habitConfig = domainSync.HabitConfig{
			Name:           doc.Habit.GetName(),
			Description:    doc.Habit.GetDescription(),
			IconKey:        doc.Habit.GetIconKey(),
			Period:         domainSync.HabitPeriod(doc.Habit.GetPeriod()),
			TimesPerPeriod: doc.Habit.GetTimesPerPeriod(),
			GoalType:       domainSync.HabitGoalType(doc.Habit.GetGoalType()),
			TotalTarget:    doc.Habit.GetTotalTarget(),
			CreatedAt:      doc.Habit.GetCreatedAt(),
			UpdatedAt:      doc.Habit.GetUpdatedAt(),
			Status:         domainSync.HabitStatus(doc.Habit.GetStatus()),
			SortOrder:      doc.Habit.GetSortOrder(),
			AlarmTime:      doc.Habit.GetAlarmTime(),
		}
	}

	return domainSync.HabitDocument{
		DDLUID:  doc.DdlUid,
		Habit:   habitConfig,
		Records: records,
	}
}

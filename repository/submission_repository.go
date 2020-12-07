package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// SubmissionRepository ..
type SubmissionRepository interface {
	Create(ctx context.Context, submission *model.Submission) error
	FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) ([]*model.Submission, int64, error)
}

type submissionRepo struct {
	db *gorm.DB
}

// NewSubmissionRepo ..
func NewSubmissionRepo(db *gorm.DB) SubmissionRepository {
	return &submissionRepo{
		db: db,
	}
}

func (s *submissionRepo) Create(ctx context.Context, submission *model.Submission) error {
	err := s.db.Create(submission).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"submission": utils.Dump(submission),
		}).Error(err)
	}

	return err
}

func (s *submissionRepo) FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) ([]*model.Submission, int64, error) {
	count := int64(0)
	err := s.db.Model(model.Submission{}).Where("assignment_id = ?", assignmentID).Count(&count).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"assignmentID": assignmentID,
		}).Error(err)
		return nil, count, err
	}

	submissions := []*model.Submission{}
	err = s.db.Where("assignment_id = ?", assignmentID).Limit(int(cursor.GetSize())).
		Offset(int(cursor.GetOffset())).Order("created_at " + cursor.GetSort()).Find(&submissions).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"cursor":       cursor,
			"assignmentID": assignmentID,
		}).Error(err)
		return nil, count, err
	}

	return submissions, count, nil
}
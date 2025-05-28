package usecase_test

import (
	"context"
	"testing"

	course "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	user "skillForce/internal/models/user"
	"skillForce/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMockCourseRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecase.NewMockCourseRepository(ctrl)
	ctx := context.Background()

	// Существующие тесты...

	t.Run("Test CreatePart", func(t *testing.T) {
		expectedPartID := 1
		mockRepo.EXPECT().CreatePart(ctx, &course.CoursePart{}, 1).Return(expectedPartID, nil)

		partID, err := mockRepo.CreatePart(ctx, &course.CoursePart{}, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedPartID, partID)
	})

	t.Run("Test CreateTextLesson", func(t *testing.T) {
		mockRepo.EXPECT().CreateTextLesson(ctx, &course.LessonPoint{}, 1).Return(nil)

		err := mockRepo.CreateTextLesson(ctx, &course.LessonPoint{}, 1)
		assert.NoError(t, err)
	})

	t.Run("Test CreateVideoLesson", func(t *testing.T) {
		mockRepo.EXPECT().CreateVideoLesson(ctx, &course.LessonPoint{}, 1).Return(nil)

		err := mockRepo.CreateVideoLesson(ctx, &course.LessonPoint{}, 1)
		assert.NoError(t, err)
	})

	t.Run("Test DeleteCourseFromFavourites", func(t *testing.T) {
		mockRepo.EXPECT().DeleteCourseFromFavourites(ctx, 1, 1).Return(nil)

		err := mockRepo.DeleteCourseFromFavourites(ctx, 1, 1)
		assert.NoError(t, err)
	})

	t.Run("Test GetBucketByLessonId", func(t *testing.T) {
		expectedBucket := &course.LessonBucket{}
		mockRepo.EXPECT().GetBucketByLessonId(ctx, 1).Return(expectedBucket, nil)

		bucket, err := mockRepo.GetBucketByLessonId(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedBucket, bucket)
	})

	t.Run("Test GetBucketCourses", func(t *testing.T) {
		expectedCourses := []*course.Course{}
		mockRepo.EXPECT().GetBucketCourses(ctx).Return(expectedCourses, nil)

		courses, err := mockRepo.GetBucketCourses(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCourses, courses)
	})

	t.Run("Test GetBucketLessons", func(t *testing.T) {
		expectedLessons := []*course.LessonPoint{}
		mockRepo.EXPECT().GetBucketLessons(ctx, 1, 1, 1).Return(expectedLessons, nil)

		lessons, err := mockRepo.GetBucketLessons(ctx, 1, 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedLessons, lessons)
	})

	t.Run("Test GetCourseParts", func(t *testing.T) {
		expectedParts := []*course.CoursePart{}
		mockRepo.EXPECT().GetCourseParts(ctx, 1).Return(expectedParts, nil)

		parts, err := mockRepo.GetCourseParts(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedParts, parts)
	})

	t.Run("Test GetCoursesPurchases", func(t *testing.T) {
		expectedPurchases := map[int]int{1: 1}
		mockRepo.EXPECT().GetCoursesPurchases(ctx, []*course.Course{}).Return(expectedPurchases, nil)

		purchases, err := mockRepo.GetCoursesPurchases(ctx, []*course.Course{})
		assert.NoError(t, err)
		assert.Equal(t, expectedPurchases, purchases)
	})

	t.Run("Test GetCoursesRaitings", func(t *testing.T) {
		expectedRatings := map[int]float32{1: 4.5}
		mockRepo.EXPECT().GetCoursesRaitings(ctx, []*course.Course{}).Return(expectedRatings, nil)

		ratings, err := mockRepo.GetCoursesRaitings(ctx, []*course.Course{})
		assert.NoError(t, err)
		assert.Equal(t, expectedRatings, ratings)
	})

	t.Run("Test GetCoursesTags", func(t *testing.T) {
		expectedTags := map[int][]string{1: {"tag1", "tag2"}}
		mockRepo.EXPECT().GetCoursesTags(ctx, []*course.Course{}).Return(expectedTags, nil)

		tags, err := mockRepo.GetCoursesTags(ctx, []*course.Course{})
		assert.NoError(t, err)
		assert.Equal(t, expectedTags, tags)
	})

	t.Run("Test GetFavouriteCourses", func(t *testing.T) {
		expectedCourses := []*course.Course{}
		mockRepo.EXPECT().GetFavouriteCourses(ctx, 1).Return(expectedCourses, nil)

		courses, err := mockRepo.GetFavouriteCourses(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedCourses, courses)
	})

	t.Run("Test GetLastLessonHeader", func(t *testing.T) {
		expectedHeader := &dto.LessonDtoHeader{}
		mockRepo.EXPECT().GetLastLessonHeader(ctx, 1, 1).Return(expectedHeader, 1, "test", true, nil)

		header, id, name, completed, err := mockRepo.GetLastLessonHeader(ctx, 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedHeader, header)
		assert.Equal(t, 1, id)
		assert.Equal(t, "test", name)
		assert.True(t, completed)
	})

	t.Run("Test GetLessonBlocks", func(t *testing.T) {
		expectedBlocks := []string{"block1", "block2"}
		mockRepo.EXPECT().GetLessonBlocks(ctx, 1).Return(expectedBlocks, nil)

		blocks, err := mockRepo.GetLessonBlocks(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedBlocks, blocks)
	})

	t.Run("Test GetLessonFooters", func(t *testing.T) {
		expectedFooters := []int{1, 2}
		mockRepo.EXPECT().GetLessonFooters(ctx, 1).Return(expectedFooters, nil)

		footers, err := mockRepo.GetLessonFooters(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedFooters, footers)
	})

	t.Run("Test GetLessonHeaderByLessonId", func(t *testing.T) {
		expectedHeader := &dto.LessonDtoHeader{}
		mockRepo.EXPECT().GetLessonHeaderByLessonId(ctx, 1, 1).Return(expectedHeader, nil)

		header, err := mockRepo.GetLessonHeaderByLessonId(ctx, 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedHeader, header)
	})

	t.Run("Test GetLessonTest", func(t *testing.T) {
		expectedTest := &dto.Test{}
		mockRepo.EXPECT().GetLessonTest(ctx, 1, 1).Return(expectedTest, nil)

		test, err := mockRepo.GetLessonTest(ctx, 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedTest, test)
	})

	t.Run("Test GetLessonVideo", func(t *testing.T) {
		expectedVideos := []string{"video1", "video2"}
		mockRepo.EXPECT().GetLessonVideo(ctx, 1).Return(expectedVideos, nil)

		videos, err := mockRepo.GetLessonVideo(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedVideos, videos)
	})

	t.Run("Test GetPartBuckets", func(t *testing.T) {
		expectedBuckets := []*course.LessonBucket{}
		mockRepo.EXPECT().GetPartBuckets(ctx, 1).Return(expectedBuckets, nil)

		buckets, err := mockRepo.GetPartBuckets(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedBuckets, buckets)
	})

	t.Run("Test GetQuestionTestLesson", func(t *testing.T) {
		expectedQuestion := &dto.QuestionTest{}
		mockRepo.EXPECT().GetQuestionTestLesson(ctx, 1, 1).Return(expectedQuestion, nil)

		question, err := mockRepo.GetQuestionTestLesson(ctx, 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedQuestion, question)
	})

	t.Run("Test GetUserById", func(t *testing.T) {
		expectedUser := &user.User{}
		mockRepo.EXPECT().GetUserById(ctx, 1).Return(expectedUser, nil)

		user, err := mockRepo.GetUserById(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})
}

// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/course/course_repository.go

// Package usecase is a generated GoMock package.
package usecase

import (
	context "context"
	io "io"
	reflect "reflect"
	coursemodels "skillForce/internal/models/course"
	dto "skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"

	gomock "github.com/golang/mock/gomock"
)

// MockCourseRepository is a mock of CourseRepository interface.
type MockCourseRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCourseRepositoryMockRecorder
}

// MockCourseRepositoryMockRecorder is the mock recorder for MockCourseRepository.
type MockCourseRepositoryMockRecorder struct {
	mock *MockCourseRepository
}

// NewMockCourseRepository creates a new mock instance.
func NewMockCourseRepository(ctrl *gomock.Controller) *MockCourseRepository {
	mock := &MockCourseRepository{ctrl: ctrl}
	mock.recorder = &MockCourseRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCourseRepository) EXPECT() *MockCourseRepositoryMockRecorder {
	return m.recorder
}

// AddUserToCourse mocks base method.
func (m *MockCourseRepository) AddUserToCourse(ctx context.Context, userId, courseId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToCourse", ctx, userId, courseId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToCourse indicates an expected call of AddUserToCourse.
func (mr *MockCourseRepositoryMockRecorder) AddUserToCourse(ctx, userId, courseId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToCourse", reflect.TypeOf((*MockCourseRepository)(nil).AddUserToCourse), ctx, userId, courseId)
}

// GetBucketByLessonId mocks base method.
func (m *MockCourseRepository) GetBucketByLessonId(ctx context.Context, lessonId int) (*coursemodels.LessonBucket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketByLessonId", ctx, lessonId)
	ret0, _ := ret[0].(*coursemodels.LessonBucket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBucketByLessonId indicates an expected call of GetBucketByLessonId.
func (mr *MockCourseRepositoryMockRecorder) GetBucketByLessonId(ctx, lessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketByLessonId", reflect.TypeOf((*MockCourseRepository)(nil).GetBucketByLessonId), ctx, lessonId)
}

// GetBucketCourses mocks base method.
func (m *MockCourseRepository) GetBucketCourses(ctx context.Context) ([]*coursemodels.Course, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketCourses", ctx)
	ret0, _ := ret[0].([]*coursemodels.Course)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBucketCourses indicates an expected call of GetBucketCourses.
func (mr *MockCourseRepositoryMockRecorder) GetBucketCourses(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketCourses", reflect.TypeOf((*MockCourseRepository)(nil).GetBucketCourses), ctx)
}

// GetBucketLessons mocks base method.
func (m *MockCourseRepository) GetBucketLessons(ctx context.Context, userId, courseId, bucketId int) ([]*coursemodels.LessonPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketLessons", ctx, userId, courseId, bucketId)
	ret0, _ := ret[0].([]*coursemodels.LessonPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBucketLessons indicates an expected call of GetBucketLessons.
func (mr *MockCourseRepositoryMockRecorder) GetBucketLessons(ctx, userId, courseId, bucketId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketLessons", reflect.TypeOf((*MockCourseRepository)(nil).GetBucketLessons), ctx, userId, courseId, bucketId)
}

// GetCourseById mocks base method.
func (m *MockCourseRepository) GetCourseById(ctx context.Context, courseId int) (*coursemodels.Course, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCourseById", ctx, courseId)
	ret0, _ := ret[0].(*coursemodels.Course)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCourseById indicates an expected call of GetCourseById.
func (mr *MockCourseRepositoryMockRecorder) GetCourseById(ctx, courseId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCourseById", reflect.TypeOf((*MockCourseRepository)(nil).GetCourseById), ctx, courseId)
}

// GetCourseParts mocks base method.
func (m *MockCourseRepository) GetCourseParts(ctx context.Context, courseId int) ([]*coursemodels.CoursePart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCourseParts", ctx, courseId)
	ret0, _ := ret[0].([]*coursemodels.CoursePart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCourseParts indicates an expected call of GetCourseParts.
func (mr *MockCourseRepositoryMockRecorder) GetCourseParts(ctx, courseId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCourseParts", reflect.TypeOf((*MockCourseRepository)(nil).GetCourseParts), ctx, courseId)
}

// GetCoursesPurchases mocks base method.
func (m *MockCourseRepository) GetCoursesPurchases(ctx context.Context, bucketCoursesWithoutPurchases []*coursemodels.Course) (map[int]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoursesPurchases", ctx, bucketCoursesWithoutPurchases)
	ret0, _ := ret[0].(map[int]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoursesPurchases indicates an expected call of GetCoursesPurchases.
func (mr *MockCourseRepositoryMockRecorder) GetCoursesPurchases(ctx, bucketCoursesWithoutPurchases interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoursesPurchases", reflect.TypeOf((*MockCourseRepository)(nil).GetCoursesPurchases), ctx, bucketCoursesWithoutPurchases)
}

// GetCoursesRaitings mocks base method.
func (m *MockCourseRepository) GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*coursemodels.Course) (map[int]float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoursesRaitings", ctx, bucketCoursesWithoutRating)
	ret0, _ := ret[0].(map[int]float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoursesRaitings indicates an expected call of GetCoursesRaitings.
func (mr *MockCourseRepositoryMockRecorder) GetCoursesRaitings(ctx, bucketCoursesWithoutRating interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoursesRaitings", reflect.TypeOf((*MockCourseRepository)(nil).GetCoursesRaitings), ctx, bucketCoursesWithoutRating)
}

// GetCoursesTags mocks base method.
func (m *MockCourseRepository) GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*coursemodels.Course) (map[int][]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoursesTags", ctx, bucketCoursesWithoutTags)
	ret0, _ := ret[0].(map[int][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoursesTags indicates an expected call of GetCoursesTags.
func (mr *MockCourseRepositoryMockRecorder) GetCoursesTags(ctx, bucketCoursesWithoutTags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoursesTags", reflect.TypeOf((*MockCourseRepository)(nil).GetCoursesTags), ctx, bucketCoursesWithoutTags)
}

// GetLastLessonHeader mocks base method.
func (m *MockCourseRepository) GetLastLessonHeader(ctx context.Context, userId, courseId int) (*dto.LessonDtoHeader, int, string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastLessonHeader", ctx, userId, courseId)
	ret0, _ := ret[0].(*dto.LessonDtoHeader)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(bool)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// GetLastLessonHeader indicates an expected call of GetLastLessonHeader.
func (mr *MockCourseRepositoryMockRecorder) GetLastLessonHeader(ctx, userId, courseId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastLessonHeader", reflect.TypeOf((*MockCourseRepository)(nil).GetLastLessonHeader), ctx, userId, courseId)
}

// GetLessonBlocks mocks base method.
func (m *MockCourseRepository) GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLessonBlocks", ctx, currentLessonId)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLessonBlocks indicates an expected call of GetLessonBlocks.
func (mr *MockCourseRepositoryMockRecorder) GetLessonBlocks(ctx, currentLessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLessonBlocks", reflect.TypeOf((*MockCourseRepository)(nil).GetLessonBlocks), ctx, currentLessonId)
}

// GetLessonById mocks base method.
func (m *MockCourseRepository) GetLessonById(ctx context.Context, lessonId int) (*coursemodels.LessonPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLessonById", ctx, lessonId)
	ret0, _ := ret[0].(*coursemodels.LessonPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLessonById indicates an expected call of GetLessonById.
func (mr *MockCourseRepositoryMockRecorder) GetLessonById(ctx, lessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLessonById", reflect.TypeOf((*MockCourseRepository)(nil).GetLessonById), ctx, lessonId)
}

// GetLessonFooters mocks base method.
func (m *MockCourseRepository) GetLessonFooters(ctx context.Context, currentLessonId int) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLessonFooters", ctx, currentLessonId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLessonFooters indicates an expected call of GetLessonFooters.
func (mr *MockCourseRepositoryMockRecorder) GetLessonFooters(ctx, currentLessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLessonFooters", reflect.TypeOf((*MockCourseRepository)(nil).GetLessonFooters), ctx, currentLessonId)
}

// GetLessonHeaderByLessonId mocks base method.
func (m *MockCourseRepository) GetLessonHeaderByLessonId(ctx context.Context, userId, currentLessonId int) (*dto.LessonDtoHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLessonHeaderByLessonId", ctx, userId, currentLessonId)
	ret0, _ := ret[0].(*dto.LessonDtoHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLessonHeaderByLessonId indicates an expected call of GetLessonHeaderByLessonId.
func (mr *MockCourseRepositoryMockRecorder) GetLessonHeaderByLessonId(ctx, userId, currentLessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLessonHeaderByLessonId", reflect.TypeOf((*MockCourseRepository)(nil).GetLessonHeaderByLessonId), ctx, userId, currentLessonId)
}

// GetLessonVideo mocks base method.
func (m *MockCourseRepository) GetLessonVideo(ctx context.Context, currentLessonId int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLessonVideo", ctx, currentLessonId)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLessonVideo indicates an expected call of GetLessonVideo.
func (mr *MockCourseRepositoryMockRecorder) GetLessonVideo(ctx, currentLessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLessonVideo", reflect.TypeOf((*MockCourseRepository)(nil).GetLessonVideo), ctx, currentLessonId)
}

// GetPartBuckets mocks base method.
func (m *MockCourseRepository) GetPartBuckets(ctx context.Context, partId int) ([]*coursemodels.LessonBucket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPartBuckets", ctx, partId)
	ret0, _ := ret[0].([]*coursemodels.LessonBucket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPartBuckets indicates an expected call of GetPartBuckets.
func (mr *MockCourseRepositoryMockRecorder) GetPartBuckets(ctx, partId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPartBuckets", reflect.TypeOf((*MockCourseRepository)(nil).GetPartBuckets), ctx, partId)
}

// GetUserById mocks base method.
func (m *MockCourseRepository) GetUserById(ctx context.Context, userId int) (*usermodels.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", ctx, userId)
	ret0, _ := ret[0].(*usermodels.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockCourseRepositoryMockRecorder) GetUserById(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockCourseRepository)(nil).GetUserById), ctx, userId)
}

// GetVideoRange mocks base method.
func (m *MockCourseRepository) GetVideoRange(ctx context.Context, name string, start, end int64) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVideoRange", ctx, name, start, end)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVideoRange indicates an expected call of GetVideoRange.
func (mr *MockCourseRepositoryMockRecorder) GetVideoRange(ctx, name, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideoRange", reflect.TypeOf((*MockCourseRepository)(nil).GetVideoRange), ctx, name, start, end)
}

// GetVideoUrl mocks base method.
func (m *MockCourseRepository) GetVideoUrl(ctx context.Context, lessonId int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVideoUrl", ctx, lessonId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVideoUrl indicates an expected call of GetVideoUrl.
func (mr *MockCourseRepositoryMockRecorder) GetVideoUrl(ctx, lessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideoUrl", reflect.TypeOf((*MockCourseRepository)(nil).GetVideoUrl), ctx, lessonId)
}

// IsUserPurchasedCourse mocks base method.
func (m *MockCourseRepository) IsUserPurchasedCourse(ctx context.Context, userId, courseId int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserPurchasedCourse", ctx, userId, courseId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserPurchasedCourse indicates an expected call of IsUserPurchasedCourse.
func (mr *MockCourseRepositoryMockRecorder) IsUserPurchasedCourse(ctx, userId, courseId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserPurchasedCourse", reflect.TypeOf((*MockCourseRepository)(nil).IsUserPurchasedCourse), ctx, userId, courseId)
}

// MarkLessonAsNotCompleted mocks base method.
func (m *MockCourseRepository) MarkLessonAsNotCompleted(ctx context.Context, userId, lessonId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkLessonAsNotCompleted", ctx, userId, lessonId)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkLessonAsNotCompleted indicates an expected call of MarkLessonAsNotCompleted.
func (mr *MockCourseRepositoryMockRecorder) MarkLessonAsNotCompleted(ctx, userId, lessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkLessonAsNotCompleted", reflect.TypeOf((*MockCourseRepository)(nil).MarkLessonAsNotCompleted), ctx, userId, lessonId)
}

// MarkLessonCompleted mocks base method.
func (m *MockCourseRepository) MarkLessonCompleted(ctx context.Context, userId, courseId, lessonId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkLessonCompleted", ctx, userId, courseId, lessonId)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkLessonCompleted indicates an expected call of MarkLessonCompleted.
func (mr *MockCourseRepositoryMockRecorder) MarkLessonCompleted(ctx, userId, courseId, lessonId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkLessonCompleted", reflect.TypeOf((*MockCourseRepository)(nil).MarkLessonCompleted), ctx, userId, courseId, lessonId)
}

// SendWelcomeCourseMail mocks base method.
func (m *MockCourseRepository) SendWelcomeCourseMail(ctx context.Context, user *usermodels.User, courseId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendWelcomeCourseMail", ctx, user, courseId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendWelcomeCourseMail indicates an expected call of SendWelcomeCourseMail.
func (mr *MockCourseRepositoryMockRecorder) SendWelcomeCourseMail(ctx, user, courseId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendWelcomeCourseMail", reflect.TypeOf((*MockCourseRepository)(nil).SendWelcomeCourseMail), ctx, user, courseId)
}

// Stat mocks base method.
func (m *MockCourseRepository) Stat(ctx context.Context, name string) (dto.VideoMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat", ctx, name)
	ret0, _ := ret[0].(dto.VideoMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *MockCourseRepositoryMockRecorder) Stat(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockCourseRepository)(nil).Stat), ctx, name)
}

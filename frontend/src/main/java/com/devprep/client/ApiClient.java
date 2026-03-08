package com.devprep.client;

import com.devprep.client.dto.*;

import java.util.List;
import java.util.Optional;

public interface ApiClient {

    List<TopicDto> getTopics();

    Optional<TopicWithQuestionsDto> getTopicBySlug(String slug, Level level);

    PaginatedQuestionsDto getQuestions(String topic, String tag, Level level, int page, int limit);

    Optional<QuestionDetailDto> getQuestionBySlug(String slug);

    List<TagDto> getTags();

    void updateProgress(String slug, ProgressStatus status);

    List<UserProgressDto> getMyProgress();

    List<TopicProgressDto> getMyProgressByTopic();

    Optional<ProgressStatus> getQuestionProgress(String slug);

    boolean toggleBookmark(String slug);

    List<BookmarkDto> getMyBookmarks();

    boolean isBookmarked(String slug);

    void recordView(String slug);

    List<ViewHistoryDto> getMyHistory();
}
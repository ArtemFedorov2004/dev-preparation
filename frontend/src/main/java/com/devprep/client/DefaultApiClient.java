package com.devprep.client;

import com.devprep.client.dto.*;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.MediaType;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.RestClient;
import org.springframework.web.util.UriComponentsBuilder;

import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Optional;

@Slf4j
@RequiredArgsConstructor
public class DefaultApiClient implements ApiClient {

    private static final ParameterizedTypeReference<List<TopicDto>> TOPICS_TYPE_REFERENCE =
            new ParameterizedTypeReference<>() {
            };

    private static final ParameterizedTypeReference<List<TagDto>> TAGS_TYPE_REFERENCE =
            new ParameterizedTypeReference<>() {
            };

    private static final ParameterizedTypeReference<List<UserProgressDto>> PROGRESS_LIST_TYPE_REFERENCE =
            new ParameterizedTypeReference<>() {
            };

    private static final ParameterizedTypeReference<List<TopicProgressDto>> TOPIC_PROGRESS_TYPE_REFERENCE =
            new ParameterizedTypeReference<>() {
            };

    private static final ParameterizedTypeReference<List<BookmarkDto>> BOOKMARKS_TYPE_REFERENCE =
            new ParameterizedTypeReference<>() {
            };

    private static final ParameterizedTypeReference<List<ViewHistoryDto>> HISTORY_TYPE_REFERENCE =
            new ParameterizedTypeReference<>() {
            };

    private static final ParameterizedTypeReference<Map<String, Object>> MAP_TYPE_REFERENCE =
            new ParameterizedTypeReference<>() {
            };

    private final RestClient restClient;

    @Override
    public List<TopicDto> getTopics() {
        try {
            return restClient.get()
                    .uri("/api/v1/topics")
                    .retrieve()
                    .body(TOPICS_TYPE_REFERENCE);
        } catch (Exception e) {
            log.error("Failed to fetch topics: {}", e.getMessage(), e.getCause());
            return Collections.emptyList();
        }
    }

    @Override
    public Optional<TopicWithQuestionsDto> getTopicBySlug(String slug, Level level) {
        try {
            String uri = UriComponentsBuilder.fromPath("/api/v1/topics/{slug}")
                    .queryParamIfPresent("level", Optional.ofNullable(level).map(Level::getValue))
                    .buildAndExpand(slug)
                    .toUriString();
            return Optional.ofNullable(restClient.get()
                    .uri(uri)
                    .retrieve()
                    .body(TopicWithQuestionsDto.class));
        } catch (HttpClientErrorException.NotFound exception) {
            return Optional.empty();
        } catch (Exception e) {
            log.error("Error calling GET /api/v1/topics/{}", slug, e);
            return Optional.empty();
        }
    }

    @Override
    public PaginatedQuestionsDto getQuestions(String topic, String tag, Level level, int page, int limit) {
        try {
            String uri = UriComponentsBuilder.fromPath("/api/v1/questions")
                    .queryParamIfPresent("topic", Optional.ofNullable(topic).filter(s -> !s.isBlank()))
                    .queryParamIfPresent("tag", Optional.ofNullable(tag).filter(s -> !s.isBlank()))
                    .queryParamIfPresent("level", Optional.ofNullable(level).map(Level::getValue))
                    .queryParam("page", page)
                    .queryParam("limit", limit)
                    .build()
                    .toUriString();

            PaginatedQuestionsDto result = restClient.get()
                    .uri(uri)
                    .retrieve()
                    .body(PaginatedQuestionsDto.class);
            return result != null ? result : emptyPaginated(page, limit);
        } catch (Exception e) {
            log.error("Error calling GET /api/v1/questions", e);
            return emptyPaginated(page, limit);
        }
    }

    @Override
    public Optional<QuestionDetailDto> getQuestionBySlug(String slug) {
        try {
            return Optional.ofNullable(restClient.get()
                    .uri("/api/v1/questions/{slug}", slug)
                    .retrieve()
                    .body(QuestionDetailDto.class));
        } catch (HttpClientErrorException.NotFound exception) {
            return Optional.empty();
        } catch (Exception e) {
            log.error("Error calling GET /api/v1/questions/{}", slug, e);
            return Optional.empty();
        }
    }

    @Override
    public List<TagDto> getTags() {
        try {
            return restClient.get()
                    .uri("/api/v1/tags")
                    .retrieve()
                    .body(TAGS_TYPE_REFERENCE);
        } catch (Exception e) {
            log.error("Error calling GET /api/v1/tags", e);
            return Collections.emptyList();
        }
    }

    @Override
    public void updateProgress(String slug, ProgressStatus status) {
        try {
            restClient.post()
                    .uri("/api/v1/questions/{slug}/progress", slug)
                    .contentType(MediaType.APPLICATION_JSON)
                    .body(new UpdateProgressRequest(status.getValue()))
                    .retrieve()
                    .toBodilessEntity();
        } catch (HttpClientErrorException.NotFound e) {
            log.warn("Question not found when updating progress: slug={}", slug);
        } catch (Exception e) {
            log.error("Error POST /api/v1/questions/{}/progress", slug, e);
        }
    }

    @Override
    public List<UserProgressDto> getMyProgress() {
        try {
            List<UserProgressDto> result = restClient.get()
                    .uri("/api/v1/me/progress")
                    .retrieve()
                    .body(PROGRESS_LIST_TYPE_REFERENCE);
            return result != null ? result : Collections.emptyList();
        } catch (Exception e) {
            log.error("Error GET /api/v1/me/progress", e);
            return Collections.emptyList();
        }
    }

    @Override
    public List<TopicProgressDto> getMyProgressByTopic() {
        try {
            List<TopicProgressDto> result = restClient.get()
                    .uri("/api/v1/me/progress/by-topic")
                    .retrieve()
                    .body(TOPIC_PROGRESS_TYPE_REFERENCE);
            return result != null ? result : Collections.emptyList();
        } catch (Exception e) {
            log.error("Error GET /api/v1/me/progress/by-topic", e);
            return Collections.emptyList();
        }
    }

    @Override
    public Optional<ProgressStatus> getQuestionProgress(String slug) {
        try {
            UserProgressDto dto = restClient.get()
                    .uri("/api/v1/questions/{slug}/progress", slug)
                    .retrieve()
                    .body(UserProgressDto.class);
            return Optional.ofNullable(dto).map(UserProgressDto::getStatus);
        } catch (HttpClientErrorException.NotFound e) {
            return Optional.empty();
        } catch (Exception e) {
            log.error("Error GET /api/v1/questions/{}/progress", slug, e);
            return Optional.empty();
        }
    }

    @Override
    public boolean toggleBookmark(String slug) {
        try {
            Map<String, Object> response = restClient.post()
                    .uri("/api/v1/questions/{slug}/bookmark", slug)
                    .retrieve()
                    .body(MAP_TYPE_REFERENCE);
            return response != null && Boolean.TRUE.equals(response.get("bookmarked"));
        } catch (HttpClientErrorException.NotFound e) {
            log.warn("Question not found when toggling bookmark: slug={}", slug);
            return false;
        } catch (Exception e) {
            log.error("Error POST /api/v1/questions/{}/bookmark", slug, e);
            return false;
        }
    }

    @Override
    public List<BookmarkDto> getMyBookmarks() {
        try {
            List<BookmarkDto> result = restClient.get()
                    .uri("/api/v1/me/bookmarks")
                    .retrieve()
                    .body(BOOKMARKS_TYPE_REFERENCE);
            return result != null ? result : Collections.emptyList();
        } catch (Exception e) {
            log.error("Error GET /api/v1/me/bookmarks", e);
            return Collections.emptyList();
        }
    }

    @Override
    public boolean isBookmarked(String slug) {
        try {
            Map<String, Object> response = restClient.get()
                    .uri("/api/v1/questions/{slug}/bookmark", slug)
                    .retrieve()
                    .body(MAP_TYPE_REFERENCE);
            return response != null && Boolean.TRUE.equals(response.get("bookmarked"));
        } catch (HttpClientErrorException.NotFound e) {
            return false;
        } catch (Exception e) {
            log.error("Error GET /api/v1/questions/{}/bookmark", slug, e);
            return false;
        }
    }

    @Override
    public void recordView(String slug) {
        try {
            restClient.post()
                    .uri("/api/v1/questions/{slug}/view", slug)
                    .retrieve()
                    .toBodilessEntity();
        } catch (HttpClientErrorException.NotFound e) {
            log.warn("Question not found when recording view: slug={}", slug);
        } catch (Exception e) {
            log.warn("Error POST /api/v1/questions/{}/view: {}", slug, e.getMessage());
        }
    }

    @Override
    public List<ViewHistoryDto> getMyHistory() {
        try {
            List<ViewHistoryDto> result = restClient.get()
                    .uri("/api/v1/me/history")
                    .retrieve()
                    .body(HISTORY_TYPE_REFERENCE);
            return result != null ? result : Collections.emptyList();
        } catch (Exception e) {
            log.error("Error GET /api/v1/me/history", e);
            return Collections.emptyList();
        }
    }

    private PaginatedQuestionsDto emptyPaginated(int page, int limit) {
        PaginatedQuestionsDto empty = new PaginatedQuestionsDto();
        empty.setData(Collections.emptyList());
        PaginatedQuestionsDto.PaginationDto pagination = new PaginatedQuestionsDto.PaginationDto();
        pagination.setPage(page);
        pagination.setLimit(limit);
        pagination.setTotal(0);
        empty.setPagination(pagination);
        return empty;
    }
}
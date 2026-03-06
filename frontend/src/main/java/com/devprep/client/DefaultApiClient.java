package com.devprep.client;

import com.devprep.client.dto.*;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.ResourceAccessException;
import org.springframework.web.client.RestClient;
import org.springframework.web.util.UriComponentsBuilder;

import java.util.Collections;
import java.util.List;
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

    private final RestClient restClient;

    @Override
    public List<TopicDto> getTopics() {
        try {
            return restClient.get()
                    .uri("/api/v1/topics")
                    .retrieve()
                    .body(TOPICS_TYPE_REFERENCE);
        } catch (ResourceAccessException e) {
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
        } catch (ResourceAccessException e) {
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
        } catch (ResourceAccessException e) {
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
        } catch (ResourceAccessException e) {
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
        } catch (ResourceAccessException e) {
            log.error("Error calling GET /api/v1/tags", e);
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
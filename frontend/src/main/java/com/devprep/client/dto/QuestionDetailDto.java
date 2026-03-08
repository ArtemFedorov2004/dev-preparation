package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public record QuestionDetailDto(
        int id, String slug, String title, String answer,
        Level level, TopicDto topic, List<TagDto> tags
) {
}
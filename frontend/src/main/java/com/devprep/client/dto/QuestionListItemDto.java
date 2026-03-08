package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public record QuestionListItemDto(int id, String slug, String title,
                                  Level level, TopicDto topic, List<TagDto> tags
) {
}
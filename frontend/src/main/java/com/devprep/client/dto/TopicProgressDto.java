package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public record TopicProgressDto(
        TopicDto topic, int total, int learned, int needReview, int dontKnow
) {
}
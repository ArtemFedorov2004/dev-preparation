package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record TopicProgressDto(
        TopicDto topic, int total, int learned,
        @JsonProperty("need_review")
        int needReview,
        @JsonProperty("dont_know")
        int dontKnow
) {
}